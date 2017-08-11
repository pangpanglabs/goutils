package echomiddleware

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/kafka"
	"github.com/sirupsen/logrus"
)

const BehaviorLoggerName = "BehaviorLogger"

type BehaviorLogContext struct {
	producer *kafka.Producer
	logger   *logrus.Logger

	Service       string                 `json:"service,omitempty"`
	Timestamp     string                 `json:"timestamp,omitempty"`
	RequestID     string                 `json:"request_id,omitempty"`
	RemoteIP      string                 `json:"remote_ip,omitempty"`
	Host          string                 `json:"host,omitempty"`
	Uri           string                 `json:"uri,omitempty"`
	Method        string                 `json:"method,omitempty"`
	Path          string                 `json:"path,omitempty"`
	Referer       string                 `json:"referer,omitempty"`
	UserAgent     string                 `json:"user_agent,omitempty"`
	Status        int                    `json:"status,omitempty"`
	Latency       int64                  `json:"latency,omitempty"`
	RequestLength int64                  `json:"request_length,omitempty"`
	BytesSent     int64                  `json:"bytes_sent,omitempty"`
	Params        interface{}            `json:"params,omitempty"`
	Controller    string                 `json:"controller,omitempty"`
	Action        string                 `json:"action,omitempty"`
	Body          string                 `json:"body,omitempty"`
	BizAttr       map[string]interface{} `json:"bizAttr,omitempty"`
}

func BehaviorLogger(serviceName string, config KafkaConfig) echo.MiddlewareFunc {
	var producer *kafka.Producer
	if p, err := kafka.NewProducer(config.Brokers, config.Topic, func(c *sarama.Config) {
		c.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
		c.Producer.Compression = sarama.CompressionSnappy   // Compress messages
		c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	}); err != nil {
		logrus.Error("Create Kafka Producer Error", err)
	} else {
		producer = p
	}

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			requestId := req.Header.Get(echo.HeaderXRequestID)
			if requestId == "" {
				requestId = res.Header().Get(echo.HeaderXRequestID)
			}
			controller, action := getControllerAndAction(c)
			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			behaviorLogger := &BehaviorLogContext{
				producer: producer,
				logger:   logger,

				Service:    serviceName,
				RequestID:  requestId,
				RemoteIP:   c.RealIP(),
				Host:       req.Host,
				Uri:        req.RequestURI,
				Method:     req.Method,
				Path:       path,
				Referer:    req.Referer(),
				UserAgent:  req.UserAgent(),
				Controller: controller,
				Action:     action,
				BizAttr:    map[string]interface{}{},
			}

			c.SetRequest(req.WithContext(context.WithValue(req.Context(),
				"request_id", requestId,
			)))
			c.SetRequest(req.WithContext(context.WithValue(req.Context(),
				BehaviorLoggerName, behaviorLogger,
			)))

			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			params := map[string]interface{}{}
			for k, v := range c.QueryParams() {
				params[k] = v[0]
			}
			for _, name := range c.ParamNames() {
				params[name] = c.Param(name)
			}

			behaviorLogger.Timestamp = start.UTC().Format(time.RFC3339)
			behaviorLogger.Status = res.Status
			behaviorLogger.Latency = int64(stop.Sub(start))
			behaviorLogger.RequestLength, _ = strconv.ParseInt(req.Header.Get(echo.HeaderContentLength), 10, 64)
			behaviorLogger.BytesSent = res.Size
			behaviorLogger.Params = params
			behaviorLogger.Controller = controller
			behaviorLogger.Action = action

			behaviorLogger.write()

			return
		}
	}
}

func NewNopLogger() *BehaviorLogContext {
	return &BehaviorLogContext{
		BizAttr: map[string]interface{}{},
	}
}

func (c *BehaviorLogContext) WithBizAttr(key string, value interface{}) *BehaviorLogContext {
	c.BizAttr[key] = value
	return c
}
func (c *BehaviorLogContext) WithBizAttrs(attrs map[string]interface{}) *BehaviorLogContext {
	for k, v := range attrs {
		c.BizAttr[k] = v
	}
	return c
}
func (c *BehaviorLogContext) WithCallURLInfo(method, uri string, param interface{}, responseStatus int) *BehaviorLogContext {
	c.Method = method
	c.Uri = uri
	c.Params = param
	c.Status = responseStatus
	if url, err := url.Parse(uri); err == nil {
		c.Path = url.Path
		c.Host = url.Host
	}

	return c
}
func (c *BehaviorLogContext) Log(action string) {
	c.Timestamp = time.Now().UTC().Format(time.RFC3339)
	c.Action = action
	c.write()
}
func (c *BehaviorLogContext) write() {
	if c.producer != nil {
		c.producer.Send(c)
	}
	c.logger.WithFields(logrus.Fields{
		"service":        c.Service,
		"timestamp":      c.Timestamp,
		"request_id":     c.RequestID,
		"remote_ip":      c.RemoteIP,
		"host":           c.Host,
		"uri":            c.Uri,
		"method":         c.Method,
		"path":           c.Path,
		"referer":        c.Referer,
		"user_agent":     c.UserAgent,
		"status":         c.Status,
		"latency":        c.Latency,
		"request_length": c.RequestLength,
		"bytes_sent":     c.BytesSent,
		"params":         c.Params,
		"controller":     c.Controller,
		"action":         c.Action,
		"body":           c.Body,
		"bizAttr":        c.BizAttr,
	}).Info()
}
func getControllerAndAction(c echo.Context) (controller, action string) {
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
		var handlerName string
		for _, r := range c.Echo().Routes() {
			if r.Path == c.Path() && r.Method == c.Request().Method {
				handlerName = r.Name
			}
		}
		handlerSplitIndex := strings.LastIndex(handlerName, ".")
		controller, action = handlerName[:handlerSplitIndex], handlerName[handlerSplitIndex+1:]
	}
	return
}
