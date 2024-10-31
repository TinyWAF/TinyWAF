package ruleengine

import (
	"net/http"
	"sync"
)

var memory *sync.Map
var ipLastRequest *sync.Map

type rememberedRequest struct{}

func init() {
	memory = &sync.Map{}
	ipLastRequest = &sync.Map{}
}

// Save up to 10 previous requests to memory per source IP. Use some kind of
// limited size array that pushes old items out when new ones are added
func RememberRequest(r *http.Request) {
}
