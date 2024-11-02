package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// @todo: check validation rules are correct
type MainConfig struct {
	Listen struct {
		Hosts           []ListenHost `validate:"required,gt=0"`
		HealthcheckPath string

		Upstream struct {
			Destination string `validate:"required"`
		}
	}

	Log struct {
		File   string `validate:"required,filepath"`
		Levels struct {
			Debug bool `validate:"boolean"`
			Warn  bool `validate:"boolean"`
			Block bool `validate:"boolean"`
		}
	}

	RequestMemory struct {
		Enabled       bool `validate:"boolean"`
		MaxAgeMinutes int  `validate:"required"`
		MaxSize       int  `validate:"required"`
	}

	Html struct {
		Blocked     string `validate:"omitempty,file"`
		Ratelimit   string `validate:"omitempty,file"`
		Unavailable string `validate:"omitempty,file"`
	}

	RuleFiles struct {
		WarnInsteadOfBlock bool
		Request            struct {
			Src       []string `validate:"dive,filepath"`
			Overrides []RuleOverride
		}
		Response struct {
			Src       []string `validate:"dive,filepath"`
			Overrides []RuleOverride
		}
	}
}

type ListenHost struct {
	Host         string `validate:"required,host_port"`
	UpstreamPort uint   `validate:"omitempty,gt=0"`
	Tls          struct {
		CertificatePath string `validate:"omitempty,file"`
		KeyPath         string `validate:"omitempty,file"`
	}
}

type RuleOverride struct {
	Path   string
	Rule   string
	Action string
}

func LoadConfig() (MainConfig, error) {
	return loadConfig(viper.New())
}

func loadConfig(v *viper.Viper) (MainConfig, error) {
	// @todo: log which config files are loaded

	v.SetConfigType("yaml")
	v.SetConfigName("tinywaf")       // file called tinywaf.yml|yaml
	v.AddConfigPath("/etc/tinywaf/") // in this dir, or...
	v.AddConfigPath("./data/")       // in data directory

	err := v.ReadInConfig()
	if err != nil {
		return MainConfig{}, err
	}
	config := MainConfig{}
	err = v.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(config); err != nil {
		return config, err
	}

	return config, nil
}
