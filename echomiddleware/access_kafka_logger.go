package echomiddleware

import (
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/kafka"
)

type TeeReadCloser struct {
	io.Reader
}

func (r TeeReadCloser) Close() error {
	return nil
}

func AccessLogger(serviceName string, config KafkaConfig) echo.MiddlewareFunc {
	if len(config.Brokers) == 0 {
		return nopMiddleware
	}
	producer, err := kafka.NewProducer(config.Brokers, config.Topic, func(c *sarama.Config) {
		c.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
		c.Producer.Compression = sarama.CompressionSnappy   // Compress messages
		c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	})
	if err != nil {
		log.Println(err)
		return nopMiddleware
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()

			var buf bytes.Buffer
			tee := io.TeeReader(req.Body, &buf)
			req.Body.Close()
			req.Body = TeeReadCloser{tee}

			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			request_id := req.Header.Get(echo.HeaderXRequestID)
			if request_id == "" {
				request_id = res.Header().Get(echo.HeaderXRequestID)
			}
			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			request_length, _ := strconv.ParseInt(req.Header.Get(echo.HeaderContentLength), 10, 64)

			params := map[string]interface{}{}
			for k, v := range c.QueryParams() {
				params[k] = v[0]
			}
			for _, name := range c.ParamNames() {
				params[name] = c.Param(name)
			}

			var handlerName string
			for _, r := range c.Echo().Routes() {
				if r.Path == c.Path() && r.Method == c.Request().Method {
					handlerName = r.Name
				}
			}
			handlerSplitIndex := strings.LastIndex(handlerName, ".")

			msg := map[string]interface{}{
				"service":        serviceName,
				"timestamp":      start.UTC().Format(time.RFC3339),
				"request_id":     request_id,
				"remote_ip":      c.RealIP(),
				"host":           req.Host,
				"uri":            req.RequestURI,
				"method":         req.Method,
				"path":           path,
				"referer":        req.Referer(),
				"user_agent":     req.UserAgent(),
				"status":         res.Status,
				"latency":        int64(stop.Sub(start)),
				"request_length": request_length,
				"bytes_sent":     res.Size,
				"params":         params,
				"controller":     handlerName[:handlerSplitIndex],
				"action":         handlerName[handlerSplitIndex+1:],
				"body":           buf.String(),
			}
			producer.Send(msg)
			return
		}
	}
}
