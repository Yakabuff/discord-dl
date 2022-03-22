package models

type Edit struct {
	MessageId string
	EditTime  int64
	Content   string
}

func NewEdit(MessageId string, EditTime int64, Content string) Edit {
	edit := Edit{MessageId: MessageId, EditTime: EditTime, Content: Content}
	return edit
}
