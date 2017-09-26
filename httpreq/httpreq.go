package httpreq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pangpanglabs/goutils/behaviorlog"
)

const DefaultMaxIdleConnsPerHost = 100

var defaultClient *http.Client

func init() {
	defaultTransportPointer, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer

	// http.DefaultMaxIdleConnsPerHost = 2
	defaultTransport.MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost

	defaultClient = &http.Client{Transport: &defaultTransport}
}

type HttpReq struct {
	req *http.Request
	err error
}

type HttpRespError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *HttpRespError) Error() string {
	return fmt.Sprint(e.Status, e.Body)
}

func New(method, url string, param interface{}) *HttpReq {
	var body io.Reader
	if param != nil {
		b, err := json.Marshal(param)
		if err != nil {
			return &HttpReq{err: err}
		}
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return &HttpReq{err: err}
	}
	req.Header.Set("Content-Type", "application/json")
	return &HttpReq{
		req: req,
	}
}

func (r *HttpReq) WithToken(token string) *HttpReq {
	if r.err != nil {
		return r
	}

	if token != "" {
		r.req.Header.Set("Authorization", "Bearer "+token)
	}

	return r
}

func (r *HttpReq) WithRequestID(requestID string) *HttpReq {
	if r.err != nil {
		return r
	}

	if requestID != "" {
		r.req.Header.Set(behaviorlog.HeaderXRequestID, requestID)
	}

	return r
}
func (r *HttpReq) WithActionID(actionID string) *HttpReq {
	if r.err != nil {
		return r
	}

	if actionID != "" {
		r.req.Header.Set(behaviorlog.HeaderXActionID, actionID)
	}

	return r
}
func (r *HttpReq) WithBehaviorLogContext(logContext *behaviorlog.LogContext) *HttpReq {
	if r.err != nil {
		return r
	}

	if logContext == nil {
		return r
	}

	r = r.WithRequestID(logContext.RequestID)
	r = r.WithActionID(logContext.ActionID)

	return r
}
func (r *HttpReq) Call(v interface{}) (int, error) {
	if r.err != nil {
		return 0, r.err
	}

	resp, err := defaultClient.Do(r.req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if v != nil {
		if err := json.Unmarshal(b, &v); err != nil {
			return resp.StatusCode, errors.New(string(b))
		}
	}
	return resp.StatusCode, nil
}
