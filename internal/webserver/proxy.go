package webserver

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/TinyWAF/TinyWAF/internal/logger"
	"github.com/TinyWAF/TinyWAF/internal/ruleengine"
)

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy, targetUrl *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		inspectionId := ruleengine.GenerateInspectionId()

		// Pass the inspection ID upstream and make it available for use in proxy error handler
		r.Header.Add(wafInspectionIdHeaderName, inspectionId)

		logger.Debug("%v :: Request access: %v from %v", inspectionId, r.URL, r.RemoteAddr)

		if loadedCfg.RequestMemory.Enabled {
			ruleengine.RememberRequest(r)
		}

		inspection := ruleengine.InspectRequest(r, inspectionId)
		if inspection.ShouldBlock {
			logger.Block(
				"%v :: Denied request from IP '%v', rule '%v'",
				inspection.InspectionId,
				inspection.RequestorIp,
				inspection.TriggerdByRuleId,
			)
			respondBlocked(inspection, w)
			return
		}

		if inspection.ShouldWarn {
			logger.Warn(
				"%v :: Bypass denied request from IP '%v', rule '%v'",
				inspection.InspectionId,
				inspection.RequestorIp,
				inspection.TriggerdByRuleId,
			)
		}

		if inspection.ShouldRateLimit {
			// @todo: log rate limit info (depending on config)
			logger.Block(
				"%v :: RATELIMITED request from IP '%v', rule '%v'",
				inspection.InspectionId,
				inspection.RequestorIp,
				inspection.TriggerdByRuleId,
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

		if loadedCfg.RequestMemory.Enabled {
			proxy.ModifyResponse = func(res *http.Response) error {
				// @todo: run response rule analysis
				return nil
			}
		}

		// If we got this far, the request is allowed to continue upstream
		logger.Debug("%v :: Pass request: %v from %v", inspectionId, r.URL, r.RemoteAddr)

		// Note that ServeHttp is non blocking and uses a goroutine under the hood
		proxy.ServeHTTP(w, r)
	}
}
