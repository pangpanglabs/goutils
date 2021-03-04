package httpreq_test

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
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
			WithUserAgent("test user agent").
			WithCookie(map[string]string{"a": "b"}).
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

func TestParamJson(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["maxResultCount"]
		if !ok || len(keys) < 1 {
			log.Println("Url Param 'maxResultCount' is missing")
			return
		}
		test.Equals(t, "2", keys[0])
		response, _ := json.Marshal(map[string]interface{}{
			"success": true,
			"result":  keys[0],
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
		url := s.URL + "?maxResultCount=2"
		statusCode, err := httpreq.New(http.MethodGet, url, nil).
			Call(&v)
		test.Ok(t, err)
		test.Equals(t, statusCode, 200)
		test.Equals(t, v.Result, "2")
	})

}

func TestParamXml(t *testing.T) {
	type Fruit struct {
		Price int64 `xml:"price"`
	}
	type ApiResult struct {
		XMLName xml.Name `xml:"xml"`
		Result  int64    `xml:"result"`
		Success bool     `xml:"success"`
		Error   struct {
			Code    int         `xml:"code,omitempty"`
			Details interface{} `xml:"details,omitempty"`
			Message string      `xml:"message,omitempty"`
		} `xml:"error"`
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			fmt.Fprintf(w, returnXmlError(err.Error()))
			return
		}
		var f Fruit
		err = xml.Unmarshal(b, &f)
		if err != nil {
			fmt.Fprintf(w, returnXmlError(err.Error()))
			return
		}
		test.Equals(t, int64(34), f.Price)
		response, _ := xml.Marshal(&ApiResult{
			Success: true,
			Result:  f.Price,
			Error: struct {
				Code    int         `xml:"code,omitempty"`
				Details interface{} `xml:"details,omitempty"`
				Message string      `xml:"message,omitempty"`
			}{0, nil, ""},
		})
		fmt.Fprintf(w, string(response))
		return
	}))
	defer s.Close()

	t.Run("POST", func(t *testing.T) {
		xmlStr := `<xml>
		<price>34</price>
		</xml>`
		var v ApiResult
		statusCode, err := httpreq.New(http.MethodPost, s.URL, xmlStr, func(httpReq *httpreq.HttpReq) error {
			httpReq.ReqDataType = httpreq.XmlType
			httpReq.RespDataType = httpreq.XmlType
			return nil
		}).Call(&v)
		test.Ok(t, err)
		test.Equals(t, statusCode, 200)
		test.Equals(t, int64(34), v.Result)
	})

}

func TestRequestXml_ResponseByteArray(t *testing.T) {
	type ApiResult struct {
		XMLName xml.Name    `xml:"xml"`
		Result  int64       `xml:"result"`
		Success bool        `xml:"success"`
		Error   interface{} `xml:"error"`
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<xml><result>34</result><success>true</success></xml>")
		return
	}))
	defer s.Close()

	t.Run("POST", func(t *testing.T) {
		var v []byte
		statusCode, err := httpreq.New(http.MethodPost, s.URL, nil, func(httpReq *httpreq.HttpReq) error {
			httpReq.ReqDataType = httpreq.XmlType
			httpReq.RespDataType = httpreq.ByteArrayType
			return nil
		}).
			Call(&v)
		test.Ok(t, err)
		test.Equals(t, statusCode, 200)
		test.Equals(t, "<xml><result>34</result><success>true</success></xml>", string(v))
	})
}

func testContext() context.Context {
	return (&behaviorlog.LogContext{
		RequestID: "requestID-1",
		ActionID:  "actionID-1",
	}).ToCtx(context.Background())
}

func returnXmlError(errStr string) string {
	type ApiResult struct {
		XMLName xml.Name    `xml:"xml"`
		Result  int64       `xml:"result"`
		Success bool        `xml:"success"`
		Error   interface{} `xml:"error"`
	}
	errMsg, _ := xml.Marshal(&ApiResult{
		Success: false,
		Result:  0,
		Error: struct {
			Code    int         `xml:"code,omitempty"`
			Details interface{} `xml:"details,omitempty"`
			Message string      `xml:"message,omitempty"`
		}{0, nil, errStr},
	})
	return string(errMsg)
}

func TestRawResponse(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["maxResultCount"]
		if !ok || len(keys) < 1 {
			log.Println("Url Param 'maxResultCount' is missing")
			return
		}
		test.Equals(t, "2", keys[0])
		response, _ := json.Marshal(map[string]interface{}{
			"success": true,
			"result":  keys[0],
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
		url := s.URL + "?maxResultCount=2"
		resp, err := httpreq.New(http.MethodGet, url, nil).
			RawCall()
		defer resp.Body.Close()
		test.Ok(t, err)
		test.Equals(t, resp.StatusCode, 200)
		b, err := ioutil.ReadAll(resp.Body)
		test.Ok(t, err)
		err = json.Unmarshal(b, &v)
		test.Ok(t, err)
		test.Equals(t, "2", v.Result)
	})

}
