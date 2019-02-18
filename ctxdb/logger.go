package ctxdb

import (
	"context"
	"time"

	"github.com/go-xorm/core"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/ctxbase"
	"github.com/pangpanglabs/goutils/kafka"
)

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
		if cb := ctxbase.FromCtx(ctx); cb != nil {
			log.ActionID = cb.ActionID
			log.RequestID = cb.RequestID
		} else if cl := behaviorlog.FromCtx(ctx); cl != nil {
			log.ActionID = cl.ActionID
			log.RequestID = cl.RequestID
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
