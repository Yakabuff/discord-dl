package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var dbConn *Db

func setup() {
	db, err := Init_db("test.db")
	dbConn = db
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cleanup() {
	dbConn.DbConnection.Close()
	os.Remove("test.db")
}

func teardown() {
	cleanup()
	setup()
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	// cleanup()
}

func TestInsertGuildID(t *testing.T) {
	e := dbConn.InsertGuildID("123")
	if e != nil {
		t.Errorf("Failed to insert guild id")
	}
}

func TestInsertChannelID(t *testing.T) {
	e := dbConn.InsertChannelID("123")
	if e != nil {
		t.Errorf("Failed to insert channel id")
	}
}

func TestInsertGuildName(t *testing.T) {

	e := dbConn.InsertGuildID("123123123")
	if e != nil {
		t.Errorf("Failed to insert guild id")
		t.Log(e.Error())
	}

	e = dbConn.InsertGuildNames("123123123", "chan1")
	if e != nil {
		t.Errorf("Failed to upsert channel name chan1 1")
		t.Errorf(e.Error())
	}

	e = dbConn.InsertGuildNames("123123123", "chan1") //Fail

	if assert.Error(t, e) {
		res := (e.Error() == "guildNameTrigger violated")

		assert.Equal(t, true, res)
	}
	e = dbConn.InsertGuildNames("123123123", "chan2") //pass
	if e != nil {
		t.Errorf("Failed to upsert channel name chan2")
		t.Errorf(e.Error())
	}

	e = dbConn.InsertGuildNames("123123123", "chan1") //Pass
	if e != nil {
		t.Errorf("Failed to upsert channel name chan1 3")
		t.Errorf(e.Error())
	}

	e = dbConn.InsertGuildNames("123123123", "chan1") //Fail

	if assert.Error(t, e) {
		res := (e.Error() == "guildNameTrigger violated")

		assert.Equal(t, true, res)
	}
	e = dbConn.InsertGuildNames("123123123", "chan1") // Fail

	if assert.Error(t, e) {
		res := (e.Error() == "guildNameTrigger violated")

		assert.Equal(t, true, res)
	}
	teardown()
}

func TestUpsertGuildMeta(t *testing.T) {
	e := dbConn.InsertGuildID("123123123")
	if e != nil {
		t.Errorf("Failed to insert guild id")
		t.Log(e.Error())
	}

	e = dbConn.InsertGuildNames("123123123", "chan1")
	if e != nil {
		t.Errorf("Failed to upsert channel name chan1 1")
		t.Errorf(e.Error())
	}
	e = dbConn.InsertGuildIcons("123123123", "iconhash1")
	if e != nil {
		t.Errorf("Failed to upsert icon")
		t.Errorf(e.Error())
	}

	e = dbConn.InsertGuildBanner("123123123", "bannerhash1")
	if e != nil {
		t.Errorf("Failed to upsert banner")
		t.Errorf(e.Error())
	}

	e = dbConn.UpdateGuildMetaTransaction("123123123")
	if e != nil {
		t.Errorf(e.Error())
	}
	///////////////////////////////////////
	e = dbConn.InsertGuildNames("123123123", "chan2")
	if e != nil {
		t.Errorf("Failed to upsert channel name chan2 1")
		t.Errorf(e.Error())
	}
	e = dbConn.InsertGuildIcons("123123123", "iconhash2")
	if e != nil {
		t.Errorf("Failed to upsert icon")
		t.Errorf(e.Error())
	}

	e = dbConn.InsertGuildBanner("123123123", "bannerhash2")
	if e != nil {
		t.Errorf("Failed to upsert banner")
		t.Errorf(e.Error())
	}

	e = dbConn.UpdateGuildMetaTransaction("123123123")
	if e != nil {
		t.Errorf(e.Error())
	}
	////////////////////////////////////////
	e = dbConn.InsertGuildNames("123123123", "chan2")
	if assert.Error(t, e) {
		res := (e.Error() == "guildNameTrigger violated")

		assert.Equal(t, true, res)
	}
	e = dbConn.InsertGuildIcons("123123123", "iconhash2")
	if assert.Error(t, e) {
		res := (e.Error() == "guildIconTrigger violated")

		assert.Equal(t, true, res)
	}

	e = dbConn.InsertGuildBanner("123123123", "bannerhash2")
	if assert.Error(t, e) {
		res := (e.Error() == "guildBannerTrigger violated")

		assert.Equal(t, true, res)
	}

	e = dbConn.UpdateGuildMetaTransaction("123123123")
	if e != nil {
		t.Errorf(e.Error())
	}
	teardown()
}

func TestInsertGuildIcon(t *testing.T) {

}

func TestInsertGuildBanner(t *testing.T) {

}

func TestUpsertChannelMeta(t *testing.T) {
	e := dbConn.InsertGuildID("234234234")
	if e != nil {
		t.Errorf("Failed to insert guild id")
		t.Log(e.Error())
	}
	e = dbConn.InsertChannelID("123123123")
	if e != nil {
		t.Errorf("Failed to insert guild id")
		t.Log(e.Error())
	}

	e = dbConn.InsertChannelNames("123123123", "chan1name")
	if e != nil {
		t.Errorf("Failed to upsert channel name chan1 1")
		t.Errorf(e.Error())
	}
	e = dbConn.InsertChannelTopic("123123123", "chan1topic")
	if e != nil {
		t.Errorf("Failed to upsert icon")
		t.Errorf(e.Error())
	}
	e = dbConn.UpdateChannelMetaTransaction("123123123", true, "234234234")
	if e != nil {
		t.Errorf(e.Error())
	}
	// teardown()
}
