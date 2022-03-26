package models

type Args struct {
	Mode                Mode
	DownloadMedia       bool
	Before              string
	After               string
	FastUpdate          bool
	Token               string
	Output              string
	Guild               string
	Channel             string
	Listen              bool
	Deploy              bool
	MediaLocation       string
	DeployPort          int
	Progress            bool
	BlacklistedChannels []string
	Input               string
}
type Mode int

const (
	NONE Mode = iota
	INPUT
	GUILD
	CHANNEL
	INVALID
)
