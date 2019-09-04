package httpreq

import (
	"net/http"
	"time"
)

type ClientConfig struct {
	Timeout time.Duration
}

func NewClient(options ...ClientConfig) *http.Client {
	client := *defaultClient
	for _, option := range options {
		if option.Timeout != 0 {
			client.Timeout = option.Timeout
		}

	}

	return &client
}
