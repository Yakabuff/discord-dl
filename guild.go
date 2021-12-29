package main
import(
	"github.com/bwmarrin/discordgo"
	"fmt"
	"os"
	"log"
)
func guild_download(dg *discordgo.Session, a args) error {
	//get all channels from guild into array
	channels, err := dg.GuildChannels(a.guild)
	if err != nil{
		fmt.Println("Could not find guild")
		os.Exit(1)
	}
	//download messages from every channel
	for _, c := range channels{
		if c.Type == discordgo.ChannelTypeGuildText{
			log.Printf("Archiving %s %s", a.guild, c.ID)
			err := download_range_date(dg, a.after, a.before, c.ID, a.guild)
			if err != nil{
				return err;
			}
		}
	}
	

	return nil;
}