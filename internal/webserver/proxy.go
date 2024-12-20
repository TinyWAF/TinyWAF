package webserver

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/TinyWAF/TinyWAF/internal/logger"
	"github.com/TinyWAF/TinyWAF/internal/ruleengine"
	"github.com/TinyWAF/TinyWAF/internal/telemetry"
)

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy, targetUrl *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		telemetry.AddRequest()
		inspectionId := ruleengine.GenerateInspectionId()

		// Pass the inspection ID upstream and make it available for use in proxy error handler
		r.Header.Add(wafInspectionIdHeaderName, inspectionId)

		logger.Debug("%v :: Request access: %v from %v", inspectionId, r.URL, r.RemoteAddr)

		inspection := ruleengine.InspectRequest(r, inspectionId)
		if inspection.ShouldBlock {
			telemetry.AddBlocked()
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

		proxy.ModifyResponse = func(res *http.Response) error {
			res.Header.Add(wafInspectionIdHeaderName, inspection.InspectionId)

			for _, headerToRemove := range loadedCfg.Listen.StripResponseHeaders {
				res.Header.Del(headerToRemove)
			}

			return nil
		}

		// If we got this far, the request is allowed to continue upstream
		logger.Debug("%v :: Pass request: %v from %v", inspectionId, r.URL, r.RemoteAddr)

		// Note that ServeHttp is non blocking and uses a goroutine under the hood
		proxy.ServeHTTP(w, r)
	}
}
