package batcher

import (
	"discord/config"
	"discord/models"
	"fmt"
)

func (b *Batcher) BatchMessage(data *models.Message) error {
	b.MessageRecords = append(
		b.MessageRecords, fmt.Sprintf(
			"('%s', '%s', '%s', '%s', '%s', toDateTime(%d))",
			data.GuildId, data.ChannelId, data.ChannelName, data.UserId, data.Username, data.EventTime,
		))
	if len(b.MessageRecords) >= config.ClickhouseMessageBatchSize {
		err := b.query.InsertBatchMessage(b.MessageRecords)
		b.MessageRecords = b.MessageRecords[:0]
		return err
	}
	return nil
}
