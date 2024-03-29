package archiver

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yakabuff/discord-dl/common"
	"github.com/yakabuff/discord-dl/db"
	"github.com/yakabuff/discord-dl/models"
)

func (a Archiver) InsertMessage(m *discordgo.Message, fastUpdate bool, downloadMedia bool) error {
	//Try adding message to DB
	//fast update, unique constraint error -> log and skip message(return fast_update error)
	//fast update, non unique constraint error -> return error
	//non fast update, unique constriant error -> log, continue(do not return anything) and download edits, attachments, embeds
	//non fast update non unique constraint error -> return error

	//Add to edits table if does not exist
	//Add to attachments if not exist
	//Add to embed if not exist

	timestamp, _ := discordgo.SnowflakeTimestamp(m.ID)
	timestamp_unix := timestamp.Unix()
	id := m.ID
	content := m.Content
	author_id := m.Author.ID
	author_username := m.Author.Username
	channel_id := m.ChannelID
	guild_id := m.GuildID
	var reply_to string
	var editedTimestamp int64 = -1
	var threadId string
	//Check if message is a reply
	if m.MessageReference != nil {
		reply_to = (*m).MessageReference.MessageID
	}
	if m.EditedTimestamp != nil {
		editedTimestamp = m.EditedTimestamp.Unix()
	}
	if m.Thread != nil {
		threadId = m.Thread.ID
	}
	msg := models.NewMessage(id, channel_id, guild_id, timestamp_unix, content, author_id, author_username, reply_to, editedTimestamp, threadId)

	errMsg := a.Db.InsertMessage(msg)

	if fastUpdate == true && errMsg != nil {
		//check for unique constraint err. If found, exit program
		if !errors.Is(errMsg, db.UniqueConstraintError) {
			log.Error(errMsg)
			return errMsg
		} else {
			log.Info("Fast update triggered")
			//return fast update error
			return models.FastUpdateError
		}

	} else if fastUpdate == false && errMsg != nil {
		if !errors.Is(errMsg, db.UniqueConstraintError) {
			log.Println(errMsg)
			return errMsg
		}
	}

	//Check if it is edited message.
	//If message is edited, insert edit. Check if uniqueConstraintError
	if m.EditedTimestamp != nil {
		edit := models.NewEdit(id, editedTimestamp, content)
		errEdit := a.InsertEdit(edit)
		if errEdit != nil {
			log.Error(errEdit)
			return errEdit
		}
	}

	//TODO: If hash exists for embed, don't redownload. Sometimes, embed images change eg: github -> num stars goes up in image
	//If it has embed, download embed
	for _, i := range m.Embeds {
		//If image != null, download image, add URL to embed, download
		//If thumbnail != null downlaod thumbnail, add URL to embed, download
		//If video != null download video, add URL to embed, download
		//Iterate through fields for every embed.  Combine fields, seperate with \n

		var fields string
		if i.Fields != nil {
			for _, j := range i.Fields {
				fields = fields + j.Name + "\n" + j.Value + "\n"
			}
		}
		var imageURL string
		if i.Image != nil {
			imageURL = i.Image.URL
		}
		var videoURL string
		if i.Video != nil {
			videoURL = i.Video.URL
		}
		var thumbnailURL string
		if i.Thumbnail != nil {
			thumbnailURL = i.Thumbnail.URL
		}
		var authorName string
		var authorURL string
		if i.Author != nil {
			authorName = i.Author.Name
			authorURL = i.Author.URL
		}

		var footerText string
		if i.Footer != nil {
			footerText = i.Footer.Text
		}
		var dateRetrieved string = fmt.Sprintf("%d", time.Now().Unix())
		embed := models.NewEmbed(m.ID,
			dateRetrieved,
			i.URL,
			i.Title,
			i.Description,
			i.Timestamp,
			thumbnailURL,
			"",
			imageURL,
			"",
			videoURL,
			"",
			footerText,
			authorName,
			authorURL,
			fields,
		)

		//Download embed media
		if i.Image != nil {
			sum, err := common.DownloadFile(i.Image.URL, m.ChannelID, a.Args.MediaLocation, downloadMedia)
			if err != nil {
				var e *fs.PathError
				if errors.As(err, &e) {
					log.Fatal(err)
				}
				log.Error(err)
			}

			embed.EmbedImageHash = sum
		}

		if i.Thumbnail != nil {
			sum, err := common.DownloadFile(i.Thumbnail.URL, m.ChannelID, a.Args.MediaLocation, downloadMedia)
			if err != nil {
				log.Error(err)
			}
			embed.EmbedThumbnailHash = sum
		}
		//Download videos in embeds from discord ONLY.
		if i.Video != nil && strings.HasPrefix(i.Video.URL, "https://cdn.discordapp.com") {
			sum, err := common.DownloadFile(i.Video.URL, m.ChannelID, a.Args.MediaLocation, downloadMedia)
			if err != nil {
				log.Error(err)
			}

			embed.EmbedVideoHash = sum
		}

		errEmbed := a.InsertEmbed(embed)
		if errEmbed != nil {
			log.Error(errEmbed)
			return errEmbed
		}
	}

	for _, i := range m.Attachments {
		attachment := models.NewAttachment(i.ID, m.ID, i.Filename, i.URL, "")

		//Download embed media
		hash, err := common.DownloadFile(i.URL, m.ChannelID, a.Args.MediaLocation, downloadMedia)
		if err != nil {
			log.Error(err)
		}

		attachment.AttachmentHash = hash

		errAttachment := a.InsertAttachment(attachment)
		if errAttachment != nil {
			return errAttachment
		}
	}
	return nil
}

//Note on threads.
//if MessageType=21, this signifies thread top message.  this message has a threads field which is a channel ID
//All messages in the thread has the channelID of that thread and not the channelID the thread is in.
//Note on media
//

func (a Archiver) InsertEdit(edit models.Edit) error {
	errEdit := a.Db.InsertEdit(edit)
	if !errors.Is(errEdit, db.UniqueConstraintError) {
		return errEdit
	}
	return nil
}

func (a Archiver) InsertEmbed(embed models.Embed) error {
	errEmbed := a.Db.InsertEmbed(embed)
	if !errors.Is(errEmbed, db.UniqueConstraintError) {
		return errEmbed
	}
	return nil
}

func (a Archiver) InsertAttachment(attachment models.Attachment) error {
	errAttachment := a.Db.InsertAttachment(attachment)
	if !errors.Is(errAttachment, db.UniqueConstraintError) {
		return errAttachment
	}
	return nil
}

func (a Archiver) ProcessMessages(m *discordgo.Message, fastUpdate bool, downloadMedia bool, id string) error {

	if a.Args.Output != "" && a.Args.Export == false {
		//If database, insertMessage
		err := a.InsertMessage(m, fastUpdate, downloadMedia)
		return err
	} else {
		//If export, export to file
		err := a.ExportMessage(m, downloadMedia, id)
		return err
	}
}
