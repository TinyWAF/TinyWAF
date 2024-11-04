package ruleengine

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/TinyWAF/TinyWAF/internal/config"
	"github.com/TinyWAF/TinyWAF/internal/logger"
)

type InspectionResult struct {
	InspectionId     string
	TriggerdByRuleId string
	RequestorIp      string
	ShouldBlock      bool
	ShouldWarn       bool
}

func InspectRequest(r *http.Request, inspectionId string) InspectionResult {
	for _, ruleGroup := range loadedRules.RequestRules {
		for _, rule := range ruleGroup.Rules {
			if matchesRule(r, ruleGroup.Group, rule, inspectionId) {
				return InspectionResult{
					InspectionId:     inspectionId,
					TriggerdByRuleId: fmt.Sprintf("%s:%s", ruleGroup.Group, rule.Id),
					RequestorIp:      r.RemoteAddr,
					ShouldBlock:      rule.Action == config.RuleActionBlock && !loadedCfg.Rulesets.InspectOnly,
					ShouldWarn:       rule.Action == config.RuleActionWarn || loadedCfg.Rulesets.InspectOnly,
				}
			}
		}
	}

	// Allow this request
	return InspectionResult{
		InspectionId: inspectionId,
		RequestorIp:  r.RemoteAddr,
	}
}

func matchesRule(r *http.Request, ruleGroupName string, rule config.Rule, inspectionId string) bool {
	fqRuleName := fmt.Sprintf("%s:%s", ruleGroupName, rule.Id)
	logger.Debug("%v :: Evaluating rule '%s'...", inspectionId, fqRuleName)

	// If the method doesn't match, don't bother doing anything else
	if len(rule.WhenMethods) > 0 && !slices.Contains(rule.WhenMethods, strings.ToLower(r.Method)) {
		return false
	}

	// Loop over the possible inspection values (eg. URL, headers, body)
	for _, inspect := range rule.Inspect {
		switch inspect {
		case config.RuleInspectUrl:
			if runOperators(r.RequestURI, rule.Operators, fqRuleName, inspectionId) {
				return true
			}

		case config.RuleInspectHeaders:
			for _, field := range rule.Fields {
				header := r.Header.Get(field)
				if runOperators(header, rule.Operators, fqRuleName, inspectionId) {
					return true
				}
			}

		// @todo: figure out how to make this work nicely
		// @todo: marshal the cookies array to json to apply operators?
		// case config.RuleInspectCookies:
		// 	for _, field := range rule.Fields {
		// 		cookies := r.Cookies()
		// 		if runOperators(cookie, rule.Operators) {
		// 			return true
		// 		}
		// 	}

		case config.RuleInspectIp:
			if runOperators(r.RemoteAddr, rule.Operators, fqRuleName, inspectionId) {
				return true
			}

		case config.RuleInspectBody:
			body := []byte{}
			_, err := r.Body.Read(body)
			if err != nil {
				logger.Debug("%v :: Failed to read request body: %v", inspectionId, err.Error())
				return false
			}

			// @todo: allow checking specific request fields - how? json object dot notation?
			if runOperators(string(body), rule.Operators, fqRuleName, inspectionId) {
				return true
			}
		}
	}

	return false
}

func runOperators(field string, operator config.Operators, fqRuleName string, inspectionId string) bool {
	field = strings.ToLower(field)

	if operator.Contains != "" {
		parts := strings.Split(strings.ToLower(operator.Contains), "|")

		for _, part := range parts {
			if strings.Contains(field, part) {
				return true
			}
		}

		return false
	}

	if operator.NotContains != "" {
		match := false
		parts := strings.Split(strings.ToLower(operator.NotContains), "|")
		for _, part := range parts {
			match = strings.Contains(field, part)
		}
		return !match
	}

	if operator.Exactly != "" {
		return field == operator.Exactly
	}

	if operator.NotExactly != "" {
		return field != operator.NotExactly
	}

	if operator.Regex != "" {
		matched, err := regexp.Match(operator.Regex, []byte(field))
		if err != nil {
			logger.Error("%s :: Failed to parse regex in rule '%s': %s", inspectionId, fqRuleName, err.Error())
		}
		return matched
	}

	if operator.NotRegex != "" {
		matched, err := regexp.Match(operator.NotRegex, []byte(field))
		if err != nil {
			logger.Error("%s :: Failed to parse regex in rule '%s': %s", inspectionId, fqRuleName, err.Error())
		}
		return !matched
	}

	logger.Error("%s :: Rule '%s' has no operators!", inspectionId, fqRuleName)
	return false
}

func GenerateInspectionId() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 16)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
