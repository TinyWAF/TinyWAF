package ruleengine

import "github.com/TinyWAF/TinyWAF/internal/config"

// Start the rule engine. Load rules from files defined in config
func Init(cfg *config.MainConfig) error {
	_, err := config.LoadRules(cfg)
	if err != nil {
		return err
	}

	// @todo: implement a way to reload config without restarting app - like apache can

	return nil
}
