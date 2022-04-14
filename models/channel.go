package models

type Channel struct {
	ChannelID string
	Name      string
	Topic     string
	GuildID   string
}

type Channels struct {
	Channels []Channel
}
