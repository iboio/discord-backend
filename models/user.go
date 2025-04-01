package models

type HashUser struct {
	GuildId     string `json:"guildId"`
	Username    string `json:"username"`
	EventTime   string `json:"eventTime"`
	ChannelName string `json:"channelName"`
}
