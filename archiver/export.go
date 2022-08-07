package archiver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yakabuff/discord-dl/common"
	"github.com/yakabuff/discord-dl/models"
)

func (a Archiver) ExportMessage(m *discordgo.Message, downloadMedia bool, jobId string) error {
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

	var embeds []models.Embed
	var attachments []models.Attachment

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

		embeds = append(embeds, embed)
	}

	for _, i := range m.Attachments {
		attachment := models.NewAttachment(i.ID, m.ID, i.Filename, i.URL, "")

		//Download embed media
		hash, err := common.DownloadFile(i.URL, m.ChannelID, a.Args.MediaLocation, downloadMedia)
		if err != nil {
			log.Error(err)
		}

		attachment.AttachmentHash = hash

		attachments = append(attachments, attachment)
	}
	msg := models.NewMessageJson(id, channel_id, guild_id, timestamp_unix, content, author_id, author_username, reply_to, fmt.Sprintf("%d", editedTimestamp), threadId, embeds, attachments)

	err := WriteMessageJson(msg, jobId)
	return err
}

// We will need to incrementally append to the file.
// It will be impossible to put 1million messages all in a struct
// We will need to decompose the channel into messages and append them individually
// 1) Files will be <channel snowflake_jobid>.json
// 2) Open file. Check if it is empty.
// 3) If not empty, check if there is content. Check json array ([] symbols)
// 4) If empty, construct json document symbols
// 5) If not empty, delete ] symbol.  Add message struct. Delete , symbol.  Re add ] symbol
// file.Seek(0, 2) to go to end of file
func WriteMessageJson(msg models.MessageJson, jobID string) error {
	name := msg.ChannelId + "_" + jobID + ".json"
	var file *os.File
	//If empty, create, else read last char of file. if ] character, assume it is json. delete char, append and re add ].  if ] not found, append [, append msg and append ]

	_, err := os.Stat(name)

	if errors.Is(err, os.ErrNotExist) {

		file, err = os.Create(name)
		if err != nil {
			return err
		}
		defer file.Close()
		//Write [ character
		file.WriteString("[\n")
		//Write message (no trailing comma)
		b, err := json.MarshalIndent(msg, "", "\t")
		if err != nil {
			return err
		}
		file.WriteString(string(b) + "\n")
		//Write ] character
		file.WriteString("]")

	} else {
		file, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			return err
		}
		defer file.Close()
		// If file already exists, we can assume that there are already messages in the file
		// Delete last ] character. If does not exist, return error

		err := deleteLastSquareBracket(file)
		if err != nil {
			log.Error(err)
			return err
		}
		// Write message (prefixed comma)
		b, err := json.MarshalIndent(msg, "", "\t")
		if err != nil {
			return err
		}
		file.WriteString("," + "\n" + string(b) + "\n")
		// Write ] character
		file.WriteString("]")
	}

	return nil
}

//Verify last line ']' and first line '['
//Second line = '{' and second last line = '}'
func verifyExportStructure(f *os.File) bool {
	valid := false
	char := make([]byte, 1)
	f, _ = os.Open(f.Name())
	f.Seek(-1, io.SeekEnd)
	f.Read(char)

	if char[0] == 93 {
		valid = true
	} else {
		valid = false
	}

	f.Seek(-2, io.SeekEnd)
	f.Read(char)

	if char[0] == 125 {
		valid = true
	} else {
		valid = false
	}

	f.Seek(1, io.SeekStart)
	f.Read(char)

	if char[0] == 91 {
		valid = true
	} else {
		valid = false
	}

	f.Seek(2, io.SeekStart)
	f.Read(char)

	if char[0] == 123 {
		valid = true
	} else {
		valid = false
	}

	return valid
}

//Truncate file by 2 bytes. (\n and ])
func deleteLastSquareBracket(file *os.File) error {
	if verifyExportStructure(file) == false {
		return errors.New("Invalid export structure")
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	filesize := stat.Size()

	return os.Truncate(file.Name(), filesize-2)
}
