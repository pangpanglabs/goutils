package echomiddleware

import (
	"context"

	"github.com/labstack/echo"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hillfolk/goutils/ctxaws"

)
var ContextAWSName ContextDBType = "AWS"

func ContextAws(session *session.Session) echo.MiddlewareFunc {

	ctxaws := ctxaws.New(session)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()

			session := ctxaws.NewSession(ctx)


			c.SetRequest(req.WithContext(context.WithValue(ctx, ContextAWSName, session)))

			return next(c)
		}
	}
}
