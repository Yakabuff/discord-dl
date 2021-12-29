package main

import(
    // "errors"
    "testing"
    "github.com/bwmarrin/discordgo"
    "os"
    "log"
)

func TestAddMessageFresh(t *testing.T){
    t.Logf("Creating DB")
    e := os.Remove("bigbrother.db")
    if e != nil{
        t.Fatal(e)
    }
    db, err := init_db();

    if err != nil{
        t.Fail()
    }
    if db == nil {
        t.Fail()
    }

    var timestamp discordgo.Timestamp = "timestamp"
    var author *discordgo.User = &discordgo.User{Username: "username", ID: "author_id"}
    var reply *discordgo.MessageReference = &discordgo.MessageReference{MessageID: "reply_id"}
    //to-do attachments
    msg := discordgo.Message{
    ID: "034asdf", 
    ChannelID: "testChannelID", 
    GuildID: "guild_id", 
    Content: "message",
    EditedTimestamp: timestamp,
    Author: author,
    MessageReference: reply,
    }
    t.Logf("HELLO")
    t.Logf(msg.ChannelID)

    err = addMessage(db, msg)

    if err != nil{
        log.Println(err)
    }

}

