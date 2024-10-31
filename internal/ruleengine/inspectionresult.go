package ruleengine

type InspectionResult struct {
	InspectionId    string
	ShouldBlock     bool
	ShouldRateLimit bool
}
