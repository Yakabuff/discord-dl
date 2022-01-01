package main
import(
	"github.com/bwmarrin/discordgo"
	"log"
)
func messageListen(dg *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == dg.State.User.ID {
		return
	}
	log.Println("DETECTED NEW MESSAGE !! Inserting message " + m.Message.ID + " from" + m.Message.ChannelID)
	log.Println(m.Message.Embeds)
	// err := addMessage(db, m.Message, false)
	// log.Println("Listen: msg added")
	// if err != nil{
	// 	log.Println("Could not insert message " + m.Message.ID + " from " + m.ChannelID)
	// }
}



