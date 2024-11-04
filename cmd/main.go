package main

import (
	"log"

	"github.com/TinyWAF/TinyWAF/internal"
	"github.com/TinyWAF/TinyWAF/internal/config"
	"github.com/TinyWAF/TinyWAF/internal/logger"
	"github.com/TinyWAF/TinyWAF/internal/ruleengine"
	"github.com/TinyWAF/TinyWAF/internal/webserver"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load TinyWAF config: %v", err.Error())
		return
	}

	setupDirs(&cfg)
	logger.Init(&cfg)

	err = ruleengine.Init(&cfg)
	if err != nil {
		logger.Fatal("Failed to load TinyWAF rules: %v", err.Error())
		return
	}

	done := make(chan struct{})

	err = webserver.Start(&cfg)
	if err != nil {
		logger.Fatal("Failed to start TinyWAF: %v", err.Error())
		return
	}

	// Wait to receive something from the done channel before closing the app
	select {
	case <-done:
		return
	}
}

func setupDirs(config *internal.MainConfig) {
	// @todo: create log dir based on path of `config.Log.Outfile`
}
