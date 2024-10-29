package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// @todo: check validation rules are correct
type MainConfig struct {
	Listen struct {
		Ips        []string     `validate:"required,gt=0"`
		Ports      []ListenPort `validate:"required,gt=0"`
		Websockets bool         `validate:"required"`
		Tls        struct {
			// @todo: TLS certificate config
		}
	}

	Upstream struct {
		Destination string `validate:"required"`
		Headers     struct {
			CopyAll bool `validate:"required"`
		}
	}

	// @todo: validations
	Log struct {
		File   string
		Levels struct {
			debug  bool
			access bool
			warn   bool
			block  bool
		}
	}
}

type ListenPort struct {
	Source uint `validate:"required"`
	Target uint
}

// var MainConfiguration *MainConfig

func LoadConfig() (MainConfig, error) {
	return loadConfig(viper.New())
}

func loadConfig(v *viper.Viper) (MainConfig, error) {
	// @todo: log which config files are loaded
	// @todo: load rules config

	v.SetConfigType("yaml")
	v.SetConfigName("tinywaf")       // file called tinywaf.yml|yaml
	v.AddConfigPath("/etc/tinywaf/") // in this dir, or...
	v.AddConfigPath("./data/")       // in data directory

	// @todo: listen.ports array should accept ints or ListenPort struct (if int, then source+target is same)

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
