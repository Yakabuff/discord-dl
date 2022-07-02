package db

import "testing"

func TestGetChannelsFromGuild(t *testing.T) {
	e := dbConn.InsertChannelID("123")
	if e != nil {
		t.Errorf("Failed to insert channel id")
	}
	e = dbConn.InsertGuildID("456")
	if e != nil {
		t.Errorf("Failed to insert channel id")
	}

	e = dbConn.InsertChannelTopic("123", "mytopic")
	if e != nil {
		t.Errorf("Failed to insert channel id")
	}

	e = dbConn.InsertChannelNames("123", "myname")
	if e != nil {
		t.Errorf("Failed to insert channel id")
	}
	e = dbConn.UpdateChannelMetaTransaction("123", false, "456")
	if e != nil {
		t.Errorf("Failed to insert channel id")
	}

	chans, e := dbConn.GetChannelsFromGuild("456")
	if e != nil {
		t.Errorf("Failed to insert channel id")
	}

	if len(chans) == 1 {
		if chans[0].ChannelID != "123" && chans[0].GuildID != "456" && chans[0].Name != "myname" && chans[0].Topic != "mytopic" {
			t.Error("Invalid values")
		}
	} else {
		t.Error("Invalid values")
	}

}
