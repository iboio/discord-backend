package models

type Message struct {
	GuildId     string `json:"guildId"`
	ChannelId   string `json:"channelId"`
	ChannelName string `json:"channelName"`
	UserId      string `json:"userId"`
	Username    string `json:"username"`
	EventTime   int64  `json:"eventTime"`
}
