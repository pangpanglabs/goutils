# goutils/httpreq

## Getting Started

```golang
var v struct {
	Result  interface{} `json:"result"`
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
}
statusCode, err := httpreq.New(http.MethodGet, "http://127.0.0.1", nil).
	Call(&v)
fmt.Println(statusCode, err,v)
```

## Basic usage

```golang
var v struct {
	Result  interface{} `json:"result"`
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
}
statusCode, err := httpreq.New(http.MethodGet, "http://127.0.0.1", nil).
	WithToken("token-1").
	WithRequestID("requestID-1").
	WithActionID("actionID-1").
	Call(&v)
```

```golang
body :=`{
	"price":12
}`
statusCode, err := httpreq.New(http.MethodPost, "http://127.0.0.1", body).
	WithToken("token-1").
	WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
	Call(&v)
```

```golang
statusCode, err := httpreq.New(http.MethodGet, "http://127.0.0.1", nil).
	Call(&v)
```

```golang
body :=`{
	"price":12
}`
statusCode, err := httpreq.New(http.MethodPost, "http://127.0.0.1", body).
	Call(&v)
```

## Advanced usage

The types of request and response support are as follows:
- JsonType
- FormType
- XmlType
- ByteArrayType

```golang
var v struct {
	Result  interface{} `json:"result"`
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
}
statusCode, err := httpreq.New(http.MethodGet, "http://127.0.0.1", nil, func(httpReq *httpreq.HttpReq) error {
		httpReq.RespDataType = httpreq.JsonType
		return nil
	}).
	Call(&v)
```

```golang
body :=`{
	"price":12
}`
statusCode, err := httpreq.New(http.MethodPost, "http://127.0.0.1", body, func(httpReq *httpreq.HttpReq) error {
		httpReq.ReqDataType = httpreq.FormType
		httpReq.RespDataType = httpreq.XmlType
		return nil
	}).
	Call(&v)
```