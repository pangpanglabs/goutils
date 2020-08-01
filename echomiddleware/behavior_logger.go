package echomiddleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hillfolk/goutils/behaviorlog"
	"github.com/hillfolk/goutils/kafka"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	passwordRegex = regexp.MustCompile(`"(password|passwd)":(\s)*"(.*)"`)
)

func BehaviorLogger(serviceName string, config kafka.Config, options ...func(*behaviorlog.LogContext)) echo.MiddlewareFunc {
	hostname, err := os.Hostname()
	logrus.WithError(err).Error("Fail to get hostname")

	var producer *kafka.Producer
	if p, err := kafka.NewProducer(config.Brokers, config.Topic,
		kafka.WithDefault(),
		kafka.WithTLS(config.SSL)); err != nil {
		logrus.Error("Create Kafka Producer Error", err)
	} else {
		producer = p
	}

	var echoRouter echoRouter

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			behaviorLogger := behaviorlog.New(serviceName, req, behaviorlog.KafkaProducer(producer))
			if len(options) >= 0 {
				for _, option := range options {
					option(behaviorLogger)
				}
			}

			behaviorLogger.Hostname = hostname

			var body []byte
			if shouldWriteBodyLog(req, behaviorLogger) {
				body, _ = ioutil.ReadAll(req.Body)
				req.Body.Close()
				req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}

			c.SetRequest(req.WithContext(context.WithValue(req.Context(),
				behaviorlog.LogContextName, behaviorLogger,
			)))

			if err = next(c); err != nil {
				c.Error(err)
				behaviorLogger.Err = err.Error()
			}

			res := c.Response()

			behaviorLogger.Status = res.Status
			behaviorLogger.BytesSent = res.Size
			behaviorLogger.Controller, behaviorLogger.Action = echoRouter.getControllerAndAction(c)
			if body != nil {
				body := passwordRegex.ReplaceAll(body, []byte(`"$1": "*"`))
				var bodyParam interface{}
				d := json.NewDecoder(bytes.NewBuffer(body))
				d.UseNumber()
				if err := d.Decode(&bodyParam); err == nil {
					behaviorLogger.Body = bodyParam
				} else {
					behaviorLogger.Body = string(body)
				}

			}

			for _, name := range c.ParamNames() {
				behaviorLogger.Params[name] = c.Param(name)
			}
			behaviorLogger.Write()
			return
		}
	}
}

func shouldWriteBodyLog(req *http.Request, logContext *behaviorlog.LogContext) bool {
	if logContext != nil && logContext.BodyHide {
		return false
	}
	if req.Method != http.MethodPost &&
		req.Method != http.MethodPut &&
		req.Method != http.MethodPatch &&
		req.Method != http.MethodDelete {
		return false
	}

	contentType := req.Header.Get(echo.HeaderContentType)
	if !strings.HasPrefix(strings.ToLower(contentType), echo.MIMEApplicationJSON) {
		return false
	}

	return true

}

type echoRouter struct {
	once   sync.Once
	routes map[string]string
}

func (er *echoRouter) getControllerAndAction(c echo.Context) (controller, action string) {
	er.once.Do(func() { er.initialize(c) })

	if v := c.Get("controller"); v != nil {
		if controllerName, ok := v.(string); ok {
			controller = controllerName
		}
	}
	if v := c.Get("action"); v != nil {
		if actionName, ok := v.(string); ok {
			action = actionName
		}
	}

	if controller == "" || action == "" {
		handlerName := er.routes[fmt.Sprintf("%s+%s", c.Path(), c.Request().Method)]
		controller, action = er.convertHandlerNameToControllerAndAction(handlerName)
	}
	return
}

func (echoRouter) convertHandlerNameToControllerAndAction(handlerName string) (controller, action string) {
	handlerSplitIndex := strings.LastIndex(handlerName, ".")
	if handlerSplitIndex == -1 || handlerSplitIndex >= len(handlerName) {
		controller, action = "", handlerName
	} else {
		controller, action = handlerName[:handlerSplitIndex], handlerName[handlerSplitIndex+1:]
	}

	// 1. find this pattern: "(controller)"
	controller = controller[strings.Index(controller, "(")+1:]
	if index := strings.Index(controller, ")"); index > 0 {
		controller = controller[:index]
	}
	// 2. remove pointer symbol
	controller = strings.TrimPrefix(controller, "*")
	// 3. split by "/"
	if index := strings.LastIndex(controller, "/"); index > 0 {
		controller = controller[index+1:]
	}

	// remove function symbol
	action = strings.TrimRight(action, ")-fm")
	return
}

func (er *echoRouter) initialize(c echo.Context) {
	er.routes = make(map[string]string)
	for _, r := range c.Echo().Routes() {
		path := r.Path
		if len(path) == 0 || path[0] != '/' {
			path = "/" + path
		}
		er.routes[fmt.Sprintf("%s+%s", path, r.Method)] = r.Name
	}
}
