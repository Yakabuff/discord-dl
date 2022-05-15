package models

//Job args
type JobArgs struct {
	Mode       Mode
	Before     string
	After      string
	FastUpdate bool
	Guild      string
	Channel    string
}

//Archiver system args
type ArchiverArgs struct {
	Mode                Mode
	DownloadMedia       bool
	Token               string
	Output              string
	MediaLocation       string
	DeployPort          int
	Listen              bool
	Deploy              bool
	Input               string
	BlacklistedChannels []string
}

type Mode int

const (
	NONE Mode = iota
	INPUT
	GUILD
	CHANNEL
	INVALID
	EXPORT
)
