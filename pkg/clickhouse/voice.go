package clickhouse

import (
	"context"
	"discord/config"
	"discord/models"
	"fmt"
	"strings"
	"time"
)

func (c *ConnectionClickhouse) InsertBatchVoice() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := fmt.Sprintf(
		"INSERT INTO discord.voice (guildId, channelId, channelName, userId, username, count, eventTime) VALUES %s",
		strings.Join(c.VoiceRecords, ","))
	err := c.Connection.Exec(ctx, query)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (c *ConnectionClickhouse) BatchVoice(data *models.VoiceQuery) error {
	c.VoiceRecords = append(
		c.VoiceRecords, fmt.Sprintf(
			"('%s', '%s', '%s', '%s', '%s', '%d', toDateTime(%d))",
			data.GuildId,
			data.ChannelId,
			data.ChannelName,
			data.UserId,
			data.Username,
			data.Count,
			data.EventTime))
	if len(c.VoiceRecords) >= config.ClickhouseVoiceLogBatchSize {
		err := c.InsertBatchVoice()
		if err != nil {
			c.VoiceRecords = c.VoiceRecords[:0]
			return err
		}
		c.VoiceRecords = c.VoiceRecords[:0]
	}
	return nil
}
