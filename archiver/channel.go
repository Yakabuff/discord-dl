package archiver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vbauerster/mpb/v7"
	"github.com/yakabuff/discord-dl/common"
	"github.com/yakabuff/discord-dl/db"
	"github.com/yakabuff/discord-dl/job"
)

func (a Archiver) ChannelDownload(channel string, fastUpdate bool, after string, before string, state job.JobState) error {
	// case 1: only after flag.
	// case 2: only before flag
	// case 3: after AND before floag
	// case 4: --fast-update flag
	// case 5: contains no flag
	// log.Println("Downloading channel " + channel)

	c, err := a.Dg.Channel(channel)
	if err != nil {
		if err.(*discordgo.RESTError).Message.Code == 50001 {
			return nil
		}
		log.Println(err)

		return err
	}

	if c == nil {
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
		if fastUpdate {
			err = a.DownloadMessages(after, before, channel, guildID, true, state)
		} else {
			//Both before and after flags
			if before != "" && after != "" {
				err = a.DownloadMessages(after, before, channel, guildID, false, state)
			} else if before == "" && after == "" {
				//Both before and after empty
				err = a.DownloadMessages("", "", channel, guildID, false, state)
			} else {
				//Before empty OR after empty
				if strings.Contains(before, "-") && after == "" {
					err = a.DownloadRangeDate(after, before, channel, guildID, state)
				} else if strings.Contains(after, "-") && before == "" {
					err = a.DownloadRangeDate(after, before, channel, guildID, state)
				} else if !strings.Contains(after, "-") && before == "" {
					err = a.DownloadMessages(after, before, channel, guildID, false, state)
				} else {
					err = a.DownloadMessages(after, before, channel, guildID, false, state)
				}
			}
		}

	} else {
		return errors.New("Channel is not a text channel")
	}

	return err
}

//Convert after and before to Time if applicable
//Convert Time for after and before to Discord Unix
//Run DownloadMessage on before and after unix times
func (a Archiver) DownloadRangeDate(after string, before string, channel_id string, guild_id string, state job.JobState) error {

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
			err := a.DownloadMessages(before_id, after_id, channel_id, guild_id, false, state)
			if err != nil {
				return err
			}
		} else {
			err := a.DownloadMessages(before_id, "", channel_id, guild_id, false, state)
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
			a.DownloadMessages("", after_id, channel_id, guild_id, false, state)
		} else {
			a.DownloadMessages("", "", channel_id, guild_id, false, state)
		}
	}

	return nil
}

func (a Archiver) DownloadMessages(before_id string, after_id string, channel_id string, guild_id string, fast_update bool, state job.JobState) error {
	var _beforeid string = before_id
	var _afterid string = after_id
	messages, error := a.Dg.ChannelMessages(channel_id, 100, before_id, "", "")
	if error != nil {
		log.Error(error)

		return error
	}
	//Start archiving messages
	var in_range bool = true
	for len(messages) != 0 && in_range {
		for _, m := range messages {
			// timestamp, _ := discordgo.SnowflakeTimestamp(m.ID)
			id := m.ID
			// content := m.Content;
			// author_id := m.Author.ID;
			// author_username := m.Author.Username;
			m.GuildID = guild_id
			if after_id != "" {
				if id > after_id {
					// log.Printf("Downloading message %s %s\n", m.Timestamp, m.ID)
					before_id = id
					//insert into db
					err := a.InsertMessage(m, fast_update, a.Args.DownloadMedia)
					if err != nil {
						if errors.Is(err, db.UniqueConstraintError) {
							return nil
						}
						log.Error(err)
						return err
					}

					CalculateChannelProgress(_afterid, _beforeid, m.ID, state.Progress, state.Bar)
					if *state.Status == job.CANCELLED {
						return nil
					}

					//Fetch threads if exist
					if m.Type == 21 {
						log.Println("Thread spotted. Traversing thread: " + m.Thread.ID)
						err := a.DownloadMessages("", "", m.Thread.ID, m.Thread.GuildID, fast_update, state)
						if err != nil {
							return err
						}
					}

				} else {
					in_range = false
					break
				}
			} else {
				// log.Printf("Downloading message %s %s %s\n", m.Timestamp, m.ID, m.ChannelID)
				err := a.InsertMessage(m, fast_update, a.Args.DownloadMedia)
				before_id = id
				if err != nil {
					if errors.Is(err, db.UniqueConstraintError) {
						return nil
					}
					return err
				}

				CalculateChannelProgress(_afterid, _beforeid, m.ID, state.Progress, state.Bar)
				if *state.Status == job.CANCELLED {
					return nil
				}

				if m.Thread != nil {
					// log.Println("Thread spotted. Traversing thread: " + m.Thread.ID)
					a.IndexChannel(m.Thread.ID)
					err := a.DownloadMessages("", "", m.Thread.ID, m.Thread.GuildID, fast_update, state)
					if err != nil {
						return err
					}
				}
			}
		}
		if *state.Status == job.CANCELLED {
			return nil
		}
		messages, error = a.Dg.ChannelMessages(channel_id, 100, before_id, "", "")
		if error != nil {
			return error
		}
	}
	return nil
}

func CalculateChannelProgress(afterID string, beforeID string, currMessageID string, progress *int, bar *mpb.Bar) {
	//If only after empty: after == discordEpoch, before = before.  after = after - after. before = before - after. curr = curr - after
	//If only before empty: after == after, before = curr unix time.  after = after - after. before = before - after. curr = curr - after
	//If both empty..
	//If neither empty: after = after, before = before.  after = after - after. before = before - after. curr = curr - after
	var startTime int64
	var endTime int64
	// if after == 0, assume discord epoch
	if afterID == "" && beforeID != "" {
		endTime = common.DiscordEpoch
		startTime, _ = common.SnowflakeToUnixTime(beforeID)
	} else if afterID != "" && beforeID == "" {
		endTime, _ = common.SnowflakeToUnixTime(afterID)
		startTime = time.Now().Unix()
	} else if afterID == "" && beforeID == "" {
		endTime = common.DiscordEpoch / 1000
		startTime = time.Now().Unix()
	} else {
		endTime, _ = common.SnowflakeToUnixTime(afterID)
		startTime, _ = common.SnowflakeToUnixTime(beforeID)
	}

	currTime, _ := common.SnowflakeToUnixTime(currMessageID)
	//Normalize numbers to calculate percentage
	var quotient float64 = float64(((float64(startTime) - float64(endTime)) - (float64(currTime) - float64(endTime)))) / (float64(startTime) - float64(endTime))
	pct := quotient * 100
	*progress = int(pct)
	bar.SetCurrent(int64(pct))
	// log.Printf("Progress is %d, start time: %d, end time: %d, curr time: %d", *progress, startTime, endTime, currTime)
}
