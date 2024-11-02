package main

import (
	"log"

	"github.com/TinyWAF/TinyWAF/internal"
	"github.com/TinyWAF/TinyWAF/internal/config"
	"github.com/TinyWAF/TinyWAF/internal/ruleengine"
	"github.com/TinyWAF/TinyWAF/internal/webserver"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load TinyWAF config: %v", err.Error())
		return
	}

	// @todo: set up log vars based on log levels defined in config

	setupDirs(cfg)

	err = ruleengine.Init(&cfg)
	if err != nil {
		log.Fatalf("Failed to load TinyWAF rules: %v", err.Error())
		return
	}

	// @todo: do periodic cleanup of RuleEngine request memory

	err = webserver.Start(&cfg)
	if err != nil {
		log.Fatalf("Failed to start TinyWAF: %v", err.Error())
		return
	}
}

func setupDirs(config *internal.MainConfig) {
	// @todo: create log dir based on path of `config.Log.Outfile`
}
