// +build test
package archiver

import (
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yakabuff/discord-dl/models"
)

func (a Archiver) CreateConnection() (error, *discordgo.Session) {
	log.Println("Logging in token" + a.Args.Token)
	dg, err := discordgo.New(a.Args.Token)
	if err != nil {
		log.Println(err.Error())
		return err, nil
	}
	err = dg.Open()

	if err != nil {
		log.Println(err.Error())
		return err, nil
	}
	u, err := dg.User("@me")

	log.Printf("discord-dl has succesfully logged into %s#%s %s\n", u.Username, u.Discriminator, u.ID)

	return nil, dg
}

func (a Archiver) addHandlers() {
	a.Dg.Identify.Intents = discordgo.IntentsGuildMessages
	log.Println("Adding handlers")
	a.Dg.AddHandler(a.messageListen)
	a.Dg.AddHandler(a.messageUpdateListen)
}

func (a Archiver) messageListen(dg *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == dg.State.User.ID {
		return
	}
	log.Println("[LISTEN] Detected new message. Fetching message " + m.ID + " from" + m.ChannelID)
	//If message contains something that resembles a URL, wait a few seconds for discord to get embed info
	//https://github.com/bwmarrin/discordgo/issues/1066
	if strings.Contains(m.Content, "https://") || strings.Contains(m.Content, "http://") {
		go func(ID string, ChannelID string) {
			time.Sleep(time.Second * 5)
			m, err := dg.ChannelMessage(ChannelID, ID)
			if err != nil {
				log.Println("Could not fetch " + m.ID + " from " + m.ChannelID)
			}
			err = a.InsertMessage(m, false)
			if err != nil {
				log.Println("Could not insert message " + m.ID + " from " + m.ChannelID)
			}
		}(m.ID, m.ChannelID)
	} else {
		err := a.InsertMessage(m.Message, false)
		if err != nil {
			log.Println("Could not insert message " + m.Message.ID + " from " + m.ChannelID)
		}
	}
}

func (a Archiver) messageUpdateListen(dg *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.Author.ID == dg.State.User.ID {
		return
	}
	log.Println("[LISTEN] Detected new message edit. Fetching message " + m.ID + " from" + m.ChannelID)
	//filter out all messages that do not have an edit timestamp. Only listen for content edits
	edited_timestamp := m.EditedTimestamp.Unix()

	if m.EditedTimestamp != nil {
		edit := models.NewEdit(m.ID, edited_timestamp, m.Content)
		a.InsertEdit(edit)
	}
}
