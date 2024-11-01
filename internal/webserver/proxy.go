package webserver

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/TinyWAF/TinyWAF/internal/ruleengine"
)

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy, targetUrl *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[ TinyWAF ] Request received at %s at %s\n", r.URL, time.Now().UTC())

		ruleengine.RememberRequest(r)

		inspection := ruleengine.InspectRequest(r)
		if inspection.ShouldBlock {
			// @todo: log block info (depending on config)
			log.Printf(
				"BLOCKED request from IP '%v': rule '%v', InspectionID:'%v'",
				inspection.RequestorIp,
				inspection.TriggerdByRuleId,
				inspection.InspectionId,
			)
			respondBlocked(inspection, w)
			return
		}
		if inspection.ShouldRateLimit {
			// @todo: log rate limit info (depending on config)
			log.Printf(
				"RATELIMITED request from IP '%v': rule '%v', InspectionID:'%v'",
				inspection.RequestorIp,
				inspection.TriggerdByRuleId,
				inspection.InspectionId,
			)
			respondRateLimited(inspection, w)
			return
		}

		// @todo: check how this works with SSL
		// Update the headers to allow for SSL redirection
		// r.URL.Host = targetUrl.Host
		// r.URL.Scheme = targetUrl.Scheme
		// r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		// r.Host = targetUrl.Host

		// log.Println("Request headers")
		// for k := range r.Header {
		// 	log.Println(k, r.Header.Get(k))
		// }
		proxy.ModifyResponse = func(res *http.Response) error {
			// @todo: run response rule analysis
			return nil
		}

		// If we got this far, the request is allowed to continue upstream
		fmt.Printf("[ TinyWAF ] Forwarding request to %s at %s\n", r.URL, time.Now().UTC())

		// Note that ServeHttp is non blocking and uses a goroutine under the hood
		proxy.ServeHTTP(w, r)
	}
}
