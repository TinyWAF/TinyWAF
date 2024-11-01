package ruleengine

import "github.com/TinyWAF/TinyWAF/internal/config"

// Start the rule engine. Load rules from files defined in config
	// @todo: load firewall rules
	err := loadRules()
func Init(cfg *config.MainConfig) error {

	return nil
}
