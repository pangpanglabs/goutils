package cronjob

import (
	"context"
	"github.com/hillfolk/goutils/ctxdb"
	"github.com/hillfolk/goutils/echomiddleware"
	"github.com/hillfolk/goutils/kafka"
	"xorm.io/xorm"
)

func ContextDB(service string, xormEngine *xorm.Engine, kafkaConfig kafka.Config) Middleware {
	ctxdb := ctxdb.New(xormEngine, service, kafkaConfig)

	return func(job, action string, next HandlerFunc) HandlerFunc {
		return func(ctx context.Context) error {
			session := ctxdb.NewSession(ctx)
			defer session.Close()

			ctx = context.WithValue(ctx, echomiddleware.ContextDBName, session)

			return next(ctx)
		}
	}
}
