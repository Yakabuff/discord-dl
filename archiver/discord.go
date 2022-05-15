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

	user, _ := a.Dg.User("@me")
	// Ignore all messages created by the bot itself
	if m.Author.ID == user.ID {
		return
	}
	log.Println("[LISTEN] Detected new message. Fetching message " + m.ID + " from" + m.ChannelID)
	guildID := m.GuildID
	//If message contains something that resembles a URL, wait a few seconds for discord to get embed info
	//https://github.com/bwmarrin/discordgo/issues/1066
	if strings.Contains(m.Content, "https://") || strings.Contains(m.Content, "http://") {
		go func(ID string, ChannelID string) {
			time.Sleep(time.Second * 5)
			m, err := dg.ChannelMessage(ChannelID, ID)
			if err != nil {
				log.Println("Could not fetch " + m.ID + " from " + m.ChannelID)
				return
			}
			m.GuildID = guildID
			err = a.InsertMessage(m, false, a.Args.DownloadMedia)
			if err != nil {
				log.Println("Could not insert message " + m.ID + " from " + m.ChannelID)
			}
		}(m.ID, m.ChannelID)
	} else {
		m, err := dg.ChannelMessage(m.ChannelID, m.ID)
		if err != nil {
			log.Println("Could not fetch " + m.ID + " from " + m.ChannelID)
			return
		}
		m.GuildID = guildID
		err = a.InsertMessage(m, false, a.Args.DownloadMedia)
		if err != nil {
			log.Println("Could not insert message " + m.ID + " from " + m.ChannelID)
		}
	}
}

func (a Archiver) messageUpdateListen(dg *discordgo.Session, m *discordgo.MessageUpdate) {
	//Note: If message with link is sent, it does not  return all fields.... Get message ID and channelID and retrieve message this way instead.

	message, err := dg.ChannelMessage(m.ChannelID, m.ID)
	if err != nil {
		log.Println("Failed to get edit: " + m.ID + " " + m.ChannelID)
	}

	if message.Author.ID == dg.State.User.ID {
		return
	}
	log.Println("[LISTEN] Detected new message edit. Fetching message " + m.ID + " from" + m.ChannelID)
	//filter out all messages that do not have an edit timestamp. Only listen for content edits

	if m.EditedTimestamp != nil {
		edited_timestamp := message.EditedTimestamp.Unix()
		edit := models.NewEdit(message.ID, edited_timestamp, message.Content)
		a.InsertEdit(edit)
	}
}
