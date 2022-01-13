package echomiddleware

import (
	"context"
	"github.com/hillfolk/goutils/ctxdb"
	"github.com/hillfolk/goutils/kafka"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"xorm.io/xorm"
)

type ContextDBType string

var ContextDBName ContextDBType = "DB"

func ContextDB(service string, xormEngine *xorm.Engine, kafkaConfig kafka.Config) echo.MiddlewareFunc {
	return ContextDBWithName(service, ContextDBName, xormEngine, kafkaConfig)
}
func ContextDBWithName(service string, contexDBName ContextDBType, xormEngine *xorm.Engine, kafkaConfig kafka.Config) echo.MiddlewareFunc {
	db := ctxdb.New(xormEngine, service, kafkaConfig)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()

			session := db.NewSession(ctx)
			defer session.Close()

			c.SetRequest(req.WithContext(context.WithValue(ctx, contexDBName, session)))

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
