package ruleengine

// var rules

// Start the rule engine
// Load rules from files defined in config
func Init() error {
	err := loadRules()

	return err
}
