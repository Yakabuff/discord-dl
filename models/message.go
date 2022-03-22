package models

import (
	"errors"
)

type MessageOut struct {
	MessageId        string
	ChannelId        string
	GuildId          string
	MessageTimestamp int64
	Content          string
	SenderId         string
	SenderName       string
	ReplyTo          string
	EditTime         string
	ThreadId         string
	Edits            []Edit
	Embeds           []Embed
	Attachments      []Attachment
}

type Message struct {
	MessageId        string
	ChannelId        string
	GuildId          string
	MessageTimestamp int64
	Content          string
	SenderId         string
	SenderName       string
	ReplyTo          string
	EditTime         int64
	ThreadId         string
}

func NewMessage(MessageId string, ChannelId string, GuildId string, MessageTimestamp int64, Content string, SenderId string, SenderName string, ReplyTo string, EditTime int64, ThreadId string) Message {
	msg := Message{
		MessageId:        MessageId,
		ChannelId:        ChannelId,
		GuildId:          GuildId,
		MessageTimestamp: MessageTimestamp,
		Content:          Content,
		SenderId:         SenderId,
		SenderName:       SenderName,
		ReplyTo:          ReplyTo,
		EditTime:         EditTime,
		ThreadId:         ThreadId}
	return msg
}

var FastUpdateError = errors.New("Reached downloaded message")
