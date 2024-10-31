package ruleengine

import "github.com/TinyWAF/TinyWAF/internal/config"

// var rules

// Start the rule engine. Load rules from files defined in config
func Init(config config.MainConfig) error {
	// @todo: load firewall rules
	err := loadRules()

	return err
}
