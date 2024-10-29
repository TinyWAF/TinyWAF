package webserver

import (
	"fmt"
	"net/http"
)

func handleHealthCheckRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}
