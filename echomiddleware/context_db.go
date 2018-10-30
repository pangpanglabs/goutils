package echomiddleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pangpanglabs/goutils/behaviorlog"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/kafka"
)

const ContextDBName = "DB"

func ContextDB(service string, db *xorm.Engine, kafkaConfig KafkaConfig) echo.MiddlewareFunc {
	db.ShowExecTime()
	if len(kafkaConfig.Brokers) != 0 {
		if producer, err := kafka.NewProducer(kafkaConfig.Brokers, kafkaConfig.Topic, func(c *sarama.Config) {
			c.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
			c.Producer.Compression = sarama.CompressionGZIP     // Compress messages
			c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

		}); err == nil {
			db.SetLogger(&dbLogger{serviceName: service, Producer: producer})
			db.ShowSQL()
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			session := db.NewSession()
			defer session.Close()

			func(session interface{}, ctx context.Context) {
				if s, ok := session.(interface{ SetContext(context.Context) }); ok {
					s.SetContext(ctx)
				}
			}(session, req.Context())

			c.SetRequest(req.WithContext(context.WithValue(req.Context(), ContextDBName, session)))

			switch req.Method {
			case "POST", "PUT", "DELETE", "PATCH":
				if err := session.Begin(); err != nil {
					log.Println(err)
				}
				if err := next(c); err != nil {
					session.Rollback()
					return err
				}
				if c.Response().Status >= 500 {
					session.Rollback()
					return nil
				}
				if err := session.Commit(); err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			default:
				return next(c)
			}

			return nil
		}
	}
}

type SqlLog struct {
	Service   string      `json:"service,omitempty"`
	RequestID string      `json:"requestId,omitempty"`
	ActionID  string      `json:"actionId,omitempty"`
	Sql       interface{} `json:"sql,omitempty"`
	Args      interface{} `json:"args,omitempty"`
	Took      interface{} `json:"took,omitempty"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
}
type dbLogger struct {
	serviceName string
	*kafka.Producer
}

func (logger *dbLogger) Write(v []interface{}) {
	if len(v) == 0 {
		return
	}
	log := SqlLog{
		Service:   logger.serviceName,
		Sql:       v[0],
		Timestamp: time.Now(),
	}
	if ctx, ok := v[len(v)-1].(context.Context); ok {
		if logContext := behaviorlog.FromCtx(ctx); logContext != nil {
			log.ActionID = logContext.ActionID
			log.RequestID = logContext.RequestID
		}
		v = v[:len(v)-1]
	}

	if len(v) == 3 {
		log.Args = v[1]
		log.Took = v[2]
	} else if len(v) == 2 {
		log.Took = v[1]
	}

	if d, ok := log.Took.(time.Duration); ok {
		log.Timestamp = log.Timestamp.Add(-d)
	}

	logger.Send(&log)
}
func (logger *dbLogger) Infof(format string, v ...interface{})  { logger.Write(v) }
func (logger *dbLogger) Errorf(format string, v ...interface{}) {}
func (logger *dbLogger) Debugf(format string, v ...interface{}) {}
func (logger *dbLogger) Warnf(format string, v ...interface{})  {}

func (logger *dbLogger) Debug(v ...interface{})   {}
func (logger *dbLogger) Error(v ...interface{})   {}
func (logger *dbLogger) Info(v ...interface{})    {}
func (logger *dbLogger) Warn(v ...interface{})    {}
func (logger *dbLogger) SetLevel(l core.LogLevel) {}
func (logger *dbLogger) ShowSQL(show ...bool)     {}
func (logger *dbLogger) Level() core.LogLevel     { return 0 }
func (logger *dbLogger) IsShowSQL() bool          { return true }
