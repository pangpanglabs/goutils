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

### Timeout

1. Declare global `*http.Client` variable
```golang
var client = httpreq.NewClient(httpreq.ClientConfig{
	Timeout: time.Second,
})
```
2. Call with `*http.Client`
```golang
_, err := httpreq.New(http.MethodGet, url, nil).CallWithClient(&resp, client)
```

### Retry
Use https://github.com/matryer/try

1. Import `try` package

```golang
import try "gopkg.in/matryer/try.v1"
```

2. Use `try` package
```golang
err := try.Do(func(attempt int) (bool, error) {
	_, err := httpreq.New(http.MethodGet, url, nil).Call(&resp)
	return attempt < 5, err // try 5 times
})
```

### ResponseType

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