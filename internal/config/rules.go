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
	Id          string            `validate:"required"`
	Inspect     []string          `validate:"required,gt=0,dive,oneof=url headers body ip"`
	WhenMethods []string          `validate:"omitempty,gt=0"`
	Fields      []string          `validate:"omitempty,gt=0"`
	Operators   map[string]string `validate:"required,gt=0,dive,keys,oneof=contains notcontains exactly notexactly regex notregex,endkeys"`
	Ratelimit   struct {
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

var RuleOperatorContains = "contains"
var RuleOperatorNotContains = "notcontains"
var RuleOperatorExactly = "exactly"
var RuleOperatorNotExactly = "notexactly"
var RuleOperatorRegex = "regex"
var RuleOperatorNotRegex = "notregex"

func LoadRules(cfg *MainConfig) (Rules, error) {
	rules := Rules{}
	requestRuleGroups := []RuleGroup{}
	responseRuleGroups := []RuleGroup{}
	numRulesLoaded := 0

	for _, globPattern := range cfg.RuleFiles.Request.Src {
		rules, numLoaded := loadRulesFromGlob(globPattern)
		requestRuleGroups = append(requestRuleGroups, rules...)
		numRulesLoaded += numLoaded
	}

	for _, globPattern := range cfg.RuleFiles.Response.Src {
		rules, numLoaded := loadRulesFromGlob(globPattern)
		responseRuleGroups = append(responseRuleGroups, rules...)
		numRulesLoaded += numLoaded
	}

	if numRulesLoaded == 0 {
		return rules, ErrNoFirewallRulesLoaded
	}

	rules.RequestRules = requestRuleGroups
	rules.ResponseRules = responseRuleGroups

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(rules)
	if err != nil {
		return rules, err
	}

	log.Printf("Loaded %v rules successfully", numRulesLoaded)

	return rules, nil
}

func loadRulesFromGlob(globPattern string) ([]RuleGroup, int) {
	v := viper.New()
	loadedRuleGroups := []RuleGroup{}
	numRulesLoaded := 0

	ruleFilePaths, err := filepath.Glob(globPattern)
	if err != nil {
		log.Printf("ERROR: Failed to glob request rule files matching '%v', skipping: %v", globPattern, err.Error())
		return loadedRuleGroups, numRulesLoaded
	}

	for _, filePath := range ruleFilePaths {
		log.Printf("Loading ruleset from '%v'...", filePath)

		v.SetConfigFile(filePath)
		v.MergeInConfig()

		rulesForFile := RuleGroup{}
		err := v.Unmarshal(&rulesForFile)
		if err != nil {
			// Unable to parse yaml file
			log.Printf("ERROR: Failed to parse yaml in rule file '%v', skipping: %v", filePath, err.Error())
			continue
		}

		numRulesLoaded += len(rulesForFile.Rules)
		loadedRuleGroups = append(loadedRuleGroups, rulesForFile)
	}

	return loadedRuleGroups, numRulesLoaded
}
