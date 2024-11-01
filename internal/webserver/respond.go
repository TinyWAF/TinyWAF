package webserver

import (
	"fmt"
	"net/http"

	"github.com/TinyWAF/TinyWAF/internal/ruleengine"
)

func respondUnavailable(w http.ResponseWriter) {
	// @todo: if config custom HTML, load it and return that as the body
	responseBody := getHtmlResponseBody(
		http.StatusText(http.StatusServiceUnavailable),
		fmt.Sprintf(
			"%v %s",
			http.StatusServiceUnavailable,
			http.StatusText(http.StatusServiceUnavailable),
		),
	)

	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprint(w, responseBody)
}

func respondBlocked(inspection ruleengine.InspectionResult, w http.ResponseWriter) {
	// @todo: if config custom HTML, load it and return that as the body
	responseBody := getHtmlResponseBody(
		"Request blocked",
		fmt.Sprintf("Request blocked by firewall. InspectionID: %s", inspection.InspectionId),
	)

	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, responseBody)
}

func respondRateLimited(inspection ruleengine.InspectionResult, w http.ResponseWriter) {
	// @todo: if config custom HTML, load it and return that as the body
	responseBody := getHtmlResponseBody(
		"Too many requests",
		fmt.Sprintf("Too many requests - try again later. InspectionID: %s", inspection.InspectionId),
	)

	w.WriteHeader(http.StatusTooManyRequests)
	fmt.Fprint(w, responseBody)
}

func getHtmlResponseBody(title string, msg string) string {
	return fmt.Sprintf("<html><head><title>%s</title></head><body>%s</body></html>", title, msg)
}
