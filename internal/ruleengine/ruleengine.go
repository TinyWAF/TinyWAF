package ruleengine

import (
	"github.com/TinyWAF/TinyWAF/internal"
	"github.com/TinyWAF/TinyWAF/internal/config"
)

var loadedRules *config.Rules
var loadedCfg *internal.MainConfig

// Start the rule engine. Load rules from files defined in config
func Init(cfg *internal.MainConfig) error {
	rules, err := config.LoadRules(cfg)
	if err != nil {
		// Most likely the rule config validation failed
		return err
	}

	loadedCfg = cfg
	loadedRules = &rules

	// @todo: implement a way to reload config without restarting app - like apache can

	return nil
}
