package models

type AttachmentOut struct {
	AttachmentId       string
	MessageId          string
	AttachmentFilename string
	AttachmentUrl      string
	AttachmentHash     string
	ResourcePath       string
	ResourceType       string
}

type Attachment struct {
	AttachmentId       string
	MessageId          string
	AttachmentFilename string
	AttachmentUrl      string
	AttachmentHash     string
}

func NewAttachment(AttachmentId string,
	MessageId string,
	AttachmentFilename string,
	AttachmentUrl string,
	AttachmentHash string,
) Attachment {

	a := Attachment{AttachmentId: AttachmentId, MessageId: MessageId, AttachmentFilename: AttachmentFilename, AttachmentUrl: AttachmentUrl, AttachmentHash: AttachmentHash}
	return a
}
