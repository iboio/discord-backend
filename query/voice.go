package query

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (q *Query) InsertBatchVoice(records []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := fmt.Sprintf(
		"INSERT INTO discord.voice (guildId, channelId, channelName, userId, username, count, eventTime) VALUES %s",
		strings.Join(records, ","))
	err := q.db.Clickhouse.Exec(ctx, query)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
