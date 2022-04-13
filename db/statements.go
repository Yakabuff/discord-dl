package db

import (
	"log"

	"github.com/mattn/go-sqlite3"
	"github.com/yakabuff/discord-dl/models"
)

func (db Db) InsertMessage(m models.Message) error {
	stmt := `
	INSERT INTO messages (message_id, channel_id, guild_id, date, content, sender_id, sender_name, reply_to, edited_timestamp, thread_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := db.DbConnection.Exec(stmt, m.MessageId, m.ChannelId, m.GuildId, m.MessageTimestamp, m.Content, m.SenderId, m.SenderName, m.ReplyTo, m.EditTime, m.ThreadId)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			log.Println("Message already downloaded. Update complete")
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertEdit(m models.Edit) error {
	stmt := `
	INSERT INTO edits (message_id, edit_time, content)
	VALUES ($1, $2, $3)
	`
	_, err := db.DbConnection.Exec(stmt, m.MessageId, m.EditTime, m.Content)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			log.Println("Message already downloaded. Update complete")
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertEmbed(m models.Embed) error {
	stmt := `
	INSERT INTO embeds (message_id, embed_url, embed_title, embed_description, embed_timestamp, embed_thumbnail_url, embed_thumbnail_hash, embed_image_url, embed_image_hash, embed_video_url, embed_video_hash, embed_footer, embed_author_name, embed_author_url, embed_field)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	_, err := db.DbConnection.Exec(stmt, m.MessageId, m.EmbedUrl, m.EmbedTitle, m.EmbedDescription, m.EmbedTimestamp, m.EmbedThumbnailUrl, m.EmbedThumbnailHash, m.EmbedImageUrl, m.EmbedImageHash, m.EmbedVideoUrl, m.EmbedVideoHash, m.EmbedFooter, m.EmbedAuthorName, m.EmbedAuthorUrl, m.EmbedField)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			log.Println("Embed already downloaded.")
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertChannelID(channel string) error {
	stmt := `INSERT INTO channels(channel_id) VALUES($1)`
	_, err := db.DbConnection.Exec(stmt, channel)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			log.Println("Channel ID already registered.")
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertChannelMeta(channel string, guild string, name string, topic string, isThread int) error {
	stmt := `INSERT INTO channels_meta(channel_id, name, topic, guild_id, is_thread) VALUES($1, $2, $3, $4, $5)`
	_, err := db.DbConnection.Exec(stmt, channel, name, topic, guild, isThread)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		log.Println(int(sqliteErr.Code))
		log.Println(int(sqliteErr.ExtendedCode))
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			log.Println("Channel meta already registered.")
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertChannelHistoricalNames(channel string, name string) error {
	stmt := `INSERT INTO channel_names(channel_id, date_renamed, channel_name) VALUES($1, strftime('%s', 'now'), $2)`
	_, err := db.DbConnection.Exec(stmt, channel, name)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			log.Println("Channel name edit already registered.")
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertChannelHistoricalTopic(channel string, topic string) error {
	stmt := `INSERT INTO channel_topics(channel_id, date_renamed, channel_topic) VALUES($1, strftime('%s', 'now'), $2)`
	_, err := db.DbConnection.Exec(stmt, channel, topic)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			log.Println("Channel topic already registered.")
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) UpdateChannelMetadata(channel string, name string, topic string) error {
	log.Println("updating channel metadata")
	stmt := `UPDATE channels_meta SET topic = $1, name = $2 WHERE channel_id = $3`
	_, err := db.DbConnection.Exec(stmt, topic, name, channel)
	return err
}

func (db Db) InsertGuildID(guild string) error {
	stmt := `INSERT INTO guilds(guild_id) VALUES($1)`
	_, err := db.DbConnection.Exec(stmt, guild)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			log.Println("GuildID already registered.")
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertGuildMetadata(guild string, name string, iconHash string, bannerHash string) error {
	stmt := `INSERT INTO guilds_meta(guild_id, icon, banner, name) values($1, $2, $3, $4)`
	_, err := db.DbConnection.Exec(stmt, guild, iconHash, bannerHash, name)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		log.Println(int(sqliteErr.Code))
		log.Println(int(sqliteErr.ExtendedCode))
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			log.Println("Channel meta already registered.")
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) UpdateGuildMetadata(guild string, name string, iconHash string, bannerHash string) error {
	log.Println("updating guild metadata")
	stmt := `UPDATE guilds_meta SET icon = $1, banner = $2, name = $3 WHERE guild_id = $4`
	_, err := db.DbConnection.Exec(stmt, iconHash, bannerHash, name, guild)
	return err
}

func (db Db) InsertGuildHistoricalNames(guild string, guild_name string) error {
	stmt := `INSERT INTO guild_names(guild_id, date_renamed, guild_name) VALUES($1, strftime('%s', 'now'), $2)`
	_, err := db.DbConnection.Exec(stmt, guild, guild_name)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			log.Println("Guild name edit already registered.")
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertGuildHistoricalIcons(guild string, hash string) error {
	stmt := `INSERT INTO guild_icons(guild_id, date_renamed, guild_icon_hash) VALUES($1, strftime('%s', 'now'), $2)`
	_, err := db.DbConnection.Exec(stmt, guild, hash)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			log.Println("Guild icon already exists.")
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertGuildHistoricalBanner(guild string, hash string) error {
	stmt := `INSERT INTO guild_banners(guild_id, date_renamed, guild_banner_hash) VALUES($1, strftime('%s', 'now'), $2)`
	_, err := db.DbConnection.Exec(stmt, guild, hash)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			log.Println("Guild banner already exists.")
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertAttachment(m models.Attachment) error {
	stmt := `
	INSERT INTO attachments (message_id, attachment_id, attachment_filename, attachment_URL, attachment_hash)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err := db.DbConnection.Exec(stmt, m.MessageId, m.AttachmentId, m.AttachmentFilename, m.AttachmentUrl, m.AttachmentHash)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			//1555 == ErrConstraintPrimaryKey
			//19 == ErrConstraint
			log.Println("Attachment already downloaded.")
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}
