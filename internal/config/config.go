package config

import (
	"log"

	"github.com/TinyWAF/TinyWAF/internal"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func LoadConfig() (internal.MainConfig, error) {
	return loadConfig(viper.New())
}

func loadConfig(v *viper.Viper) (internal.MainConfig, error) {
	v.SetConfigType("yaml")
	v.SetConfigName("tinywaf")       // file called tinywaf.yml|yaml
	v.AddConfigPath("/etc/tinywaf/") // in this dir, or...
	v.AddConfigPath("./data/")       // in data directory

	err := v.ReadInConfig()
	if err != nil {
		return internal.MainConfig{}, err
	}

	log.Printf("Loading config from '%v'...", v.ConfigFileUsed())

	config := internal.MainConfig{}
	err = v.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(config); err != nil {
		return config, err
	}

	log.Println("Config loaded successfully")

	return config, nil
}
