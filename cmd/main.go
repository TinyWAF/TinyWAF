package main

import (
	"log"

	"github.com/TinyWAF/TinyWAF/internal/ruleengine"
	"github.com/TinyWAF/TinyWAF/internal/webserver"
)

func main() {
	// @todo: load tinywaf config
	// @todo: load firewall rules

	// @todo: handle first-run (eg. create log files, data dir)

	err := ruleengine.Init()
	if err != nil {
		log.Fatalf("Failed to load TinyWAF rules: %v", err.Error())
		return
	}

	// @todo: do periodic cleanup of RuleEngine request memory

	err = webserver.Start()
	if err != nil {
		log.Fatalf("Failed to start TinyWAF: %v", err.Error())
		return
	}
}
