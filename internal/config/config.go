package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// @todo: check validation rules are correct
type MainConfig struct {
	Listen struct {
		Ips             []string     `validate:"required,gt=0,dive,ip"`
		Ports           []ListenPort `validate:"required,gt=0"`
		Websockets      bool
		HealthcheckPath string
		Tls             struct {
			// @todo: TLS certificate config
		}
	}

	Upstream struct {
		Destination string `validate:"required"`
	}

	// @todo: validations
	Log struct {
		Outfile string `validate:"required,filepath"`
		Levels  struct {
			Access bool
			Warn   bool
			Block  bool
		}
	}

	RequestMemory struct {
		MaxAgeMinutes int `validate:"required"`
		MaxSize       int `validate:"required"`
	}

	Html struct {
		Blocked     string `validate:"omitempty,filepath"`
		Ratelimit   string `validate:"omitempty,filepath"`
		Unavailable string `validate:"omitempty,filepath"`
	}

	RuleFiles struct {
		Request struct {
			Src       []string `validate:"dive,filepath"`
			Overrides []RuleOverride
		}
		Response struct {
			Src       []string `validate:"dive,filepath"`
			Overrides []RuleOverride
		}
	}
}

type ListenPort struct {
	Source uint `validate:"required,gt=0"`
	Target uint `validate:"omitempty,gt=0"`
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
