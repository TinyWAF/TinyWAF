package ruleengine

import (
	"net/http"
	"sync"
	"time"
)

var ipRequestHistory *sync.Map

type rememberedRequest struct {
	timestamp int64
	request   *http.Request
}

func init() {
	ipRequestHistory = &sync.Map{}
}

// Save up to 10 previous requests to memory per source IP. Use some kind of
// limited size array that pushes old items out when new ones are added
func RememberRequest(r *http.Request) {
	var requestsForIp []rememberedRequest

	data, ok := ipRequestHistory.Load(r.RemoteAddr)
	if ok {
		// Type assertion
		requestsForIp = data.([]rememberedRequest)
	} else {
		// Create an empty slice if this is the first request for this IP
		requestsForIp = []rememberedRequest{}
	}

	requestsForIp = append(requestsForIp, rememberedRequest{
		timestamp: time.Now().Unix(),
		request:   r,
	})

	ipRequestHistory.Store(r.RemoteAddr, requestsForIp)
}
}
