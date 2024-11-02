package config

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Rules struct {
	RequestRules  []RuleGroup `validate:"required"`
	ResponseRules []RuleGroup
}

type RuleGroup struct {
	Group string `validate:"required"`
	Rules []Rule `validate:"required"`
}

type Rule struct {
	Id        string            `validate:"required"`
	Inspect   []string          `validate:"required,gt=0,dive,oneof=url headers body"`
	Fields    []string          `validate:"omitempty,gt=0"`
	Operators map[string]string `validate:"required,gt=0,dive,keys,oneof=contains notcontains exactly notexactly regex notregex,endkeys"`
	Ratelimit struct {
		MaxAllowedRequests int `validate:"omitempty,gt=0"`
		WithinMinutes      int `validate:"omitempty,gt=0"`
	}
	Action string `validate:"required,oneof=ignore warn block"`
}

var ErrNoFirewallRulesLoaded = errors.New("No firewall rules loaded")

var RuleActionIgnore = "ignore"
var RuleActionWarn = "warn"
var RuleActionBlock = "block"

var RuleInspectUrl = "url"
var RuleInspectHeaders = "headers"
var RuleInspectBody = "body"

// contains notcontains exactly notexactly regex notregex
var RuleOperatorContains = "contains"
var RuleOperatorNotContains = "notcontains"
var RuleOperatorExactly = "exactly"
var RuleOperatorNotExactly = "notexactly"
var RuleOperatorRegex = "regex"
var RuleOperatorNotRegex = "notregex"

func LoadRules(cfg *MainConfig) (Rules, error) {
	return loadRules(viper.New(), cfg)
}

func loadRules(v *viper.Viper, cfg *MainConfig) (Rules, error) {
	for _, globPattern := range cfg.RuleFiles.Request.Src {
		ruleFilePaths, err := filepath.Glob(globPattern)
		if err != nil {
			log.Printf("ERROR: Failed to glob request rule files matching '%v', skipping: %v", globPattern, err.Error())
			continue
		}

		for _, filePath := range ruleFilePaths {
			v.SetConfigFile(filePath)

			log.Printf("Loading ruleset from '%v'...", filePath)

			err := v.MergeInConfig()
			if err != nil {
				// Probably failed to read in the file
				log.Printf("ERROR: Failed to load rule file '%v', skipping: %v", filePath, err.Error())
				continue
			}
		}
	}

	// @todo: load in response rules

	rules := Rules{}

	requestRuleGroups := []RuleGroup{}
	err := v.Unmarshal(&requestRuleGroups)
	if err != nil {
		return rules, err
	}

	if len(requestRuleGroups) == 0 || (len(requestRuleGroups) == 1 && len(requestRuleGroups[0].Rules) == 0) {
		return rules, ErrNoFirewallRulesLoaded
	}

	rules.RequestRules = requestRuleGroups

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(rules)
	if err != nil {
		return rules, err
	}

	log.Println("Firewall rules loaded")

	return rules, nil
}
