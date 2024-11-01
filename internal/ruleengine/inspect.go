package ruleengine

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/TinyWAF/TinyWAF/internal/config"
)

type InspectionResult struct {
	InspectionId     string
	TriggerdByRuleId string
	RequestorIp      string
	ShouldBlock      bool
	ShouldRateLimit  bool
}

type Test struct {
	Example []string
}

func InspectRequest(r *http.Request) InspectionResult {
	inspectionId := generateInspectionId()

	for _, ruleGroup := range loadedRules.RequestRules {
		for _, rule := range ruleGroup.Rules {
			if matchesRule(r, rule) && rule.Action == config.RuleActionBlock {
				// @todo: log warning if rule.action is warn

				return InspectionResult{
					InspectionId:     inspectionId,
					TriggerdByRuleId: fmt.Sprintf("%s:%s", ruleGroup.Group, rule.Id),
					RequestorIp:      r.RemoteAddr,
					ShouldBlock:      true,
				}
			}
		}
	}

	// If we got this far the request is not blocked. Continue checking rate limits
	_, ok := GetRememberedRequestsForIp(r.RemoteAddr)
	if ok {
		// @todo: apply rate limiting
		// @todo: check for enumeration attack
	}

	// Allow this request
	return InspectionResult{
		InspectionId:    inspectionId,
		RequestorIp:     r.RemoteAddr,
		ShouldRateLimit: false,
		ShouldBlock:     false,
	}
}

func InspectResponse(r *http.Request) InspectionResult {

	// @todo
	return InspectionResult{}
}

func matchesRule(r *http.Request, rule config.Rule) bool {
	match := false

	log.Printf("Evaluating rule %s", rule.Id)

	// Loop over the possible inspection values (eg. URL, headers, body)
	for _, inspect := range rule.Inspect {
		switch strings.ToLower(inspect) {
		case config.RuleInspectUrl:
			// Loop over the operators to run for this rule
			for operatorKey, operatorValue := range rule.Operators {
				if runOperator(r.URL.String(), operatorKey, operatorValue) {
					return true
				}
			}

		case config.RuleInspectHeaders:

		case config.RuleInspectBody:
		}
	}

	return match
}

func runOperator(field string, operatorKey string, operatorValue string) bool {
	switch strings.ToLower(operatorKey) {
	case config.RuleOperatorContains:
		return strings.Contains(field, operatorValue)

	case config.RuleOperatorNotContains:

	case config.RuleOperatorExactly:

	case config.RuleOperatorNotExactly:

	case config.RuleOperatorRegex:

	case config.RuleOperatorNotRegex:

	}

	return false
}

func generateInspectionId() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 16)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
