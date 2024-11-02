package webserver

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/TinyWAF/TinyWAF/internal"
)

var loadedCfg *internal.MainConfig

func Start(config *internal.MainConfig) error {
	loadedCfg = config

	for _, listenHost := range config.Listen.Hosts {
		var err error

		targetPort := listenHost.UpstreamPort
		if targetPort == 0 {
			// If the upsttream port isn't set, use the same as the source port
			parsedListenHost, err := url.ParseRequestURI(listenHost.Host)
			if err != nil {
				log.Printf("Failed to parse listen host '%v', skipping: %v", listenHost.Host, err.Error())
				continue
			}

			i64, err := strconv.ParseUint(parsedListenHost.Port(), 10, 0)
			if err != nil {
				log.Printf("Failed to parse port from listen host '%v', skipping: %v", listenHost.Host, err.Error())
				continue
			}
			targetPort = uint(i64)
		}

		targetUrl, _ := url.Parse(fmt.Sprintf("%v:%v", config.Listen.Upstream.Destination, targetPort))

		// Register the proxy handler
		proxy := NewProxy(targetUrl)
		mux := http.NewServeMux()

		// Return custom responses for gateway error
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error: %v", err.Error())
			respondUnavailable(w)
		}

		// Register the healthcheck endpoint if set
		if config.Listen.HealthcheckPath != "" {
			mux.HandleFunc(config.Listen.HealthcheckPath, handleHealthCheckRequest)
		}

		// Register the reverse proxy handler that runs rules and forwards requests
		mux.HandleFunc("/", ProxyRequestHandler(proxy, targetUrl))

		// Start the webserver for this IP and port combination
		if listenHost.Tls.CertificatePath != "" {
			err = http.ListenAndServeTLS(
				listenHost.Host,
				listenHost.Tls.CertificatePath,
				listenHost.Tls.KeyPath,
				mux,
			)
		} else {
			err = http.ListenAndServe(listenHost.Host, mux)
		}

		if err != nil {
			return fmt.Errorf("Failed to start web server '%v': %v", listenHost.Host, err)
		}
	}

	return nil
}
