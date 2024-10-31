package webserver

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/TinyWAF/TinyWAF/internal/config"
)

func Start(config config.MainConfig) error {
	// load configurations from config file

	// Loop over listen IPs and create a server for each IP:port combination
	for _, listenIp := range config.Listen.Ips {
		for _, portConfig := range config.Listen.Ports {
			// If the target port isn't set, use the same as the source port
			targetPort := portConfig.Target
			if targetPort == 0 {
				targetPort = portConfig.Source
			}

			targetUrl, _ := url.Parse(fmt.Sprintf("%v:%v", config.Upstream.Destination, targetPort))

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
			err := http.ListenAndServe(fmt.Sprintf("%v:%v", listenIp, portConfig.Source), mux)
			if err != nil {
				return fmt.Errorf("could not start the server: %v", err)
			}
		}
	}

	return nil
}
