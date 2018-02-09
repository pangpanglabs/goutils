package httpreq_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

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

type Fruit struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Price     int64     `json:"price"`
	StoreCode string    `json:"store_code"`
}

func Test_Json_Get(t *testing.T) {
	type ArrayResult struct {
		TotalCount int64   `json:"totalCount"`
		Items      []Fruit `json:"items"`
	}
	var result struct {
		Result  ArrayResult `json:"result"`
		Success bool        `json:"success"`
		Error   interface{} `json:"error"`
	}
	baseUrl := "https://staging.p2shop.cn/fruit/v1/fruits"
	status, err := httpreq.New(http.MethodGet, baseUrl, nil).Call(&result)
	fmt.Println(status, result, err)
	test.Ok(t, err)
}

func Test_Xml_Put(t *testing.T) {
	xmlStr := `<xml>
	<price>34</price>
	</xml>`
	baseUrl := "https://staging.p2shop.cn/fruit/v1/fruits/14"
	status, err := httpreq.New(http.MethodPut, baseUrl, xmlStr, 2).WithContentType(httpreq.MIMEApplicationXMLCharsetUTF8).Call(nil)
	fmt.Println(err)
	test.Equals(t, http.StatusNoContent, status)
	test.Ok(t, err)
}

func Test_Transport(t *testing.T) {
	certPathName := fmt.Sprintf("c:/cert/%v/apiclient_cert.pem", os.Getenv("WXPAY_MCHID"))
	certPathKey := fmt.Sprintf("c:/cert/%v/apiclient_key.pem", os.Getenv("WXPAY_MCHID"))
	rootCa := fmt.Sprintf("c:/cert/%v/rootca.pem", os.Getenv("WXPAY_MCHID"))

	tport, err := httpreq.CertTransport(certPathName, certPathKey, rootCa)
	fmt.Println(tport, err)
	baseUrl := "https://api.mch.weixin.qq.com/secapi/pay/refund"
	xmlStr := `<xml>
	<out_refund_no>15802602088494275784251559636</out_refund_no>
	<total_fee>1</total_fee>
	<refund_fee>1</refund_fee>
	<out_trade_no>14201711085205823413229775520</out_trade_no>
	<sign>DB4C6EBEDF63884C272752476574B50B</sign>
	<appid>wx856df5e42a345096</appid>
	<mch_id>1293941701</mch_id>
	<nonce_str>2820116027603502902</nonce_str>
   </xml>`

	var respRefundDto struct {
		ReturnCode string `xml:"return_code,omitempty"`
		ReturnMsg  string `xml:"return_msg,omitempty"`
		AppId      string `xml:"appid,omitempty"`
		SubAppId   string `xml:"sub_appid,omitempty"`
		MchId      string `xml:"mch_id,omitempty"`

		SubMchId   string `xml:"sub_mch_id,omitempty"`
		NonceStr   string `xml:"nonce_str,omitempty"`
		Sign       string `xml:"sign,omitempty"`
		ResultCode string `xml:"result_code,omitempty"`
		ErrCode    string `xml:"err_code,omitempty"`

		ErrCodeDes    string `xml:"err_code_des,omitempty"`
		DeviceInfo    string `xml:"device_info,omitempty"`
		TransactionId string `xml:"transaction_id,omitempty"`
		OutRefundNo   string `xml:"out_refund_no,omitempty"`
		RefundId      string `xml:"refund_id,omitempty"`

		RefundFee            int64  `xml:"refund_fee,omitempty"`            //int64
		SettlementRefund_Fee int64  `xml:"settlement_refund_fee,omitempty"` //int64
		TotalFee             int64  `xml:"total_fee,omitempty"`             //int64
		SettlementTotalFee   int64  `xml:"settlement_total_fee,omitempty"`  //int64
		FeeType              string `xml:"fee_type,omitempty"`

		OutTradeNo string `xml:"out_trade_no,omitempty"`
	}
	status, err := httpreq.New(http.MethodPost, baseUrl, xmlStr, 2).
		WithContentType(httpreq.MIMEApplicationXMLCharsetUTF8).
		CallWithTransport(&respRefundDto, tport)
	fmt.Println(status, respRefundDto, err)
	test.Ok(t, err)
	test.Equals(t, http.StatusOK, status)
}
