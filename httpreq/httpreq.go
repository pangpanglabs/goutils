package httpreq

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

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
	defaultTransport.MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
	defaultClient = &http.Client{Transport: &defaultTransport}
}

type OptionFunc func(*HttpReq) error

type HttpReq struct {
	Req          *http.Request
	ReqDataType  formatType
	RespDataType formatType
	err          error
}

type HttpRespError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *HttpRespError) Error() string {
	return fmt.Sprint(e.Status, e.Body)
}

func New(method, url string, param interface{}, options ...OptionFunc) *HttpReq {
	httpReq := &HttpReq{}
	for _, option := range options {
		if option == nil {
			continue
		}
		httpReq.err = option(httpReq)
		if httpReq.err != nil {
			return httpReq
		}
	}
	var body io.Reader
	if param != nil {
		b, err := DataTypeFactory{}.New(httpReq.ReqDataType).marshal(param)
		if err != nil {
			return &HttpReq{err: err}
		}
		body = bytes.NewBuffer(b)
	}
	httpReq.Req, httpReq.err = http.NewRequest(method, url, body)
	if httpReq.err != nil {
		return httpReq
	}
	return httpReq
}

func (r *HttpReq) WithToken(token string) *HttpReq {
	if r.err != nil {
		return r
	}

	if token != "" {
		if !strings.HasPrefix(token, "Bearer ") {
			token = "Bearer " + token
		}
		r.Req.Header.Set("Authorization", token)
	}

	return r
}

func (r *HttpReq) WithRequestID(requestID string) *HttpReq {
	if r.err != nil {
		return r
	}

	if requestID != "" {
		r.Req.Header.Set(behaviorlog.HeaderXRequestID, requestID)
	}

	return r
}
func (r *HttpReq) WithActionID(actionID string) *HttpReq {
	if r.err != nil {
		return r
	}

	if actionID != "" {
		r.Req.Header.Set(behaviorlog.HeaderXActionID, actionID)
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
	if len(r.Req.Header.Get("Authorization")) == 0 {
		r = r.WithToken(logContext.AuthToken)
	}

	return r
}

func (r *HttpReq) WithUserAgent(userAgent string) *HttpReq {
	r.Req.Header.Add("User-Agent", userAgent)
	return r
}

func (r *HttpReq) WithHeader(key, value string) *HttpReq {
	r.Req.Header.Add(key, value)
	return r
}

func (r *HttpReq) WithCookie(m map[string]string) *HttpReq {
	for k, v := range m {
		r.Req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	return r
}

func (r *HttpReq) Call(v interface{}) (int, error) {
	return r.call(v, defaultClient)
}

func (r *HttpReq) CallWithClient(v interface{}, httpClient *http.Client) (int, error) {
	return r.call(v, httpClient)
}

func (r *HttpReq) CallWithTransport(v interface{}, transport *http.Transport) (int, error) {
	httpClient := &http.Client{Transport: transport}
	return r.call(v, httpClient)
}

func (r *HttpReq) call(v interface{}, httpClient *http.Client) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	if len(r.Req.Header.Get("Content-Type")) == 0 {
		r.Req.Header.Set("Content-Type", DataTypeFactory{}.New(r.ReqDataType).contentType())
	}
	resp, err := httpClient.Do(r.Req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	if v != nil {
		if err := (DataTypeFactory{}).New(r.RespDataType).unMarshal(b, v); err != nil {
			return resp.StatusCode, err
		}
	}
	return resp.StatusCode, nil

}

func (r *HttpReq) RawCall() (*http.Response, error) {
	return r.rawCall(defaultClient)
}

func (r *HttpReq) RawCallWithClient(httpClient *http.Client) (*http.Response, error) {
	return r.rawCall(httpClient)
}

func (r *HttpReq) RawCallWithTransport(transport *http.Transport) (*http.Response, error) {
	httpClient := &http.Client{Transport: transport}
	return r.rawCall(httpClient)
}

func (r *HttpReq) rawCall(httpClient *http.Client) (*http.Response, error) {
	if r.err != nil {
		return nil, r.err
	}
	if len(r.Req.Header.Get("Content-Type")) == 0 {
		r.Req.Header.Set("Content-Type", DataTypeFactory{}.New(r.ReqDataType).contentType())
	}
	resp, err := httpClient.Do(r.Req)
	if err != nil {
		return nil, err
	}
	return resp, nil

}
