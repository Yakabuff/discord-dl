package archiver

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func (a Archiver) GuildDownload(guildID string, fastUpdate bool, after string, before string) error {
	//get all channels from guild into array
	channels, err := a.Dg.GuildChannels(guildID)
	if err != nil {
		log.Println("Could not find guild")
		os.Exit(1)
	}
	//download messages from every channel
	for _, c := range channels {
		if c.Type == discordgo.ChannelTypeGuildText && !contains(a.Args.BlacklistedChannels, c.ID) {

			log.Printf("Archiving guild: %s channel: %s", guildID, c.ID)
			// err := a.ChannelDownload(c.ID, fastUpdate, after, before)
			log.Println(err)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a Archiver) GetChannelsGuild(guildID string) ([]string, error) {
	channels, err := a.Dg.GuildChannels(guildID)
	if err != nil {
		log.Println("Could not find guild")
		return nil, err
	}
	res := make([]string, len(channels))
	for _, val := range channels {
		res = append(res, val.ID)
	}
	return res, nil
}

func contains(channels []string, id string) bool {
	for _, a := range channels {
		if a == id {
			return true
		}
	}
	return false
}
