package config

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Rules struct {
	RequestRules  []RuleGroup `validate:"required"`
	ResponseRules []RuleGroup
}

type RuleGroup struct {
	Group string `validate:"required"`
	Rules []Rule `validate:"required"`
}

type Rule struct {
	Id        string            `validate:"required"`
	Inspect   []string          `validate:"required,gt=0"`
	Operators map[string]string `validate:"required,gt=0"`
	Ratelimit struct {
		MaxAllowedRequests int `validate:"omitempty,gt=0"`
		WithinMinutes      int `validate:"omitempty,gt=0"`
	}
	Action string `validate:"required"`
}
	return rules, nil
}
