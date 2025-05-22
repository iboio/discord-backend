package methods

import (
	"discord/models"
	"fmt"
	"strings"
)

func MessageQueryGenerator(params models.MessageQueryParams) (string, error) {
	var groupFormat string
	switch params.Interval {
	case "hour":
		groupFormat = "%Y-%m-%d %H:00:00"
	case "day":
		groupFormat = "%Y-%m-%d"
	case "minute":
		groupFormat = "%Y-%m-%d %H:%M:00"
	default:
		return "", fmt.Errorf("invalid interval: %s", params.Interval)
	}
	if params.Location == "" {
		return "", fmt.Errorf("location is empty")
	}
	location := params.Location
	channelList := ""
	userList := ""
	if len(params.ChannelList) != 0 {
		channelList = fmt.Sprintf("AND channelId IN (%s)", QuoteList(params.ChannelList))
	}
	if len(params.UserList) != 0 {
		userList = fmt.Sprintf("AND userId IN (%s)", QuoteList(params.UserList))
	}
	query := `
	SELECT
		formatDateTime(eventTime, '` + groupFormat + `', '` + location + `') as time_bucket,
		count() as total
	FROM discord.message
	WHERE guildId = '` + params.GuildID + `' ` + channelList + ` ` + userList + `
	AND eventTime BETWEEN toDateTime(` + fmt.Sprint(params.StartTime) + `) AND toDateTime(` + fmt.Sprint(params.EndTime) + `)
	GROUP BY time_bucket
	ORDER BY time_bucket ASC
`
	return query, nil
}
func QuoteList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = fmt.Sprintf("'%s'", item)
	}
	return strings.Join(quoted, ", ")
}
