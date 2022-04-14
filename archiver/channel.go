// +build test
package archiver

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/yakabuff/discord-dl/common"
	"github.com/yakabuff/discord-dl/db"
)

func (a Archiver) ChannelDownload(channel string) error {
	// case 1: only after flag.
	// case 2: only before flag
	// case 3: after AND before floag
	// case 4: --fast-update flag
	// case 5: contains no flag
	log.Println("Downloading channel " + channel)

	c, err := a.Dg.Channel(channel)
	if err != nil {
		if err.(*discordgo.RESTError).Message.Code == 50001 {
			log.Println("Do not have permission for channel")
			return nil
		}
		log.Println(err)
		return err
	}

	if c == nil {
		log.Println("Channel is nil")
		return nil
	}

	var guildID string
	if c.Type == discordgo.ChannelTypeGuildText ||
		c.Type == discordgo.ChannelTypeDM ||
		c.Type == discordgo.ChannelTypeGroupDM ||
		c.Type == discordgo.ChannelTypeGuildPublicThread ||
		c.Type == discordgo.ChannelTypeGuildPrivateThread {

		if c.Type == discordgo.ChannelTypeGuildText {
			guildID = c.GuildID
		}
		if a.Args.FastUpdate == true {
			a.DownloadMessages(a.Args.After, a.Args.Before, channel, guildID, true)
		} else {
			//Both before and after flags
			if a.Args.Before != "" && a.Args.After != "" {
				if strings.Contains(a.Args.Before, "-") && strings.Contains(a.Args.Before, "-") {
					a.DownloadRangeDate(a.Args.After, a.Args.Before, channel, guildID)
				} else {
					a.DownloadMessages(a.Args.After, a.Args.Before, channel, guildID, false)
				}
			} else if a.Args.Before == "" && a.Args.After == "" {
				//Both before and after empty
				a.DownloadMessages("", "", channel, guildID, false)
			} else {
				//Before empty OR after empty
				if strings.Contains(a.Args.Before, "-") && a.Args.After == "" {
					a.DownloadRangeDate(a.Args.After, a.Args.Before, channel, guildID)
				} else if strings.Contains(a.Args.After, "-") && a.Args.Before == "" {
					a.DownloadRangeDate(a.Args.After, a.Args.Before, channel, guildID)
				} else if !strings.Contains(a.Args.After, "-") && a.Args.Before == "" {
					a.DownloadMessages(a.Args.After, a.Args.Before, channel, guildID, false)
				} else {
					a.DownloadMessages(a.Args.After, a.Args.Before, channel, guildID, false)
				}
			}
		}

	} else {
		return errors.New("Channel is not a text channel")
	}

	return nil
}

//Convert after and before to Time if applicable
//Convert Time for after and before to Discord Unix
//Run DownloadMessage on before and after unix times
func (a Archiver) DownloadRangeDate(after string, before string, channel_id string, guild_id string) error {

	if before != "" {
		before_time, err := common.DateToTime(before)
		if err != nil {
			return err
		}
		//Generate initial snowflake message ID position
		bt := ((before_time.Unix()+1)*1000 - 1420070400000) << 22
		before_id := strconv.Itoa(int(bt))

		if after != "" {
			after_time, _ := common.DateToTime(after)
			at := ((after_time.Unix()-1)*1000 - 1420070400000) << 22
			after_id := strconv.Itoa(int(at))
			err := a.DownloadMessages(before_id, after_id, channel_id, guild_id, false)
			if err != nil {
				return err
			}
		} else {
			err := a.DownloadMessages(before_id, "", channel_id, guild_id, false)
			if err != nil {
				return err
			}
		}

	} else {
		if after != "" {
			after_time, _ := common.DateToTime(after)
			fmt.Println(after_time)
			at := ((after_time.Unix()-1)*1000 - 1420070400000) << 22
			after_id := strconv.Itoa(int(at))
			a.DownloadMessages("", after_id, channel_id, guild_id, false)
		} else {
			a.DownloadMessages("", "", channel_id, guild_id, false)
		}
	}

	return nil
}

func (a Archiver) DownloadMessages(before_id string, after_id string, channel_id string, guild_id string, fast_update bool) error {

	messages, error := a.Dg.ChannelMessages(channel_id, 100, before_id, "", "")
	if error != nil {
		log.Println(error)
		return error
	}
	//Start archiving messages
	var in_range bool = true
	for len(messages) != 0 && in_range {
		log.Println(len(messages))
		for _, m := range messages {
			// timestamp, _ := discordgo.SnowflakeTimestamp(m.ID)
			id := m.ID
			// content := m.Content;
			// author_id := m.Author.ID;
			// author_username := m.Author.Username;
			m.GuildID = guild_id
			if after_id != "" {
				if id > after_id {
					log.Printf("Downloading messages %s %s %s %s %s\n", m.Timestamp, m.ID, m.Content, m.Author.ID, m.Author.Username)
					before_id = id
					//insert into db
					err := a.InsertMessage(m, fast_update)
					if err != nil {
						if errors.Is(err, db.UniqueConstraintError) {
							return nil
						}
						log.Println(err)
						return err
					}
					//Fetch threads if exist
					if m.Type == 21 {
						log.Println("Thread spotted. Traversing thread: " + m.Thread.ID)
						err := a.DownloadMessages("", "", m.Thread.ID, m.Thread.GuildID, fast_update)
						if err != nil {
							return err
						}
					}

				} else {
					in_range = false
					break
				}
			} else {
				log.Printf("Downloading messagez %s %s %s %s %s %s\n", m.Timestamp, m.ID, m.ChannelID, m.Content, m.Author.ID, m.Author.Username)
				err := a.InsertMessage(m, fast_update)
				before_id = id
				if err != nil {
					if errors.Is(err, db.UniqueConstraintError) {
						return nil
					}
					return err
				}
				if m.Thread != nil {
					log.Println("Thread spotted. Traversing thread: " + m.Thread.ID)
					a.IndexChannel(m.Thread.ID)
					err := a.DownloadMessages("", "", m.Thread.ID, m.Thread.GuildID, fast_update)
					if err != nil {
						return err
					}
				}
			}
		}
		messages, error = a.Dg.ChannelMessages(channel_id, 100, before_id, "", "")
		if error != nil {
			return error
		}
	}
	return nil
}
