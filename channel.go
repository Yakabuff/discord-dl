package main
import(
	"github.com/bwmarrin/discordgo"
	"fmt"
	 "strings"
	"errors"
	"strconv"
	// "time"
	"log"
)

func channel_download(dg *discordgo.Session, a args) error{
	// case 1: only after flag. 
	// case 2: only before flag
	// case 3: after AND before floag
	// case 4: --fast-update flag
	// case 5: contains no flag

	c, err := dg.Channel(a.channel)
	guildID := c.GuildID
	if err != nil{
		return err
	}

	if c.Type == discordgo.ChannelTypeGuildText{
			
		if a.fast_update == true{
			fmt.Println("Fast update")
		}else{
			//Both before and after flags
			if a.before != "" && a.after != ""{
				if strings.Contains(a.before, "-") && strings.Contains(a.before, "-"){
					download_range_date(dg, a.after, a.before, a.channel, guildID);
				}else{
					download_messages(dg, a.after, a.before, a.channel, guildID);
				}
			}else if a.before == "" && a.after == ""{
				//Both before and after empty
				download_messages(dg, "","", a.channel, guildID);
			}else{
				//Before empty OR after empty
				if strings.Contains(a.before, "-") && a.after == ""{
					download_range_date(dg, a.after, a.before, a.channel, guildID);
				}else if strings.Contains(a.after, "-") && a.before == ""{
					download_range_date(dg, a.after, a.before, a.channel, guildID);
				}else if !strings.Contains(a.after, "-") && a.before == ""{
					download_messages(dg, a.after, a.before, a.channel, guildID);
				}else{
					download_messages(dg, a.after, a.before, a.channel, guildID);
				}
			}
		}

	}else{
		return errors.New("Channel is not a text channel")
	}

	return nil;
}

func download_range_date(dg *discordgo.Session, after string, before string, channel_id string, guild_id string) error {

	if before != ""{
		before_time, _ := DateToTime(before)
		//Generate initial snowflake message ID position
		bt := ((before_time.Unix()+1)*1000 - 1420070400000) << 22
		before_id := strconv.Itoa(int(bt))

		if after != ""{
			after_time, _ := DateToTime(after)
			at := ((after_time.Unix()-1)*1000 - 1420070400000) << 22
			after_id := strconv.Itoa(int(at))
			download_messages(dg, before_id, after_id, channel_id, guild_id)
		}else{
			download_messages(dg, before_id, "", channel_id, guild_id)
		}
		
	}else{
		if after != ""{
			after_time, _ := DateToTime(after)
			fmt.Println(after_time)
			at := ((after_time.Unix()-1)*1000 - 1420070400000) << 22
			after_id := strconv.Itoa(int(at))
			download_messages(dg, "", after_id, channel_id, guild_id)
		}else{
			download_messages(dg, "", "", channel_id, guild_id)
		}
	}

	return nil;
}

func download_messages(dg *discordgo.Session, before_id string, after_id string, channel_id string, guild_id string) error{
	messages, error := dg.ChannelMessages(channel_id, 100, before_id, "", "")
	var in_range bool = true
	for len(messages) != 0 && in_range{
		for _, m := range messages{
			timestamp, _ := discordgo.SnowflakeTimestamp(m.ID)
			id := m.ID
			content := m.Content;
			author_id := m.Author.ID;
			author_username := m.Author.Username;
			m.GuildID = guild_id
			if after_id != ""{
				if id > after_id {
					log.Printf("Downloading message %s %s %s %s %s\n", timestamp, id, content, author_id, author_username);
					before_id = id;
					//insert into db
					addMessage(db, m)
				}else{
					in_range = false;
					break;
				}
			}else{
				log.Printf("Downloading message %s %s %s %s %s\n", timestamp, id, content, author_id, author_username);
				addMessage(db, m)
				before_id = id;
			}
		}
		messages, error = dg.ChannelMessages(channel_id, 100, before_id, "", "")
		if error != nil{
			return error;
		}
	}
	return nil
}

