package echo

import (
	"discord/context"
	"discord/internal/api"
	mw "discord/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func StartServer(appCtx context.AppContext) error {
	e := echo.New()
	e.Use(
		middleware.CORSWithConfig(
			middleware.CORSConfig{
				AllowOrigins: []string{"http://localhost:3001", "*"}, // Allow specific origins
				AllowMethods: []string{
					http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions,
				},
				AllowHeaders:     []string{"Content-Type", "Authorization"},
				AllowCredentials: true,
			}))
	e.Use(mw.ContextMiddleware(appCtx))
	e.Use(mw.QueryMiddleware(appCtx))

	e.GET("/dummy/guild", api.DummyDataGeneratorForGuild)
	e.POST("api/test/:guildId", api.Test)
	e.GET("api/guild", api.GetGuildsHandler)
	e.GET("api/guild/:guildId", api.GetGuildOverview)
	return e.Start(":8080")
}
