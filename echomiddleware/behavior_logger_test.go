package echomiddleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/stretchr/testify/assert"
)

func TestRefactorControllerAndAction(t *testing.T) {
	datas := []struct {
		handlerName, controller, action string
	}{
		{"best/eland-show-service/controllers.(*AccountApiController).GetOpenid-fm", "AccountApiController", "GetOpenid"},
		{"best/eland-show-service/controllers.(RankingApiController).GetByType-fm", "RankingApiController", "GetByType"},
		{"main.(*ProductHandler).SearchSkus-fm", "ProductHandler", "SearchSkus"},
		{"main.(ProductHandler).SearchSkus-fm", "ProductHandler", "SearchSkus"},
		{"main.(*handler).(main.postSystmeContents)-fm", "handler", "postSystmeContents"},
		{"main.(*EventHandler).CreatePriceUpdatedEvent-fm", "EventHandler", "CreatePriceUpdatedEvent"},
		{"order-settings-service/api.CheckShopGrabForLogin", "api", "CheckShopGrabForLogin"},
		{"p2saas/smart-seller/kit.Ping", "kit", "Ping"},
		// {"main.main.func1.3.1", "main.main.func1.3", "1"},
	}
	for _, data := range datas {
		t.Run(data.handlerName, func(t *testing.T) {
			controller, action := echoRouter{}.convertHandlerNameToControllerAndAction(data.handlerName)
			assert.Equal(t, controller, data.controller)
			assert.Equal(t, action, data.action)
			fmt.Println(controller, action)
		})
	}
}

func TestAddRequestBody(t *testing.T) {
	datas := []struct{ body, passwordFieldName string }{
		{`{"a":41431341324143,"password": "123"}`, "password"},
		{`{"a":"b","passwd": "123"}`, "passwd"},
		{`{"a":"b","password":
			"123"}`, "password"},
	}

	for _, data := range datas {
		t.Run(data.body, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(echo.POST, "/post", strings.NewReader(data.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			var logContext *behaviorlog.LogContext

			h := BehaviorLogger("test", KafkaConfig{})(func(c echo.Context) error {
				logContext = behaviorlog.FromCtx(c.Request().Context())

				var v interface{}
				if err := c.Bind(&v); err != nil {
					return c.String(http.StatusBadRequest, err.Error())
				}
				return c.JSON(http.StatusOK, v)
			})
			h(c)

			assert.Equal(t, logContext.Params[data.passwordFieldName], "*")
			assert.JSONEq(t, rec.Body.String(), data.body)
			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestParams(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/post?a=query-a&b=query-b&c=query-c&d=query-d", strings.NewReader(`{"a":"body-a","b":"body-b","c":"body-c"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("a", "b")
	c.SetParamValues("path-a", "path-b")

	var logContext *behaviorlog.LogContext

	BehaviorLogger("ping", KafkaConfig{})(func(c echo.Context) error {
		logContext = behaviorlog.FromCtx(c.Request().Context())
		return c.String(http.StatusOK, "ping")
	})(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	fmt.Println("logContext.Params:", logContext.Params)

	// params priority
	// 1: path param
	// 2: request body
	// 3: query param
	assert.Equal(t, logContext.Params["a"], "path-a")
	assert.Equal(t, logContext.Params["b"], "path-b")
	assert.Equal(t, logContext.Params["c"], "body-c")
	assert.Equal(t, logContext.Params["d"], "query-d")

}
