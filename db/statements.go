package db

import (
	"database/sql"
	"time"

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
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertEmbed(m models.Embed) error {
	stmt := `
	INSERT INTO embeds (message_id, embed_date_retrieved, embed_url, embed_title, embed_description, embed_timestamp, embed_thumbnail_url, embed_thumbnail_hash, embed_image_url, embed_image_hash, embed_video_url, embed_video_hash, embed_footer, embed_author_name, embed_author_url, embed_field)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := db.DbConnection.Exec(stmt, m.MessageId, m.EmbedDateRetrieved, m.EmbedUrl, m.EmbedTitle, m.EmbedDescription, m.EmbedTimestamp, m.EmbedThumbnailUrl, m.EmbedThumbnailHash, m.EmbedImageUrl, m.EmbedImageHash, m.EmbedVideoUrl, m.EmbedVideoHash, m.EmbedFooter, m.EmbedAuthorName, m.EmbedAuthorUrl, m.EmbedField)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
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
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
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
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) UpdateChannelMetadata(channel string, name string, topic string) error {
	stmt := `UPDATE channels_meta SET topic = $1, name = $2 WHERE channel_id = $3`
	_, err := db.DbConnection.Exec(stmt, topic, name, channel)
	return err
}

func (db Db) InsertGuildID(guild string) error {
	stmt := `INSERT INTO guilds(guild_id) VALUES($1)`
	_, err := db.DbConnection.Exec(stmt, guild)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
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
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) UpdateGuildMetadata(guild string, name string, iconHash string, bannerHash string) error {
	stmt := `UPDATE guilds_meta SET icon = $1, banner = $2, name = $3 WHERE guild_id = $4`
	_, err := db.DbConnection.Exec(stmt, iconHash, bannerHash, name, guild)
	return err
}

func (db Db) InsertGuildHistoricalNames(guild string, guild_name string) error {
	stmt := `INSERT INTO guild_names(guild_id, date_renamed, guild_name) VALUES($1, strftime('%s', 'now'), $2)`
	_, err := db.DbConnection.Exec(stmt, guild, guild_name)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
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
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertGuildNames(gid string, guild_name string) error {

	now := time.Now().UnixNano()
	stmt := `INSERT INTO guild_names VALUES($1, $2, $3);`
	_, err := db.DbConnection.Exec(stmt, gid, now/1000000, guild_name)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			//1555 == ErrConstraintPrimaryKey
			//19 == ErrConstraint
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertGuildIcons(gid string, icon_hash string) error {

	now := time.Now().UnixNano()
	stmt := `INSERT INTO guild_icons VALUES($1, $2, $3);`
	_, err := db.DbConnection.Exec(stmt, gid, now/1000000, icon_hash)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			//1555 == ErrConstraintPrimaryKey
			//19 == ErrConstraint
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) InsertGuildBanner(gid string, banner_hash string) error {

	now := time.Now().UnixNano()
	stmt := `INSERT INTO guild_banners VALUES($1, $2, $3);`
	_, err := db.DbConnection.Exec(stmt, gid, now/1000000, banner_hash)

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			//1555 == ErrConstraintPrimaryKey
			//19 == ErrConstraint
			return UniqueConstraintError
		} else {
			return err
		}
	}
	return err
}

func (db Db) UpdateGuildMetaTransaction(gid string) error {

	insert := `INSERT INTO guilds_meta values($1, $2, $3, $4) ON CONFLICT(guild_id) DO UPDATE SET guild_id = $1, name = $4, icon = $2, banner = $3;`

	var name, icon, banner string
	var row *sql.Row

	tx, err := db.DbConnection.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	row = tx.QueryRow("select guild_name from guild_names order by date_renamed DESC LIMIT 1;")
	err = row.Scan(&name)
	if err != nil {
		tx.Rollback()
		return err
	}
	row = tx.QueryRow("select guild_icon_hash from guild_icons order by date_renamed DESC LIMIT 1;")
	err = row.Scan(&icon)
	if err != nil {
		tx.Rollback()
		return err
	}
	row = tx.QueryRow("select guild_banner_hash from guild_banners order by date_renamed DESC LIMIT 1;")
	err = row.Scan(&banner)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(insert, gid, icon, banner, name)

	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err

}
