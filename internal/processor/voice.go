package processor

import (
	"discord/config"
	"discord/internal/methods"
	"discord/models"
	redisdb "discord/pkg/redis"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"strconv"
)

type Type struct {
	Type string `json:"type"`
}

const (
	VoiceJoin   = "join"
	VoiceLeave  = "left"
	VoiceSwitch = "switch"
)

func (h *Handler) VoiceProcessor() {
	fmt.Println("Voice Processor started")
	fmt.Println(config.JetstreamSubjectsEventVoice)
	_, err := h.js.Jetstream.Subscribe(
		config.JetstreamSubjectsEventVoice,
		func(msg *nats.Msg) {
			stateCheck := Type{}
			err := json.Unmarshal(msg.Data, &stateCheck)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(stateCheck.Type)
			if stateCheck.Type == VoiceJoin {
				err := h.JoinEvent(string(msg.Data))
				if err != nil {
					fmt.Println(err)
				}
				err = msg.Ack()
				if err != nil {
					fmt.Println(err)
				}
			} else if stateCheck.Type == VoiceLeave {
				err := h.LeftEvent(string(msg.Data))
				if err != nil {
					fmt.Println(err)
				}
				err = msg.Ack()
				if err != nil {
					fmt.Println(err)
				}
			} else if stateCheck.Type == VoiceSwitch {
				err := h.SwitchEvent(string(msg.Data))
				if err != nil {
					fmt.Println(err)
				}
				err = msg.Ack()
				if err != nil {
					fmt.Println(err)
				}
			}
		},
		nats.ManualAck(),
		nats.Durable("backend-voice-consumer"))
	if err != nil {
		fmt.Println("Error subscribing to Jetstream:", err)
	}
}
func (h *Handler) JoinEvent(rawData string) error {
	data, err := methods.StringToStruct[models.VoiceStateUpdate](rawData)
	if err != nil {
		fmt.Println(err)
		return err
	}
	redisdb.ListSet(0, data.ChannelId, data.UserId)
	redisdb.HashSet(0, data.UserId, "GuildId", data.GuildId)
	redisdb.HashSet(0, data.UserId, "ChannelName", data.ChannelName)
	redisdb.HashSet(0, data.UserId, "Username", data.Username)
	redisdb.HashSet(
		0,
		data.UserId,
		"EventTime",
		strconv.FormatInt(data.EventTime, 10))
	//h.VoiceTraffic(data, "join")
	return nil
}

func (h *Handler) LeftEvent(rawData string) error {
	data, err := methods.StringToStruct[models.VoiceStateUpdate](rawData)
	if err != nil {
		fmt.Println(err)
		return err
	}
	joinEpochTime, err := methods.StringToInt64(redisdb.HashGet(0, data.UserId, "EventTime"))
	if err != nil {
		fmt.Println(err)
	}
	eventEpochTime := data.EventTime
	count := LeftCountCalculator(joinEpochTime, eventEpochTime)
	h.UserLeftChannelVoiceLog(data, count)
	//h.VoiceTraffic(data, "left")
	redisdb.HashDelAll(0, data.UserId)
	redisdb.ListDelElement(0, data.ChannelId, data.UserId)
	return nil
}

func (h *Handler) SwitchEvent(rawData string) error {
	data, err := methods.StringToStruct[models.VoiceStateSwitch](rawData)
	if err != nil {
		fmt.Println(err)
		return err
	}
	joinEpochTime, err := methods.StringToInt64(redisdb.HashGet(0, data.UserId, "EventTime"))
	if err != nil {
		fmt.Println(err)
	}
	leftEpochTime := data.EventTime
	count := LeftCountCalculator(joinEpochTime, leftEpochTime)
	oldChannel := models.VoiceStateUpdate{
		GuildId:     data.GuildId,
		ChannelId:   data.OldChannelId,
		ChannelName: data.OldChannelName,
		UserId:      data.UserId,
		Username:    data.Username,
		EventTime:   data.EventTime,
	}
	h.UserLeftChannelVoiceLog(oldChannel, count)
	//h.VoiceTraffic(oldChannel, "left")
	redisdb.ListMove(0, data.OldChannelId, data.NewChannelId, data.UserId)
	redisdb.HashSet(0, data.UserId, "GuildId", data.GuildId)
	redisdb.HashSet(0, data.UserId, "ChannelName", data.NewChannelName)
	redisdb.HashSet(0, data.UserId, "Username", data.Username)
	redisdb.HashSet(
		0,
		data.UserId,
		"EventTime",
		methods.Int64ToString(data.EventTime))
	//newChannel := models.VoiceStateUpdate{
	//	GuildId:     data.GuildId,
	//	ChannelId:   data.NewChannelId,
	//	ChannelName: data.NewChannelName,
	//	UserId:      data.UserId,
	//	Username:    data.Username,
	//	EventTime:   data.EventTime,
	//}
	//h.VoiceTraffic(newChannel, "join")
	return nil
}

func LeftCountCalculator(joinEventTime, leftEventTime int64) int64 {
	var count int64
	cronEpochTime, err := strconv.ParseInt(redisdb.StringGet(0, "CronTime"), 10, 64)
	if err != nil {
		fmt.Println("Error parsing event time:", err)
	}
	if cronEpochTime < joinEventTime && cronEpochTime+config.Period > leftEventTime {
		count = leftEventTime - joinEventTime
	} else {
		count = leftEventTime - cronEpochTime
	}
	return count
}
func (h *Handler) UserLeftChannelVoiceLog(voiceData models.VoiceStateUpdate, count int64) {
	data := models.VoiceQuery{
		GuildId:     voiceData.GuildId,
		ChannelId:   voiceData.ChannelId,
		ChannelName: voiceData.ChannelName,
		UserId:      voiceData.UserId,
		Username:    voiceData.Username,
		Count:       int(count),
		EventTime:   voiceData.EventTime,
	}
	err := h.ch.BatchVoice(&data)
	if err != nil {
		fmt.Println(err)
	}
}

//func (h *Handler) VoiceTraffic(voiceData models.VoiceStateUpdate, status string) {
//	data := models.VoiceLogTraffic{
//		GuildId:     voiceData.GuildId,
//		ChannelId:   voiceData.ChannelId,
//		ChannelName: voiceData.ChannelName,
//		UserId:      voiceData.UserId,
//		Status:      status,
//		Username:    voiceData.Username,
//		EventTime:   voiceData.EventTime,
//	}
//	err := h.ch.InsertVoiceTrafficLog(data)
//	if err != nil {
//		fmt.Println(err)
//	}
//}
