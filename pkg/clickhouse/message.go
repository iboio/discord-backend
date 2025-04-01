package clickhouse

import (
	"context"
	"discord/config"
	"discord/models"
	"fmt"
	"strings"
	"time"
)

func (c *ConnectionClickhouse) Test() {

}

func (c *ConnectionClickhouse) InsertBatchMessage() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := fmt.Sprintf(
		"INSERT INTO discord.message (guildId, channelId, channelName, userId, username, eventTime) VALUES %s",
		strings.Join(c.MessageRecords, ","))
	err := c.Connection.Exec(ctx, query)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (c *ConnectionClickhouse) BatchMessage(data *models.Message) error {
	c.MessageRecords = append(
		c.MessageRecords, fmt.Sprintf(
			"('%s', '%s', '%s', '%s', '%s', toDateTime(%d))",
			data.GuildId,
			data.ChannelId,
			data.ChannelName,
			data.UserId,
			data.Username,
			data.EventTime))
	if len(c.MessageRecords) >= config.ClickhouseMessageBatchSize {
		err := c.InsertBatchMessage()
		if err != nil {
			c.MessageRecords = c.MessageRecords[:0]
			return err
		}
		c.MessageRecords = c.MessageRecords[:0]
	}
	return nil
}
