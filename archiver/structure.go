package archiver

import (
	"errors"
	"log"

	"github.com/yakabuff/discord-dl/common"
	"github.com/yakabuff/discord-dl/db"
)

//get channel metadata

//get guild metadasta

//could there be a situation where 2 jobs are fetching metadata at same time. both query channel names, see last name is different and date is different too
//both try to insert same name with their current time (might be slightly different but still duplicates)
//make get channel-metadata a job.. when you make a get_channel job, also make a get_channel_metadata job. don't combine them.
//then limit the number of
//idea: desconstruct jobs by type: guild/structure/channel but also guild id, channel id...  guilds should be deconstructed into channels. the only guild specific thing guild meta data

// queue: {[channel1, channel1], [channel2, channel2, channel2]}.  then process first element of every array in queue concurrently?

func (a Archiver) InsertChannelID(channel string) error {
	log.Println("inserting channel ID")

	errChannel := a.Db.InsertChannelID(channel)
	if !errors.Is(errChannel, db.UniqueConstraintError) {
		log.Println(errChannel)
		return errChannel
	}
	return nil
}

func (a Archiver) InsertChannelMetadata(channel string, guild string, name string, topic string, isThread bool) error {
	log.Println("inserting channel meta")
	var t int
	if isThread {
		t = 1
	}
	errChannel := a.Db.InsertChannelMeta(channel, guild, name, topic, t)
	if !errors.Is(errChannel, db.UniqueConstraintError) {
		log.Println(errChannel)
		return errChannel
	}
	return nil
}

func (a Archiver) InsertGuildMetadata(guild string, name string, iconHash string, bannerHash string) error {
	log.Println("inserting guild meta")
	errChannel := a.Db.InsertGuildMetadata(guild, name, iconHash, bannerHash)
	if !errors.Is(errChannel, db.UniqueConstraintError) {
		log.Println(errChannel)
		return errChannel
	}
	return nil
}

func (a Archiver) InsertChannelName(channel string, channelName string) error {
	log.Println("Inserting channel name")
	changed, err := a.CheckFieldChanged("channel_names", "channel_name", channelName)
	if err != nil {
		return err
	}
	if changed {
		errChannel := a.Db.InsertChannelHistoricalNames(channel, channelName)
		if !errors.Is(errChannel, db.UniqueConstraintError) {
			return errChannel
		}
	}
	return nil
}

func (a Archiver) InsertChannelTopic(channel string, topic string) error {
	log.Println("inserting channel topic")
	changed, err := a.CheckFieldChanged("channel_topics", "channel_topic", topic)
	if err != nil {
		return err
	}
	if changed {
		errChannel := a.Db.InsertChannelHistoricalTopic(channel, topic)
		if !errors.Is(errChannel, db.UniqueConstraintError) {
			return errChannel
		}
	}
	return nil
}

func (a Archiver) InsertGuildID(guild string) error {
	log.Println("Inserting guild")
	errGuild := a.Db.InsertGuildID(guild)
	log.Println(errGuild)
	log.Println(guild)
	if !errors.Is(errGuild, db.UniqueConstraintError) {
		return errGuild
	}
	return nil
}

func (a Archiver) InsertGuildName(guild string, guildName string) error {
	log.Println("Inserting guild name")
	changed, err := a.CheckFieldChanged("guild_names", "guild_name", guildName)
	if err != nil {
		return err
	}
	if changed {
		errGuild := a.Db.InsertGuildHistoricalNames(guild, guildName)
		if !errors.Is(errGuild, db.UniqueConstraintError) {
			return errGuild
		}
	}
	return nil
}

func (a Archiver) InsertGuildHistoricalIcons(guild string, hash string) error {
	log.Println("inserting guild icons")
	changed, err := a.CheckFieldChanged("guild_icons", "guild_icon_hash", hash)
	if err != nil {
		return err
	}
	if changed {
		err := a.Db.InsertGuildHistoricalIcons(guild, hash)
		if !errors.Is(err, db.UniqueConstraintError) {
			return err
		}
	}
	return nil
}

func (a Archiver) InsertGuildHistoricalBanner(guild string, hash string) error {
	log.Println("inserting guild banner")
	changed, err := a.CheckFieldChanged("guild_banners", "guild_banner_hash", hash)
	if err != nil {
		return err
	}
	if changed {
		err := a.Db.InsertGuildHistoricalBanner(guild, hash)
		if !errors.Is(err, db.UniqueConstraintError) {
			return err
		}
	}
	return nil
}

func (a Archiver) CheckFieldChanged(tableName string, column string, targetValue string) (bool, error) {
	log.Println("Checking if field changed " + targetValue)
	changed, err := a.Db.CheckFieldChanged(tableName, column, targetValue)
	log.Println(changed)
	if err != nil {
		log.Println(err)
	}
	return changed, err
}

//Index guild metadata: icon, name, banner
func (a Archiver) IndexGuild(guild string) error {
	g, err := a.Dg.Guild(guild)
	if err != nil {
		return err
	}
	//Insert guild ID
	err = a.InsertGuildID(g.ID)
	if err != nil {
		return err
	}
	//Insert guild name if different
	//make sure latest guild name != g.Name

	err = a.InsertGuildName(g.ID, g.Name)
	if err != nil {
		return err
	}
	var iconHash string
	var bannerHash string
	if g.ID != "" {
		iconHash, err = common.DownloadFile(g.IconURL(), g.ID, a.Args.MediaLocation)
		if err != nil {
			log.Println(err)
		}
		if iconHash != g.Icon {
			log.Println("Mismatch icon hashes")
			log.Println(iconHash + " vs " + g.Icon)
		}
		//Insert and download guild icon if different.  If g.Icon != select icon_hash from guild_icon ORDER BY date ASC LIMIT 1

		err = a.InsertGuildHistoricalIcons(g.ID, iconHash)
		if err != nil {
			return err
		}
	}
	if g.Banner != "" {
		bannerHash, err := common.DownloadFile(g.BannerURL(), g.ID, a.Args.MediaLocation)
		if err != nil {
			log.Println(err)
		}
		if bannerHash != g.Banner {
			log.Println("mismatch banner hash")
		}
		//Insert and download guild banner if different

		err = a.InsertGuildHistoricalBanner(g.ID, bannerHash)
		if err != nil {
			return err
		}
	}
	err = a.InsertGuildMetadata(g.ID, g.Name, iconHash, bannerHash)
	if err != nil {
		return err
	}
	//Update guild with new metadata (if any)
	err = a.Db.UpdateGuildMetadata(g.ID, g.Name, iconHash, bannerHash)
	return err
}

//Index channel metadata: topic, name, guild it belongs to, channel type
func (a Archiver) IndexChannel(channel string) error {
	log.Println("Indexing channel")
	c, err := a.Dg.Channel(channel)
	if err != nil {
		return err
	}

	//Insert channel ID
	err = a.InsertChannelID(c.ID)
	if err != nil {
		return err
	}

	err = a.InsertGuildID(c.GuildID)
	if err != nil {
		return err
	}

	//Insert name if different
	err = a.InsertChannelName(c.ID, c.Name)
	if err != nil {
		return err
	}
	//Insert topic if different
	err = a.InsertChannelTopic(c.ID, c.Topic)
	if err != nil {
		return err
	}
	err = a.InsertChannelMetadata(c.ID, c.GuildID, c.Name, c.Topic, c.IsThread())
	if err != nil {
		return err
	}
	//update current channel metadata if applicable
	err = a.Db.UpdateChannelMetadata(c.ID, c.Name, c.Topic)
	return err
}
