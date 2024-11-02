package ruleengine

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
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

func InspectRequest(r *http.Request) InspectionResult {
	inspectionId := generateInspectionId()

	for _, ruleGroup := range loadedRules.RequestRules {
		for _, rule := range ruleGroup.Rules {
			if matchesRule(r, ruleGroup.Group, rule) {
				return InspectionResult{
					InspectionId:     inspectionId,
					TriggerdByRuleId: fmt.Sprintf("%s:%s", ruleGroup.Group, rule.Id),
					RequestorIp:      r.RemoteAddr,
					ShouldBlock:      rule.Action == config.RuleActionBlock && !loadedCfg.RuleFiles.WarnInsteadOfBlock,
					ShouldWarn:       rule.Action == config.RuleActionWarn || loadedCfg.RuleFiles.WarnInsteadOfBlock,
				}
			}
		}
	}

	if loadedCfg.RequestMemory.Enabled {
		// If we got this far the request is not blocked. Continue checking rate limits if enabled
		_, ok := GetRememberedRequestsForIp(r.RemoteAddr)
		if ok {
			// @todo: apply rate limiting
			// @todo: check for enumeration attack
		}
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
	log.Printf("Evaluating rule '%s:%s'...", ruleGroupName, rule.Id)

	// If the method doesn't match, don't bother doing anything else
	if len(rule.WhenMethods) > 0 && !slices.Contains(rule.WhenMethods, strings.ToLower(r.Method)) {
		return false
	}

	// Loop over the possible inspection values (eg. URL, headers, body)
	for _, inspect := range rule.Inspect {
		switch inspect {
		case config.RuleInspectUrl:
			if runOperators(r.RequestURI, rule.Operators) {
				return true
			}

		case config.RuleInspectHeaders:
			for _, field := range rule.Fields {
				header := r.Header.Get(field)
				if runOperators(header, rule.Operators) {
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
			if runOperators(r.RemoteAddr, rule.Operators) {
				return true
			}

		case config.RuleInspectBody:
			body := []byte{}
			_, err := r.Body.Read(body)
			if err != nil {
				log.Printf("Failed to read request body: %v", err.Error())
				return false
			}

			// @todo: allow checking specific request fields - how? json object dot notation?
			if runOperators(string(body), rule.Operators) {
				return true
			}
		}
	}

	return false
}

func runOperators(field string, operator config.Operators) bool {
	field = strings.ToLower(field)

	if operator.Contains != "" {
		parts := strings.Split(strings.ToLower(operator.Contains), "|")

		for _, part := range parts {
			log.Println(strings.Contains(field, part))
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
			log.Printf("Failed to parse regex: %v", err.Error())
		}
		return matched
	}

	if operator.NotRegex != "" {
		matched, err := regexp.Match(operator.NotRegex, []byte(field))
		if err != nil {
			log.Printf("Failed to parse regex: %v", err.Error())
		}
		return !matched
	}

	log.Println("Rule has no operators!")
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
