package webserver

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/TinyWAF/TinyWAF/internal/config"
)

var loadedCfg *config.MainConfig

func Start(config *config.MainConfig) error {
	loadedCfg = config

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
			var err error
			if portConfig.Tls || portConfig.Source == 443 {
				err = http.ListenAndServeTLS(
					fmt.Sprintf("%v:%v", listenIp, portConfig.Source),
					loadedCfg.Listen.Tls.CertificatePath,
					loadedCfg.Listen.Tls.KeyPath,
					mux,
				)
			} else {
				err = http.ListenAndServe(
					fmt.Sprintf("%v:%v", listenIp, portConfig.Source),
					mux,
				)
			}

			if err != nil {
				return fmt.Errorf("Failed to start web server '%v:%v': %v", listenIp, portConfig.Source, err)
			}
		}
	}

	return nil
}
