package main
import(
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
	"strings"
)
func messageListen(dg *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == dg.State.User.ID {
		return
	}
	log.Println("[LISTEN] Detected new message. Fetching message " + m.ID + " from" + m.ChannelID)
	//If message contains something that resembles a URL, wait a few seconds for discord to get embed info
	//https://github.com/bwmarrin/discordgo/issues/1066
	if strings.Contains(m.Content, "https://") || strings.Contains(m.Content, "http://"){
		go func(ID string, ChannelID string){
			time.Sleep(time.Second * 5)
			m, err := dg.ChannelMessage(ChannelID, ID)
			if err != nil{
				log.Println("Could not fetch " + m.ID + " from " + m.ChannelID)
			}
			err = addMessage(db, m, false)
			if err != nil{
				log.Println("Could not insert message " + m.ID + " from " + m.ChannelID)
			}
		}(m.ID, m.ChannelID)
	}else{
		err := addMessage(db, m.Message, false)
		if err != nil{
			log.Println("Could not insert message " + m.Message.ID + " from " + m.ChannelID)
		}
	}
}



