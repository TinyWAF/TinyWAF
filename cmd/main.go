package main

import (
	"log"

	"github.com/TinyWAF/TinyWAF/internal/webserver"
)

func main() {
	// @todo: load tinywaf config
	// @todo: load firewall rules
	// @todo: connect to local sqlite db (rate limiting)

	// @todo: handle first-run (eg. create log files, data dir, sqlite DB)

	err := webserver.Start()
	if err != nil {
		log.Fatalf("Failed to start TinyWAF: %v", err.Error())
	}
}
