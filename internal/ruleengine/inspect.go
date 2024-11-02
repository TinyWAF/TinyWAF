package ruleengine

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"strings"

	"github.com/TinyWAF/TinyWAF/internal/config"
)

type InspectionResult struct {
	InspectionId     string
	TriggerdByRuleId string
	RequestorIp      string
	ShouldBlock      bool
	ShouldRateLimit  bool
	ShouldWarn       bool
}

type Test struct {
	Example []string
}

func InspectRequest(r *http.Request) InspectionResult {
	inspectionId := generateInspectionId()

	for _, ruleGroup := range loadedRules.RequestRules {
		for _, rule := range ruleGroup.Rules {
			if matchesRule(r, ruleGroup.Group, rule) {
				return InspectionResult{
					InspectionId:     inspectionId,
					TriggerdByRuleId: fmt.Sprintf("%s:%s", ruleGroup.Group, rule.Id),
					RequestorIp:      r.RemoteAddr,
					ShouldBlock:      rule.Action == config.RuleActionBlock,
					ShouldWarn:       rule.Action == config.RuleActionWarn,
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

func matchesRule(r *http.Request, ruleGroupName string, rule config.Rule) bool {
	match := false

	log.Printf("Evaluating rule '%s:%s'...", ruleGroupName, rule.Id)

	// If the method doesn't match, don't bother doing anything else
	if len(rule.WhenMethods) == 0 || !slices.Contains(rule.WhenMethods, strings.ToLower(r.Method)) {
		return false
	}

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
			for _, field := range rule.Fields {
				header := r.Header.Get(field)
				for operatorKey, operatorValue := range rule.Operators {
					if runOperator(header, operatorKey, operatorValue) {
						return true
					}
				}
			}

		case config.RuleInspectBody:
		}
	}

	return match
}

func runOperator(field string, operatorKey string, operatorValue string) bool {
	field = strings.ToLower(field)
	operatorKey = strings.ToLower(operatorKey)
	operatorValue = strings.ToLower(operatorValue)

	switch operatorKey {
	case config.RuleOperatorContains:
		return strings.Contains(field, operatorValue)

	case config.RuleOperatorNotContains:
		return !strings.Contains(field, operatorValue)

	case config.RuleOperatorExactly:
		return field == operatorValue

	case config.RuleOperatorNotExactly:
		return field != operatorValue

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
