package echomiddleware

import (
	"github.com/labstack/echo/v4"
)

func ZipkinTracer(c ZipkinConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
