package internal

type MainConfig struct {
	Listen struct {
		Hosts                []string `validate:"omitempty,gt=0"`
		TlsDomains           []string `validate:"omitempty,gt=0"`
		ForwardToLocalPort   uint     `validate:"required,gt=0"`
		StripResponseHeaders []string
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

	Stats struct {
		Enabled      bool `validate:"boolean"`
		PostUrl      string
		IntervalSecs uint
	}

	Rulesets struct {
		InspectOnly bool     `validate:"boolean"`
		Include     []string `validate:"dive,filepath"`
		Overrides   []RuleOverride
	}
}

type RuleOverride struct {
	Host    string `validate:"hostname"`
	Enable  []string
	Disable []string
}
