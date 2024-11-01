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
	Inspect   []string          `validate:"required,gt=0"`
	Operators map[string]string `validate:"required,gt=0"`
	Ratelimit struct {
		MaxAllowedRequests int `validate:"omitempty,gt=0"`
		WithinMinutes      int `validate:"omitempty,gt=0"`
	}
	Action string `validate:"required"`
}

var ErrNoFirewallRulesLoaded = errors.New("No firewall rules loaded")

func LoadRules(cfg *MainConfig) (Rules, error) {
	return loadRules(viper.New(), cfg)
}

func loadRules(v *viper.Viper, cfg *MainConfig) (Rules, error) {
	for i := range cfg.RuleFiles.Request.Src {
		ruleFilePaths, err := filepath.Glob(cfg.RuleFiles.Request.Src[i])
		if err != nil {
			log.Printf("ERROR: Failed to glob request rule files matching '%v', skipping: %v", cfg.RuleFiles.Request.Src[i], err.Error())
			continue
		}

		for j := range ruleFilePaths {
			v.SetConfigFile(ruleFilePaths[i])

			log.Printf("Loading ruleset from '%v'...", ruleFilePaths[i])

			err := v.MergeInConfig()
			if err != nil {
				// Probably failed to read in the file
				log.Printf("ERROR: Failed to load rule file '%v', skipping: %v", ruleFilePaths[j], err.Error())
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

	log.Println(requestRuleGroups)

	rules.RequestRules = requestRuleGroups

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(rules)
	if err != nil {
		return rules, err
	}

	return rules, nil
}
