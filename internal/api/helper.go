package api

import (
	"discord/context"
	"discord/query"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DummyDataGeneratorForGuild(c echo.Context) error {
	appCtx := c.Request().Context().Value("appctx")
	appCtx, ok := appCtx.(context.AppContext)
	if !ok || appCtx == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}
	q := c.Request().Context().Value("query").(*query.Query)
	err := q.DummyData(1000000)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": "Failed to generate dummy data",
			})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Dummy data generated successfully"})
}
