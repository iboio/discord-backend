package query

import (
	"context"
	"discord/internal/methods"
	"discord/models"
	redisdb "discord/pkg/redis"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"math/rand"
	"strings"
	"time"
)

func (q *Query) InsertBatchMessage(records []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := fmt.Sprintf(
		"INSERT INTO discord.message (guildId, channelId, channelName, userId, username, eventTime) VALUES %s",
		strings.Join(records, ","))
	return q.db.Clickhouse.Exec(ctx, query)
}

func (q *Query) GetUserMessageViaDays(guildId string, day int) ([]models.GuildUserMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := fmt.Sprintf(
		`
    SELECT userId, username, count(userId) as message
    FROM discord.message
    WHERE guildId = '%s' 
    AND eventTime >= toStartOfDay(now()) - INTERVAL %d DAY
    GROUP BY username, userId
    ORDER BY message DESC;`, guildId, day)
	rows, err := q.db.Clickhouse.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows driver.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var results []models.GuildUserMessage
	for rows.Next() {
		var cmc models.GuildUserMessageQuery
		if err := rows.Scan(
			&cmc.UserId,
			&cmc.Username,
			&cmc.Message); err != nil {
			return nil, err
		}
		profileUrl := "default"
		avatarHash := redisdb.HashGet(1, "user-profile-hash", cmc.UserId)
		if avatarHash != "" {
			profileUrl = fmt.Sprintf(
				"https://cdn.discordapp.com/avatars/%s/%s.webp?size=80",
				cmc.UserId, avatarHash)
		}
		guildUserMessage := models.GuildUserMessage{
			GuildUserMessageQuery: cmc,
			Profile:               profileUrl,
		}
		results = append(results, guildUserMessage)
	}
	return results, nil
}

func (q *Query) GetGuildMessageCountViaDays(guildId string, day int) ([]models.GuildMessageDailyStat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := fmt.Sprintf(
		`
        SELECT
            d AS day,
            sum(m.guildId = '%s') AS message_count
        FROM
            (
                SELECT toDate(now()) - INTERVAL number DAY AS d
                FROM numbers(%d)
            ) AS days
        LEFT JOIN discord.message AS m
            ON toDate(m.eventTime) = days.d
        GROUP BY d
        ORDER BY d ASC;
    `, guildId, day)

	rows, err := q.db.Clickhouse.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows driver.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var result []models.GuildMessageDailyStat

	for rows.Next() {
		var guildUserMessageRecord models.GuildMessageDailyStat
		if err := rows.Scan(&guildUserMessageRecord.Day, &guildUserMessageRecord.MessageCount); err != nil {
			return nil, err
		}
		result = append(result, guildUserMessageRecord)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (q *Query) GuildMessageChannelQuery(guildId string) ([]models.GuildChannelMessageCount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := fmt.Sprintf(
		"SELECT channelId, channelName, count(channelId) as message FROM discord.message WHERE guildId = '%s' GROUP BY channelId,channelName ORDER BY message DESC;",
		guildId)

	rows, err := q.db.Clickhouse.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows driver.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var results []models.GuildChannelMessageCount
	for rows.Next() {
		var cmc models.GuildChannelMessageCount
		if err := rows.Scan(&cmc.ChannelID, &cmc.ChannelName, &cmc.Message); err != nil {
			return nil, err
		}
		results = append(results, cmc)
	}
	return results, nil

}

func (q *Query) GetGuildMessageStats(params models.MessageQueryParams) ([]models.MessageStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query, err := methods.MessageQueryGenerator(params)
	if err != nil {
		return nil, err
	}

	rows, err := q.db.Clickhouse.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.MessageStats
	for rows.Next() {
		var stat models.MessageStats
		if err := rows.Scan(&stat.Timestamp, &stat.Count); err != nil {
			return nil, err
		}
		results = append(results, stat)
	}

	return results, nil
}
func (q *Query) DummyData(rowCount int) error {
	batch, err := q.db.Clickhouse.PrepareBatch(context.Background(), "INSERT INTO discord.message")
	if err != nil {
		return err
	}
	for i := 0; i < rowCount; i++ {
		guildId := fmt.Sprintf("%d", rand.Intn(3)+1)
		channelId := fmt.Sprintf("%d", rand.Intn(10)+1)
		channelName := "channel_" + channelId
		userId := fmt.Sprintf("u%d", rand.Intn(50)+1)
		username := "User_" + userId

		eventTime := time.Now().
			Add(-time.Duration(rand.Intn(30*24)) * time.Hour). // last 30 days
			Add(-time.Duration(rand.Intn(60)) * time.Minute).
			Add(-time.Duration(rand.Intn(60)) * time.Second)

		if err := batch.Append(
			guildId,
			channelId,
			channelName,
			userId,
			username,
			eventTime,
		); err != nil {
			return err
		}
	}

	return batch.Send()
}
