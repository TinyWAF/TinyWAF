package ruleengine

import "net/http"

type InspectionResult struct {
	InspectionId     string
	TriggerdByRuleId string
	RequestorIp      string
	ShouldBlock      bool
	ShouldRateLimit  bool
}

func InspectRequest(r *http.Request) InspectionResult {
	// @todo: run rules against request
	// @todo: early return if blocked

	_, ok := GetRememberedRequestsForIp(r.RemoteAddr)
	if ok {
		// @todo: apply rate limiting
		// @todo: check for enumeration attack
	}

	return InspectionResult{
		InspectionId:     "abc123",
		TriggerdByRuleId: "no-wordpress:url-wp-admin",
		RequestorIp:      r.RemoteAddr,
		ShouldRateLimit:  true,
	}
}

func InspectResponse(r *http.Request) InspectionResult {

	// @todo
	return InspectionResult{}
}
