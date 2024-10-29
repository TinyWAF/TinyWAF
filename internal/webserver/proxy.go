package webserver

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy, targetUrl *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[ TinyWAF ] Request received at %s at %s\n", r.URL, time.Now().UTC())

		log.Println("Inbound request headers", r.Host)

		// @todo: check how this works with SSL
		// Update the headers to allow for SSL redirection
		// r.URL.Host = targetUrl.Host
		// r.URL.Scheme = targetUrl.Scheme
		// r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		// r.Host = targetUrl.Host

		log.Println("Request headers")
		for k := range r.Header {
			log.Println(k, r.Header.Get(k))
		}

		// @todo: apply firewall rules here (inc. rate limiting)
		// @todo: if request fails firewall rules, use `w` to write failed response and return
		// @todo: save rate limit data in memory but persist to DB periodically to survive restarts

		fmt.Printf("[ TinyWAF ] Forwarding request to %s at %s\n", r.URL, time.Now().UTC())

		// Note that ServeHttp is non blocking and uses a go routine under the hood
		proxy.ServeHTTP(w, r)
	}
}
