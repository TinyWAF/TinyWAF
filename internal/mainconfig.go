package internal

// @todo: check validation rules are correct
type MainConfig struct {
	Listen struct {
		Hosts         []string `validate:"required,gt=0,dive,hostname_port"`
		Autotls       bool     `validate:"boolean"`
		ForwardToPort uint     `validate:"required,gt=0"`
	}

	Log struct {
		File   string `validate:"omitempty,filepath"`
		Levels struct {
			Debug bool `validate:"boolean"`
			Warn  bool `validate:"boolean"`
			Block bool `validate:"boolean"`
		}
	}

	Html struct {
		Blocked     string `validate:"omitempty,filepath"`
		Unavailable string `validate:"omitempty,filepath"`
	}

	Rulesets struct {
		InspectOnly bool     `validate:"boolean"`
		Include     []string `validate:"dive,filepath"`
		Overrides   []RuleOverride
	}
}

type RuleOverride struct {
	Host    string `validate:"hostname_port"`
	Enable  []string
	Disable []string
}
