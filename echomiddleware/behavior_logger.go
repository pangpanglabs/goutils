package echomiddleware

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/kafka"
	"github.com/sirupsen/logrus"
)

func BehaviorLogger(serviceName string, config KafkaConfig) echo.MiddlewareFunc {
	var producer *kafka.Producer
	if p, err := kafka.NewProducer(config.Brokers, config.Topic, func(c *sarama.Config) {
		c.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
		c.Producer.Compression = sarama.CompressionGZIP     // Compress messages
		c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	}); err != nil {
		logrus.Error("Create Kafka Producer Error", err)
	} else {
		producer = p
	}

	var echoRouter echoRouter

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			// fmt.Println("req.RequestURI:", req.RequestURI)
			// fmt.Println("req.Host:", req.Host)
			behaviorLogger := behaviorlog.New(serviceName, req, behaviorlog.KafkaProducer(producer))

			c.SetRequest(req.WithContext(context.WithValue(req.Context(),
				behaviorlog.LogContextName, behaviorLogger,
			)))

			if err = next(c); err != nil {
				c.Error(err)
				behaviorLogger.Err = err
			}

			res := c.Response()

			behaviorLogger.Status = res.Status
			behaviorLogger.BytesSent = res.Size
			behaviorLogger.Controller, behaviorLogger.Action = echoRouter.getControllerAndAction(c)

			params := map[string]interface{}{}
			for k, v := range c.QueryParams() {
				params[k] = v[0]
			}
			for _, name := range c.ParamNames() {
				params[name] = c.Param(name)
			}
			behaviorLogger.Params = params

			behaviorLogger.Write()
			return
		}
	}
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
		handlerSplitIndex := strings.LastIndex(handlerName, ".")
		if handlerSplitIndex == -1 || handlerSplitIndex >= len(handlerName) {
			controller, action = "", handlerName
		} else {
			controller, action = handlerName[:handlerSplitIndex], handlerName[handlerSplitIndex+1:]
		}
	}
	return
}

func (er *echoRouter) initialize(c echo.Context) {
	er.routes = make(map[string]string)
	for _, r := range c.Echo().Routes() {
		er.routes[fmt.Sprintf("%s+%s", r.Path, r.Method)] = r.Name
	}
}
