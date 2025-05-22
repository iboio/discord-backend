package api

import (
	"discord/context"
	"discord/models"
	pb "discord/proto"
	"discord/query"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Guild struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}
type GuildResponse struct {
	Guilds []Guild `json:"guilds"`
}

func GetGuildsHandler(c echo.Context) error {
	val := c.Request().Context().Value("appctx")
	appCtx, ok := val.(context.AppContext)
	if !ok || appCtx == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	resp, err := appCtx.Server().GRPC.GuildClient.GetGuilds(c.Request().Context(), &pb.GetGuildsRequest{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch guilds")
	}

	result := make([]Guild, 0, len(resp.Guilds))
	for _, g := range resp.Guilds {
		result = append(
			result, Guild{
				ID:   g.Id,
				Name: g.Name,
				Icon: g.Icon,
			})
	}

	return c.JSON(http.StatusOK, GuildResponse{Guilds: result})
}

func GetGuildOverview(c echo.Context) error {
	appCtx := c.Request().Context().Value("appctx")
	appCtx, ok := appCtx.(context.AppContext)
	if !ok || appCtx == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}
	q := c.Request().Context().Value("query").(*query.Query)
	guildId := c.Param("guildId")
	guildChannelData, err := q.GuildMessageChannelQuery(guildId)
	guildUserData, err := q.GetGuildMessageCountViaDays(guildId, 7)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch guild channel data",
			})
	}
	return c.JSON(
		http.StatusOK, map[string]interface{}{
			"guildChannelData": guildChannelData,
			"guildUserData":    guildUserData,
		})
}
func Test(c echo.Context) error {
	// Get AppContext from middleware
	appCtx := c.Request().Context().Value("appctx")
	appCtxCasted, ok := appCtx.(context.AppContext)
	if !ok || appCtxCasted == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "App context missing")
	}

	// Get Query instance from middleware
	q := c.Request().Context().Value("query")
	queryInstance, ok := q.(*query.Query)
	if !ok || queryInstance == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Query instance missing")
	}

	// Get guildId from path
	guildId := c.Param("guildId")
	if guildId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing guildId")
	}

	// Bind JSON body to struct
	var body struct {
		StartTime   int64    `json:"startTime"`
		EndTime     int64    `json:"endTime"`
		Interval    string   `json:"interval"`
		Location    string   `json:"location"`
		UserList    []string `json:"userList"`
		ChannelList []string `json:"channelList"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if body.StartTime == 0 || body.EndTime == 0 || body.Interval == "" || body.Location == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields in body")
	}

	params := models.MessageQueryParams{
		GuildID:     guildId,
		StartTime:   body.StartTime,
		EndTime:     body.EndTime,
		Interval:    body.Interval,
		Location:    body.Location,
		UserList:    body.UserList,
		ChannelList: body.ChannelList,
	}

	stats, err := queryInstance.GetGuildMessageStats(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Query error: "+err.Error())
	}

	return c.JSON(
		http.StatusOK, map[string]interface{}{
			"chartInterval": params.Interval,
			"chartData":     stats,
			"hedehode":      "test",
		})
}
