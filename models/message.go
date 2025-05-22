package models

import "time"

type Message struct {
	GuildId     string `json:"guildId"`
	ChannelId   string `json:"channelId"`
	ChannelName string `json:"channelName"`
	UserId      string `json:"userId"`
	Username    string `json:"username"`
	EventTime   int64  `json:"eventTime"`
}

type GuildChannelMessageCount struct {
	ChannelID   string
	ChannelName string
	Message     uint64
}

type GuildUserMessageQuery struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Message  uint64 `json:"message"`
}

type GuildUserMessage struct {
	GuildUserMessageQuery
	Profile string `json:"profile"`
}

type GuildMessageDailyStat struct {
	Day          time.Time `json:"day"`
	MessageCount uint64    `json:"messageCount"`
}

type MessageQueryParams struct {
	GuildID     string   `json:"guildId"`
	ChannelList []string `json:"channelList"`
	UserList    []string `json:"userList"`
	StartTime   int64    `json:"startTime"`
	EndTime     int64    `json:"endTime"`
	Interval    string   `json:"interval"`
	Location    string   `json:"location"`
	CustomQuery string   `json:"customQuery"`
}

type MessageStats struct {
	Timestamp string `json:"timestamp"`
	Count     uint64 `json:"count"`
}
