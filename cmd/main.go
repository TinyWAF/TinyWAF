package main

import (
	"log"

	"github.com/TinyWAF/TinyWAF/internal"
	"github.com/TinyWAF/TinyWAF/internal/config"
	"github.com/TinyWAF/TinyWAF/internal/logger"
	"github.com/TinyWAF/TinyWAF/internal/ruleengine"
	"github.com/TinyWAF/TinyWAF/internal/telemetry"
	"github.com/TinyWAF/TinyWAF/internal/webserver"
)

// Build script should tag the commit with the version, then run this:
// go build -i -v -ldflags="-X main.version=$(git describe --always --long --dirty)" github.com/TinyWAF/TinyWAF
var version = "undefined"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load TinyWAF config: %v", err.Error())
		return
	}

	setupDirs(&cfg)
	logger.Init(&cfg)
	telemetry.Init()

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

	if cfg.Stats.Enabled {
		telemetry.Start(&cfg, version)
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
