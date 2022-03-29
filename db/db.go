// +build test
package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yakabuff/discord-dl/models"

	// "github.com/bwmarrin/discordgo"
	"errors"

	"github.com/mattn/go-sqlite3"

	// "path/filepath"
	"strings"
	// "github.com/yakabuff/discord-dl/models"
)

const POSTGRES = "postgres"
const SQLITE = "sqlite3"

var UniqueConstraintError = errors.New("Unique constraint error")

type Db struct {
	DbConnection *sql.DB
}

func Init_db(path string) (*Db, error) {
	var err error
	//Check if DB exists
	if path == "" {
		//If path empty, fallback and check default db location
		_, err = os.Stat("archive.db")
	} else {
		_, err = os.Stat(path)
	}

	driver := determineDbType(path)
	log.Println(path)
	var dbConn *sql.DB
	var file *os.File
	if err == nil {
		//Exists
		if path == "" {
			dbConn, err = sql.Open(driver, "archive.db")
		} else {
			dbConn, err = sql.Open(driver, path)
		}
		if err != nil {
			return nil, err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		if path == "" {
			path = "archive.db"
			file, err = os.Create(path)
		} else {
			file, err = os.Create(path) // Create SQLite file
		}
		if err != nil {
			log.Println("could not create db file")
			log.Fatal(err.Error())
		}
		file.Close()
		dbConn, err = sql.Open("sqlite3", path)
		if err != nil {
			return nil, err
		}
		createTable(dbConn)
		//*message_id | channel_id | guild| | date | content | media | sender_id | reply_to //
		// 234234242  | 23489353   | 324242 | 1231 |asdfasdfs | <urL> | 234242 | 234756//
	} else {
		//Panic
		log.Fatal(err.Error())
	}
	db := Db{DbConnection: dbConn}
	return &db, err
}
func determineDbType(path string) string {
	if strings.HasPrefix(path, "postgres://") {
		return POSTGRES
	}
	return SQLITE
}
func createTable(db *sql.DB) {
	createMessages := `CREATE TABLE messages (
		"message_id" TEXT NOT NULL PRIMARY KEY,		
		"channel_id" TEXT NOT NULL,
		"guild_id" TEXT,
		"date" INTEGER NOT NULL,
		"content" TEXT NOT NULL,
		"sender_id"	TEXT NOT NULL,
		"sender_name" TEXT NOT NULL,
		"reply_to" TEXT,
		"edited_timestamp" INTEGER,
		"thread_id" TEXT
	);`

	createAttachments := `CREATE TABLE attachments (
		"message_id" TEXT NOT NULL,  
		"attachment_id" TEXT NOT NULL PRIMARY KEY,
		"attachment_filename" TEXT NOT NULL,
		"attachment_URL" TEXT NOT NULL,
		"attachment_hash" TEXT NOT NULL
	);`

	createEdits := `CREATE TABLE edits (
		"message_id" TEXT,  
		"edit_time" INTEGER,
		"content" TEXT,
		PRIMARY KEY(message_id, edit_time, content)
	);`

	//field name: name1\nname2\nname3\nnam4
	//body  body1\nbody2\body3 etc.
	createEmbeds := `CREATE TABLE embeds (
		"message_id" TEXT NOT NULL,
		"embed_url" TEXT,
		"embed_title" TEXT,
		"embed_description" TEXT,  
		"embed_timestamp" TEXT,
		"embed_thumbnail_url" TEXT,
		"embed_thumbnail_hash" TEXT,  
		"embed_image_url" TEXT,
		"embed_image_hash" TEXT,
		"embed_video_url" TEXT,
		"embed_video_hash" TEXT,
		"embed_footer" TEXT,
		"embed_author_name" TEXT,
		"embed_author_url" TEXT, 
		"embed_field" TEXT,
		UNIQUE(message_id, 
			embed_url, 
			embed_title, 
			embed_description,
			embed_timestamp,
			embed_thumbnail_url,
			embed_thumbnail_hash,
			embed_image_url,
			embed_image_hash,
			embed_video_url,
			embed_video_hash,
			embed_footer,
			embed_author_name,
			embed_author_url,
			embed_field
		)
	);`

	log.Println("Create messages table...")
	statement, err := db.Prepare(createMessages) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("message table created")

	log.Println("Create attachment table...")
	statement_attachment, err := db.Prepare(createAttachments) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement_attachment.Exec() // Execute SQL Statements
	log.Println("Attachment table created")

	log.Println("Create edits table...")
	statement_edits, err := db.Prepare(createEdits) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement_edits.Exec() // Execute SQL Statements
	log.Println("Edits table created")

	log.Println("Create embeds table...")
	statement_embeds, err := db.Prepare(createEmbeds) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement_embeds.Exec() // Execute SQL Statements
	log.Println("Embeds table created")

	message_index := `CREATE INDEX messages_index on messages(message_id, channel_id, guild_id, sender_id, date)`
	log.Println("Create messages index...")
	statement_message_index, err := db.Prepare(message_index) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement_message_index.Exec()

	attachments_index := `CREATE INDEX attachments_index on attachments(message_id, attachment_id)`
	log.Println("Create attachments index...")
	statement_attachments_index, err := db.Prepare(attachments_index) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement_attachments_index.Exec()

	edits_index := `CREATE INDEX edits_index on edits(message_id)`
	log.Println("Create attachments index...")
	statement_edits_index, err := db.Prepare(edits_index) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement_edits_index.Exec()

	embeds_index := `CREATE INDEX embeds_index on embeds(message_id)`
	log.Println("Create embeds index...")
	statement_embeds_index, err := db.Prepare(embeds_index) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement_embeds_index.Exec()
}

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
