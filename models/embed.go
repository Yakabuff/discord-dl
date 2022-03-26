package models

type EmbedOut struct {
	MessageId             string
	EmbedUrl              string
	EmbedTitle            string
	EmbedDescription      string
	EmbedTimestamp        string
	EmbedThumbnailUrl     string
	EmbedThumbnailHash    string
	EmbedImageUrl         string
	EmbedImageHash        string
	EmbedVideoUrl         string
	EmbedVideoHash        string
	EmbedFooter           string
	EmbedAuthorName       string
	EmbedAuthorUrl        string
	EmbedField            string
	ResourcePathThumbnail string
	ResourcePathImage     string
	ResourcePathVideo     string
}

type Embed struct {
	MessageId          string
	EmbedUrl           string
	EmbedTitle         string
	EmbedDescription   string
	EmbedTimestamp     string
	EmbedThumbnailUrl  string
	EmbedThumbnailHash string
	EmbedImageUrl      string
	EmbedImageHash     string
	EmbedVideoUrl      string
	EmbedVideoHash     string
	EmbedFooter        string
	EmbedAuthorName    string
	EmbedAuthorUrl     string
	EmbedField         string
}

func NewEmbed(MessageId string,
	EmbedUrl string,
	EmbedTitle string,
	EmbedDescription string,
	EmbedTimestamp string,
	EmbedThumbnailUrl string,
	EmbedThumbnailHash string,
	EmbedImageUrl string,
	EmbedImageHash string,
	EmbedVideoUrl string,
	EmbedVideoHash string,
	EmbedFooter string,
	EmbedAuthorName string,
	EmbedAuthorUrl string,
	EmbedField string) Embed {

	e := Embed{MessageId: MessageId,
		EmbedUrl:           EmbedUrl,
		EmbedTitle:         EmbedTitle,
		EmbedDescription:   EmbedDescription,
		EmbedTimestamp:     EmbedTimestamp,
		EmbedThumbnailUrl:  EmbedThumbnailUrl,
		EmbedThumbnailHash: EmbedThumbnailHash,
		EmbedImageUrl:      EmbedImageUrl,
		EmbedImageHash:     EmbedImageHash,
		EmbedVideoUrl:      EmbedVideoUrl,
		EmbedVideoHash:     EmbedVideoHash,
		EmbedFooter:        EmbedFooter,
		EmbedAuthorName:    EmbedAuthorName,
		EmbedAuthorUrl:     EmbedAuthorUrl,
		EmbedField:         EmbedField}
	return e
}
