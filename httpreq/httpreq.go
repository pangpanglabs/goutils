package httpreq

import (
	"bytes"
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
	defaultTransport.MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
	defaultClient = &http.Client{Transport: &defaultTransport}
}

type HttpReq struct {
	req      *http.Request
	dataType reqFormatType
	err      error
}

type HttpRespError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *HttpRespError) Error() string {
	return fmt.Sprint(e.Status, e.Body)
}

func New(method, url string, param interface{}, dataTypes ...reqFormatType) *HttpReq {
	var dataType reqFormatType
	if dataTypes == nil || len(dataTypes) == 0 {
		dataType = 0
	} else {
		dataType = dataTypes[0]
	}
	var body io.Reader
	if param != nil {
		b, err := DataTypeFactory{}.New(dataType).marshal(param)
		if err != nil {
			return &HttpReq{err: err}
		}
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return &HttpReq{err: err}
	}
	return &HttpReq{
		req:      req,
		dataType: dataType,
	}
}

func (r *HttpReq) WithContentType(contentType string) *HttpReq {
	if r.err != nil {
		return r
	}
	if r.dataType == 0 {
		if !(contentType == MIMEApplicationJSON || contentType == MIMEApplicationJSONCharsetUTF8) {
			r.err = fmt.Errorf("If the Content-Type is not json, the dataTypes parameter in the httpreq.New method is required")
			return r
		}
	}
	if contentType != "" {
		r.req.Header.Set("Content-Type", contentType)
	}

	return r
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
	return r.call(v, defaultClient)
}

func (r *HttpReq) CallWithClient(v interface{}, httpClient *http.Client) (int, error) {
	return r.call(v, httpClient)
}

func (r *HttpReq) CallWithTransport(v interface{}, transport *http.Transport) (int, error) {
	httpClient := &http.Client{Transport: transport}
	return r.call(v, httpClient)
}

func (r *HttpReq) SetGlobalTransport(v interface{}, transport *http.Transport) (int, error) {
	if defaultClient != nil {
		defaultClient.Transport = transport
	}
	return r.call(v, defaultClient)
}

func (r *HttpReq) call(v interface{}, httpClient *http.Client) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	if len(r.req.Header.Get("Content-Type")) == 0 {
		r.req.Header.Set("Content-Type", DataTypeFactory{}.New(r.dataType).head())
	}
	resp, err := httpClient.Do(r.req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	if v != nil {
		err = DataTypeFactory{}.New(r.dataType).unMarshal(b, v)
		if err != nil {
			return resp.StatusCode, errors.New(string(b))
		}
	}
	return resp.StatusCode, nil

}
