package models

type VoiceStateUpdate struct {
	Type        string `json:"type"`
	GuildId     string `json:"guildId"`
	ChannelId   string `json:"channelId"`
	ChannelName string `json:"channelName"`
	UserId      string `json:"userId"`
	Username    string `json:"username"`
	EventTime   int64  `json:"eventTime"`
}

type VoiceStateSwitch struct {
	Type           string `json:"type"`
	GuildId        string `json:"guildId"`
	OldChannelId   string `json:"oldChannelId"`
	OldChannelName string `json:"oldChannelName"`
	NewChannelId   string `json:"channelId"`
	NewChannelName string `json:"channelName"`
	UserId         string `json:"userId"`
	Username       string `json:"username"`
	EventTime      int64  `json:"eventTime"`
}

type VoiceQuery struct {
	GuildId     string `json:"guildId"`
	ChannelId   string `json:"channelId"`
	ChannelName string `json:"channelName"`
	UserId      string `json:"userId"`
	Username    string `json:"username"`
	Count       int    `json:"count"`
	EventTime   int64  `json:"eventTime"`
}
