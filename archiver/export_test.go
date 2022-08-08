package archiver

import (
	"os"
	"testing"

	"github.com/yakabuff/discord-dl/models"
)

func TestExport(t *testing.T) {
	msg := models.MessageJson{MessageId: "messageid", ChannelId: "channelid"}
	msg2 := models.MessageJson{MessageId: "messageid2", ChannelId: "channelid"}
	err := WriteMessageJson(msg, "420420")
	if err != nil {
		t.Error(err)
	}
	err = WriteMessageJson(msg2, "420420")
	if err != nil {

		t.Error(err)
	}
	f, _ := os.Open("channelid_420420.json")
	stat, err := f.Stat()
	if err != nil {
		t.Error(err)
	}
	filesize := stat.Size()

	if filesize != 407 {
		t.Error(err)
	}

	e := os.Remove("channelid_420420.json")
	if e != nil {
		t.Error(err)
	}
}
