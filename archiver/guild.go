package archiver

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yakabuff/discord-dl/common"
)

func (a Archiver) GetChannelsGuild(guildID string) ([]string, error) {

	channels, err := a.Dg.GuildChannels(guildID)
	if err != nil {
		log.Error("Could not find guild")
		return nil, err
	}
	res := []string{}
	for i, val := range channels {
		_, err := a.Dg.Channel(val.ID)
		if err == nil {
			if val.Type == discordgo.ChannelTypeGuildText && !common.Contains(a.Args.BlacklistedChannels, val.ID) {
				res = append(res, channels[i].ID)
			}

		}

	}
	return res, nil
}

// func contains(channels []string, id string) bool {
// 	for _, a := range channels {
// 		if a == id {
// 			return true
// 		}
// 	}
// 	return false
// }

// c, err := a.Dg.Channel(channel)
// if err != nil {
// 	if err.(*discordgo.RESTError).Message.Code == 50001 {
// 		log.Println("Do not have permission for channel")
// 		return nil
// 	}
// 	log.Println(err)

// 	return err
// }
