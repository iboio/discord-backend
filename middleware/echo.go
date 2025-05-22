package middleware

import (
	stdCtx "context"
	"discord/context"
	"discord/query"
	"github.com/labstack/echo/v4"
)

func ContextMiddleware(appCtx context.AppContext) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request().WithContext(stdCtx.WithValue(c.Request().Context(), "appctx", appCtx))
			c.SetRequest(req)
			return next(c)
		}
	}
}

func QueryMiddleware(appCtx context.AppContext) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			q := query.NewQuery(appCtx)
			req := c.Request().WithContext(stdCtx.WithValue(c.Request().Context(), "query", q))
			c.SetRequest(req)
			return next(c)
		}
	}
}
