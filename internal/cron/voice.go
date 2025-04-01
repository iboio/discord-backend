package cron

import (
	"discord/config"
	"discord/internal/methods"
	"discord/models"
	redisdb "discord/pkg/redis"
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

func (h *Handler) VoiceLogScheduler() {
	c := cron.New()
	_, err := c.AddFunc(
		"*/1 * * * *", func() {
			nowEpoch := time.Now().Unix()
			redisdb.StringSet(0, "CronTime", methods.Int64ToString(nowEpoch))
			listData, err := redisdb.GetAllList(0)
			if err != nil {
				fmt.Println("Error fetching voice list:", err)
				return
			}
			for channel := range listData {
				for _, user := range listData[channel] {
					userHash := redisdb.HashGetAll(0, user)
					if len(userHash) == 0 {
						continue
					}
					var userData models.HashUser
					if err := methods.MapToStruct(userHash, &userData); err != nil {
						fmt.Println("Error mapping user data:", err)
						continue
					}
					eventTime, err := methods.StringToInt64(userData.EventTime)
					if err != nil {
						fmt.Println("Error converting event time:", err)
					}
					count := config.Period
					if nowEpoch-eventTime < config.Period {
						count = nowEpoch - eventTime
					}
					data := models.VoiceQuery{
						GuildId:     userData.GuildId,
						ChannelId:   channel,
						ChannelName: userData.ChannelName,
						UserId:      user,
						Username:    userData.Username,
						Count:       int(count),
						EventTime:   nowEpoch,
					}
					err = h.ch.BatchVoice(&data)
					if err != nil {
						return
					}
				}
			}
		})
	if err != nil {
		return
	}
	c.Start()
	select {}
}
