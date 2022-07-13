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
			return nil
		} else {
			return err
		}
	}
	return err
}

//fix
func (db Db) InsertChannelNames(channel string, name string) error {
	now := time.Now().UnixNano()
	stmt := `INSERT INTO channel_names(channel_id, date_renamed, channel_name) VALUES($1, $2, $3)`
	_, err := db.DbConnection.Exec(stmt, channel, now/1000000, name)

	// if sqliteErr, ok := err.(sqlite3.Error); ok {
	// 	if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
	// 		return nil
	// 	} else {
	// 		return err
	// 	}
	// }
	return err
}

//fix
func (db Db) InsertChannelTopic(channel string, topic string) error {
	now := time.Now().UnixNano()
	stmt := `INSERT INTO channel_topics(channel_id, date_renamed, channel_topic) VALUES($1, $2, $3)`
	_, err := db.DbConnection.Exec(stmt, channel, now/1000000, topic)
	return err
}

func (db Db) InsertGuildID(guild string) error {
	stmt := `INSERT INTO guilds(guild_id) VALUES($1)`
	_, err := db.DbConnection.Exec(stmt, guild)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555 {
			return nil
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

	return err
}

func (db Db) InsertGuildNames(gid string, guild_name string) error {

	now := time.Now().UnixNano()
	stmt := `INSERT INTO guild_names VALUES($1, $2, $3);`
	_, err := db.DbConnection.Exec(stmt, gid, now/1000000, guild_name)

	return err
}

func (db Db) InsertGuildIcons(gid string, icon_hash string) error {

	now := time.Now().UnixNano()
	stmt := `INSERT INTO guild_icons VALUES($1, $2, $3);`
	_, err := db.DbConnection.Exec(stmt, gid, now/1000000, icon_hash)

	return err
}

func (db Db) InsertGuildBanner(gid string, banner_hash string) error {

	now := time.Now().UnixNano()
	stmt := `INSERT INTO guild_banners VALUES($1, $2, $3);`
	_, err := db.DbConnection.Exec(stmt, gid, now/1000000, banner_hash)

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
	row = tx.QueryRow("select guild_name from guild_names where guild_id=$1 order by date_renamed DESC LIMIT 1;", gid)
	err = row.Scan(&name)
	if err != nil {
		tx.Rollback()
		return err
	}
	row = tx.QueryRow("select guild_icon_hash from guild_icons where guild_id=$1 order by date_renamed DESC LIMIT 1;", gid)
	err = row.Scan(&icon)
	if err != nil {
		tx.Rollback()
		return err
	}
	row = tx.QueryRow("select guild_banner_hash from guild_banners where guild_id=$1 order by date_renamed DESC LIMIT 1;", gid)
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

func (db Db) UpdateChannelMetaTransaction(cid string, isThread bool, gid string) error {

	insert := `INSERT INTO channels_meta values($1, $2, $3, $4, $5) ON CONFLICT(channel_id) DO UPDATE SET channel_id = $1, name = $2, topic = $3;`

	var name, topic string
	var row *sql.Row

	tx, err := db.DbConnection.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	row = tx.QueryRow("select channel_name from channel_names where channel_id=$1 order by date_renamed DESC LIMIT 1;", cid)
	err = row.Scan(&name)
	if err != nil {
		tx.Rollback()
		return err
	}
	row = tx.QueryRow("select channel_topic from channel_topics where channel_id=$1 order by date_renamed DESC LIMIT 1;", cid)
	err = row.Scan(&topic)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(insert, cid, name, topic, gid, isThread)

	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err

}
