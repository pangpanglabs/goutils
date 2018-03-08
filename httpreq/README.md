# goutils/httpreq

## Getting Started

```golang
var v ApiResult
statusCode, err := httpreq.New(http.MethodGet, s.URL, nil).
	WithToken("token-1").
	WithRequestID("requestID-1").
	WithActionID("actionID-1").
	Call(&v)
```

```golang
var v ApiResult
statusCode, err := httpreq.New(http.MethodPost, s.URL, body).
	WithToken("token-1").
	WithBehaviorLogContext(behaviorlog.FromCtx(ctx)).
	Call(&v)
```

```golang
var v ApiResult
statusCode, err := httpreq.New(http.MethodPost, s.URL, body,httpreq.XmlType).
	Call(&v)
```

```golang
var v ApiResult
statusCode, err := httpreq.New(http.MethodPost, s.URL, body,httpreq.JsonType,httpreq.XmlType).
	Call(&v)
```