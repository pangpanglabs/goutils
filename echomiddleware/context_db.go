package echomiddleware

import (
	"context"
	"log"
	"net/http"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/ctxdb"
	"github.com/pangpanglabs/goutils/kafka"
)

var ContextDBName = "DB"

func ContextDB(service string, xormEngine *xorm.Engine, kafkaConfig kafka.Config) echo.MiddlewareFunc {
	return ContextDBWithName(service, ContextDBName, xormEngine, kafkaConfig)
}
func ContextDBWithName(service string, contexDBName string, xormEngine *xorm.Engine, kafkaConfig kafka.Config) echo.MiddlewareFunc {
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
