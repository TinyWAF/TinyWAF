package webserver

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/TinyWAF/TinyWAF/internal"
	"github.com/TinyWAF/TinyWAF/internal/logger"
	"golang.org/x/crypto/acme/autocert"
)

var loadedCfg *internal.MainConfig

func Start(config *internal.MainConfig) error {
	loadedCfg = config

	localProtocol := "http://"
	if config.Listen.ForwardToLocalPort == 443 {
		localProtocol = "https://"
	}

	targetUrl, _ := url.Parse(fmt.Sprintf("%s%v:%v", localProtocol, "localhost", config.Listen.ForwardToLocalPort))

	for _, listenHost := range config.Listen.Hosts {
		var err error

		// Register the proxy handler
		go func() {
			logger.Info("Listening on non-TLS: '%v'", listenHost)
			err = http.ListenAndServe(listenHost, getProxyMux(targetUrl))
			if err != nil {
				logger.Error("Failed to start non-TLS reverse proxy for '%v': %v", listenHost, err)
			}
		}()
	}

	// Remove duplicates
	slices.Sort(loadedCfg.Listen.TlsDomains)
	tlsDomains := slices.Compact(loadedCfg.Listen.TlsDomains)

	if len(tlsDomains) > 0 {

		go func() {
			logger.Info("Listening on TLS: '%v'", tlsDomains)

			err := http.Serve(autocert.NewListener(tlsDomains...), getProxyMux(targetUrl))
			if err != nil {
				logger.Error("Failed to start TLS reverse proxy for '%v': %v", tlsDomains, err)
			}
		}()
	}

	return nil
}

func getProxyMux(targetUrl *url.URL) *http.ServeMux {
	proxy := NewProxy(targetUrl)
	mux := http.NewServeMux()

	// Return custom responses for gateway error
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		// Get the inspection ID header we set earlier when we intercepted the request
		inspectionId := r.Header.Get(wafInspectionIdHeaderName)
		logger.Error("%v :: Proxy error from upstream: %v", inspectionId, err.Error())
		respondUnavailable(w)
	}

	// Register the reverse proxy handler that runs rules and forwards requests
	mux.HandleFunc("/", ProxyRequestHandler(proxy, targetUrl))

	return mux
}
