package batcher

import (
	"discord/config"
	"discord/models"
	"fmt"
)

func (b *Batcher) BatchVoice(data *models.VoiceQuery) error {
	b.VoiceRecords = append(
		b.VoiceRecords, fmt.Sprintf(
			"('%s', '%s', '%s', '%s', '%s', '%d', toDateTime(%d))",
			data.GuildId,
			data.ChannelId,
			data.ChannelName,
			data.UserId,
			data.Username,
			data.Count,
			data.EventTime))
	if len(b.VoiceRecords) >= config.ClickhouseVoiceLogBatchSize {
		err := b.query.InsertBatchVoice(b.VoiceRecords)
		if err != nil {
			b.VoiceRecords = b.VoiceRecords[:0]
			return err
		}
		b.VoiceRecords = b.VoiceRecords[:0]
	}
	return nil
}
