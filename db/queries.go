package db

import (
	"database/sql"
	"log"

	"github.com/yakabuff/discord-dl/models"
)

func (db Db) GetMessages(guild_id string, channel_id string, last_date int, after bool) (error, *models.Messages) {
	var messages models.Messages
	//use keyset pagination
	//first page: fetch first 100 messages. get date > curr time. keep track of the date of the last message returned
	//second page: fetch second batch of 100. get date > date of the last message returned in the previous batch
	//query messages -> query edits -> query embeds -> query attachments
	var stmt string
	if after == true {
		stmt = `SELECT * FROM messages where channel_id = $1 AND guild_id = $2 AND date > $3 ORDER BY date DESC LIMIT 10`
	} else {
		stmt = `SELECT * FROM messages where channel_id = $1 AND guild_id = $2 AND date < $3 ORDER BY date DESC LIMIT 10`
	}

	rows, err := db.DbConnection.Query(stmt, channel_id, guild_id, last_date)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	for rows.Next() {
		var message models.MessageOut
		err := rows.Scan(
			&message.MessageId,
			&message.ChannelId,
			&message.GuildId,
			&message.MessageTimestamp,
			&message.Content,
			&message.SenderId,
			&message.SenderName,
			&message.ReplyTo,
			&message.EditTime,
			&message.ThreadId)
		if err != nil {
			log.Println(err)
			return err, nil
		}

		err, edits := db.GetEdits(message.MessageId)
		if err != nil {
			log.Println(err)
			return err, nil
		}
		err, embeds := db.GetEmbeds(message.MessageId)
		if err != nil {
			return err, nil
		}
		err, attachments := db.GetAttachments(message.MessageId)
		if err != nil {
			return err, nil
		}

		// addEmbedResourceLink(embeds, message.ChannelId)
		// addAttachmentResourceLink(attachments, message.ChannelId)
		message.Edits = edits
		message.Embeds = embeds
		message.Attachments = attachments
		messages.Messages = append(messages.Messages, message)
	}
	return nil, &messages
}

func (db Db) GetEdits(message_id string) (error, []models.Edit) {
	var edits []models.Edit
	stmt := `SELECT * FROM edits where message_id = $1`
	rows, err := db.DbConnection.Query(stmt, message_id)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		var edit models.Edit
		err := rows.Scan(&edit.MessageId,
			&edit.EditTime,
			&edit.Content)

		if err != nil {
			return err, nil
		}
		edits = append(edits, edit)
	}
	return nil, edits
}

func (db Db) GetEmbeds(message_id string) (error, []models.EmbedOut) {
	var embeds []models.EmbedOut
	stmt := `SELECT * FROM embeds where message_id = $1`
	rows, err := db.DbConnection.Query(stmt, message_id)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	for rows.Next() {
		var embed models.EmbedOut
		err := rows.Scan(
			&embed.MessageId,
			&embed.EmbedUrl,
			&embed.EmbedTitle,
			&embed.EmbedDescription,
			&embed.EmbedTimestamp,
			&embed.EmbedThumbnailUrl,
			&embed.EmbedThumbnailHash,
			&embed.EmbedImageUrl,
			&embed.EmbedImageHash,
			&embed.EmbedVideoUrl,
			&embed.EmbedVideoHash,
			&embed.EmbedFooter,
			&embed.EmbedAuthorName,
			&embed.EmbedAuthorUrl,
			&embed.EmbedField)
		if err != nil {
			return err, nil
		}

		embeds = append(embeds, embed)
	}
	return nil, embeds
}

func (db Db) GetAttachments(message_id string) (error, []models.AttachmentOut) {
	var attachments []models.AttachmentOut
	stmt := `SELECT * FROM attachments where message_id = $1`
	rows, err := db.DbConnection.Query(stmt, message_id)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	for rows.Next() {
		var attachment models.AttachmentOut
		err := rows.Scan(
			&attachment.AttachmentId,
			&attachment.MessageId,
			&attachment.AttachmentFilename,
			&attachment.AttachmentUrl,
			&attachment.AttachmentHash)
		if err != nil {
			return err, nil
		}

		attachments = append(attachments, attachment)
	}
	return nil, attachments
}

func (db Db) GetChannelsFromGuild(guild_id string) ([]models.Channel, error) {
	var channels []models.Channel
	stmt := `SELECT * FROM channels where guild_id = $1 AND is_thread = 0`
	rows, err := db.DbConnection.Query(stmt, guild_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var channel models.Channel
		err := rows.Scan(
			&channel.ChannelID,
			&channel.Name,
			&channel.Topic)

		if err != nil {
			return nil, err
		}

		channels = append(channels, channel)
	}
	return channels, nil
}

func (db Db) GetAllGuilds() ([]models.Guild, error) {
	var guilds []models.Guild
	stmt := `SELECT * FROM guilds`
	rows, err := db.DbConnection.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var guild models.Guild
		err := rows.Scan(
			&guild.IconHash,
			&guild.BannerHash,
			&guild.Name)
		if err != nil {
			return nil, err
		}

		guilds = append(guilds, guild)
	}
	return guilds, nil
}

func (db Db) CheckFieldChanged(tableName string, column string, targetValue string) (bool, error) {
	stmt := `SELECT ` + column + ` FROM ` + tableName + ` ORDER BY date_renamed DESC LIMIT 1`
	log.Println(stmt)
	var changed bool = false
	var rowValue string
	err := db.DbConnection.QueryRow(stmt).Scan(&rowValue)
	if err == sql.ErrNoRows {
		log.Println("row not found")
		return true, nil
	}
	if err != nil {
		log.Println("something wrong")
		return false, err
	}
	log.Println(rowValue + " vs " + targetValue)
	if rowValue != targetValue {
		changed = true
	}
	return changed, err
}
