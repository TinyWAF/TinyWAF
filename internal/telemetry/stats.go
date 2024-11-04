package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TinyWAF/TinyWAF/internal"
	"github.com/TinyWAF/TinyWAF/internal/logger"
)

// Aggregate number of requests/blocked requests since the last report was sent
type statNumbers struct {
	numRequests uint
	numBlocked  uint
}

var stats statNumbers
var sleepDuration time.Duration

const defaultReportUrl string = "https://tinywaf.com/api.php?p=/v1/anonymous-stats"

func Init() {
	resetStats()
}

func AddRequest() {
	stats.numRequests++
}

func AddBlocked() {
	stats.numBlocked++
}

func Start(cfg *internal.MainConfig, appVersion string) {
	sleepDuration = time.Duration(cfg.Stats.IntervalSecs) * time.Second
	if sleepDuration == 0 {
		sleepDuration = 300 // Default 5 minutes
	}

	go func() {
		for {
			time.Sleep(sleepDuration)
			sendReport(cfg, appVersion)
		}
	}()
}

func sendReport(cfg *internal.MainConfig, appVersion string) {
	logger.Debug("Sending stats: requests: %v, blocked: %v", stats.numRequests, stats.numBlocked)

	postBody, _ := json.Marshal(map[string]uint{
		"requests": stats.numRequests,
		"blocked":  stats.numBlocked,
	})
	responseBody := bytes.NewBuffer(postBody)

	reportUrl := defaultReportUrl
	if cfg.Stats.PostUrl != "" {
		reportUrl = cfg.Stats.PostUrl
	}

	// resp, err := http.Post(reportUrl, "application/json", responseBody)

	client := &http.Client{}
	req, err := http.NewRequest("POST", reportUrl, responseBody)
	if err != nil {
		logger.Error("Failed to send anonymous usage stats to '%s' backing off: %v", reportUrl, err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", fmt.Sprintf("tinywaf/%s", appVersion))

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to send anonymous usage stats to '%s' backing off: %v", reportUrl, err)
		backOff()
		return
	}

	if resp.StatusCode != http.StatusAccepted {
		logger.Error("Failed to send anonymous usage stats to '%s' (HTTP response: %s) backing off", reportUrl, resp.Status)
		backOff()
		return
	}

	// Only reset stats if the request was successful, otherwise we want to try
	// sending them again in the next request attempt
	resetStats()
}

func resetStats() {
	stats = statNumbers{
		numRequests: 0,
		numBlocked:  0,
	}
}

// If the request failed, wait longer before trying again. Double the wait time
// up to a maximum of 1 hour.
func backOff() {
	sleepDuration = sleepDuration * 2

	if sleepDuration > time.Hour {
		sleepDuration = time.Hour
	}
}
