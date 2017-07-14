package echomiddleware

import (
	"log"

	"github.com/labstack/echo"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

func ZipkinTracer(c ZipkinConfig) echo.MiddlewareFunc {
	var (
		collector zipkin.Collector
		err       error
	)
	if c.Collector.Url != "" {
		collector, err = zipkin.NewHTTPCollector(c.Collector.Url)
		// defer collector.Close()
	} else if len(c.Kafka.Brokers) != 0 {
		collector, err = zipkin.NewKafkaCollector(
			c.Kafka.Brokers,
			zipkin.KafkaTopic(c.Kafka.Topic),
		)
		// defer collector.Close()
	} else {
		log.Println("Invalid Zipkin Config")
		return nopMiddleware
	}

	if err != nil {
		log.Println("Init Zipkin Collector Error.", err)
		return nopMiddleware
	}

	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, c.Addr, c.Service),
	)
	if err != nil {
		log.Println("Init Zipkin Tracer Error.", err)
		return nopMiddleware
	}

	operationName := "http"
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			wireContext, err := tracer.Extract(
				opentracing.TextMap,
				opentracing.HTTPHeadersCarrier(req.Header),
			)
			if err != nil && err != opentracing.ErrSpanContextNotFound {
				log.Println("Extract Tracer Error.", err)
			}
			span := tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
			defer span.Finish()

			ext.HTTPMethod.Set(span, req.Method)
			ext.HTTPUrl.Set(span, req.URL.String())
			ext.SpanKindRPCServer.Set(span)

			c.SetRequest(req.WithContext(opentracing.ContextWithSpan(req.Context(), span)))

			return next(c)
		}
	}
}
