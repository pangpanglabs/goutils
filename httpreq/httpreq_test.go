package httpreq_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pangpanglabs/goutils/behaviorlog"

	"github.com/pangpanglabs/goutils/httpreq"
	"github.com/pangpanglabs/goutils/test"
)

func TestHttpreq(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		test.Equals(t, r.Header.Get(behaviorlog.HeaderXActionID), "actionID-1")
		test.Equals(t, r.Header.Get(behaviorlog.HeaderXRequestID), "requestID-1")
		test.Equals(t, r.Header.Get("Authorization"), "Bearer token-1")
		response, _ := json.Marshal(map[string]interface{}{
			"success": true,
			"result":  1,
			"error":   nil,
		})
		fmt.Fprint(w, string(response))
		return
	}))
	defer s.Close()

	type ApiResult struct {
		Result  interface{} `json:"result"`
		Success bool        `json:"success"`
		Error   struct {
			Code    int         `json:"code,omitempty"`
			Details interface{} `json:"details,omitempty"`
			Message string      `json:"message,omitempty"`
		} `json:"error"`
	}

	t.Run("GET", func(t *testing.T) {
		var v ApiResult
		statusCode, err := httpreq.New(http.MethodGet, s.URL, nil).
			WithToken("token-1").
			WithRequestID("requestID-1").
			WithActionID("actionID-1").
			Call(&v)
		test.Ok(t, err)
		test.Equals(t, statusCode, 200)
		test.Equals(t, v.Result, float64(1))
	})
	t.Run("POST", func(t *testing.T) {
		var v ApiResult
		statusCode, err := httpreq.New(http.MethodGet, s.URL, nil).
			WithToken("token-1").
			WithRequestID("requestID-1").
			WithActionID("actionID-1").
			Call(&v)
		test.Ok(t, err)
		test.Equals(t, statusCode, 200)
		test.Equals(t, v.Result, float64(1))
	})
	t.Run("WithBehaviorLogContext", func(t *testing.T) {
		ctx := testContext()
		var v ApiResult
		statusCode, err := httpreq.New(http.MethodGet, s.URL, nil).
			WithToken("token-1").
			WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
			Call(&v)
		test.Ok(t, err)
		test.Equals(t, statusCode, 200)
		test.Equals(t, v.Result, float64(1))
	})
}

func testContext() context.Context {
	return (&behaviorlog.LogContext{
		RequestID: "requestID-1",
		ActionID:  "actionID-1",
	}).ToCtx(context.Background())
}
