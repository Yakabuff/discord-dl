package archiver

import (
	"errors"

	"github.com/mattn/go-sqlite3"
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
	errChannel := a.Db.InsertChannelID(channel)

	if errChannel != nil {
		log.Error(errChannel.Error())
		return errChannel
	}
	return nil
}

func (a Archiver) InsertGuildID(guild string) error {
	errGuild := a.Db.InsertGuildID(guild)

	if !errors.Is(errGuild, db.UniqueConstraintError) {
		return errGuild
	}
	return nil
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

	err = a.Db.InsertGuildNames(g.ID, g.Name)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if int(sqliteErr.Code) != 19 && int(sqliteErr.ExtendedCode) != 1811 {
				return err
			}
		}
	}

	if g.ID != "" {
		iconHash, err := common.DownloadFile(g.IconURL(), g.ID, a.Args.MediaLocation, true)
		if err != nil {
			log.Error(err)
		}

		//Insert and download guild icon if different.  If g.Icon != select icon_hash from guild_icon ORDER BY date ASC LIMIT 1

		err = a.Db.InsertGuildIcons(g.ID, iconHash)
		if err != nil {
			var sqliteErr sqlite3.Error
			if errors.As(err, &sqliteErr) {
				if int(sqliteErr.Code) != 19 && int(sqliteErr.ExtendedCode) != 1811 {
					return err
				}
			}
		}
	}
	var bannerHash string
	if g.BannerURL() != "" {
		bannerHash, err = common.DownloadFile(g.BannerURL(), g.ID, a.Args.MediaLocation, true)
		if err != nil {
			log.Error(err)
		}
		//Insert and download guild banner if different
	}
	err = a.Db.InsertGuildBanner(g.ID, bannerHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if int(sqliteErr.Code) != 19 && int(sqliteErr.ExtendedCode) != 1811 {
				return err
			}
		}
	}

	//Update guild with new metadata (if any)
	err = a.Db.UpdateGuildMetaTransaction(g.ID)
	return err
}

//Index channel metadata: topic, name, guild it belongs to, channel type
func (a Archiver) IndexChannel(channel string) error {

	c, err := a.Dg.Channel(channel)
	if err != nil {
		log.Error(err)
		return err
	}

	//Insert channel ID
	err = a.InsertChannelID(c.ID)
	if err != nil {
		log.Error(err)
		return err
	}

	err = a.InsertGuildID(c.GuildID)
	if err != nil {
		log.Error(err)
		return err
	}

	err = a.IndexGuild(c.GuildID)
	if err != nil {
		log.Error(err)
		return err
	}

	err = a.Db.InsertChannelNames(c.ID, c.Name)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if int(sqliteErr.Code) != 19 && int(sqliteErr.ExtendedCode) != 1811 {
				log.Error(err)
				return err
			}
		}
	}

	err = a.Db.InsertChannelTopic(c.ID, c.Topic)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if int(sqliteErr.Code) != 19 && int(sqliteErr.ExtendedCode) != 1811 {
				log.Error(err)
				return err
			}
		}
	}
	err = a.Db.UpdateChannelMetaTransaction(c.ID, c.IsThread(), c.GuildID)
	if err != nil {
		log.Error(err)
		return err
	}

	return err
}
