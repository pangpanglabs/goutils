package echomiddleware

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/kafka"
	"github.com/sirupsen/logrus"
)

const BehaviorLoggerName = "BehaviorLogger"

func BehaviorLogger(serviceName string, config KafkaConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			behaviorLogger := NewBehaviorLogContext(c, KafkaProducer(config.Brokers, config.Topic))
			c.SetRequest(req.WithContext(context.WithValue(req.Context(), BehaviorLoggerName, behaviorLogger)))
			return next(c)
		}
	}
}

type Attrs map[string]interface{}
type BehaviorLogContext struct {
	producer *kafka.Producer
	logger   *logrus.Logger

	Timestamp  string      `json:"timestamp"`
	Service    string      `json:"service"`
	RequestID  string      `json:"request_id"`
	RemoteIP   string      `json:"remote_ip"`
	Host       string      `json:"host"`
	Uri        string      `json:"uri"`
	Method     string      `json:"method"`
	Status     int         `json:"status"`
	Latency    string      `json:"latency"`
	Params     interface{} `json:"params"`
	Controller string      `json:"controller"`
	Action     string      `json:"action"`
	Body       string      `json:"body"`
	BizAttr    Attrs       `json:"bizAttr"`
}

func NewNopLogger() *BehaviorLogContext {
	return &BehaviorLogContext{}
}

func NewBehaviorLogContext(c echo.Context, options ...func(*BehaviorLogContext)) *BehaviorLogContext {
	req := c.Request()
	controller, action := getControllerAndAction(c)
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	logContext := &BehaviorLogContext{
		RequestID:  req.Header.Get(echo.HeaderXRequestID),
		RemoteIP:   c.RealIP(),
		Host:       req.Host,
		Uri:        req.RequestURI,
		Method:     req.Method,
		Controller: controller,
		Action:     action,
		BizAttr:    map[string]interface{}{},

		logger: logger,
	}

	for _, o := range options {
		if o != nil {
			o(logContext)
		}
	}
	return logContext
}

func KafkaProducer(brokers []string, topic string) func(*BehaviorLogContext) {
	if len(brokers) == 0 {
		return nil
	}
	kafkaProducer, err := kafka.NewProducer(brokers, topic, func(c *sarama.Config) {
		c.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
		c.Producer.Compression = sarama.CompressionSnappy   // Compress messages
		c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	})
	if err != nil {
		log.Println(err)
		return nil
	}
	return func(c *BehaviorLogContext) {
		c.producer = kafkaProducer
	}
}

func (c *BehaviorLogContext) WithBizAttr(key string, value interface{}) *BehaviorLogContext {
	c.BizAttr[key] = value
	return c
}
func (c *BehaviorLogContext) WithBizAttrs(attrs Attrs) *BehaviorLogContext {
	for k, v := range attrs {
		c.BizAttr[k] = v
	}
	return c
}
func (c *BehaviorLogContext) WithRequestInfo(method, uri string, param interface{}, responseStatus int) *BehaviorLogContext {
	c.Method = method
	c.Uri = uri
	c.Params = param
	c.Status = responseStatus

	return c
}
func (c *BehaviorLogContext) Log(action string) {
	c.Timestamp = time.Now().UTC().Format(time.RFC3339)
	c.Action = action
	if c.producer != nil {
		c.producer.Send(c)
	}
	c.logger.WithFields(logrus.Fields{
		"service":    c.Service,
		"timestamp":  c.Timestamp,
		"request_id": c.RequestID,
		"remote_ip":  c.RemoteIP,
		"host":       c.Host,
		"uri":        c.Uri,
		"method":     c.Method,
		"status":     c.Status,
		"latency":    c.Latency,
		"params":     c.Params,
		"controller": c.Controller,
		"action":     c.Action,
		"body":       c.Body,
		"bizAttr":    c.BizAttr,
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
