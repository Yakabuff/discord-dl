package main

import(
	_ "github.com/mattn/go-sqlite3"
	"os"
	"log"
	"database/sql"
	"github.com/bwmarrin/discordgo"
	"errors"
	"github.com/mattn/go-sqlite3"
)

func init_db() (*sql.DB, error){
	var err error;
	//Check if DB exists
	_, err = os.Stat("bigbrother.db")
	var db *sql.DB;
	if err == nil {
        //Exists
		db, err = sql.Open("sqlite3", "./bigbrother.db")
		if err != nil {
			return nil, err;
		}
    }else if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create("bigbrother.db") // Create SQLite file
		if err != nil {
			log.Println("could not create db file lul")
			log.Fatal(err.Error())
		}
		file.Close()
		db, _ = sql.Open("sqlite3", "./bigbrother.db")
		createTable(db)
		//*message_id | channel_id | guild| | date | content | media | sender_id | reply_to //
		// 234234242  | 23489353   | 324242 | 1231 |asdfasdfs | <urL> | 234242 | 234756//
    }else{
		//Panic
		log.Fatal(err.Error())
	}
	return db, err
}

func createTable(db *sql.DB){
	createMessages := `CREATE TABLE messages (
		"message_id" TEXT NOT NULL PRIMARY KEY,		
		"channel_id" TEXT NOT NULL,
		"guild_id" TEXT NOT NULL,
		"date" INTEGER NOT NULL,
		"content" TEXT NOT NULL,
		"sender_id"	TEXT NOT NULL,
		"sender_name" TEXT NOT NULL,
		"reply_to" TEXT,
		"edited_timestamp" INTEGER
	);` 

	createAttachments:= `CREATE TABLE attachments (
		"message_id" TEXT NOT NULL,  
		"attachment_id" TEXT NOT NULL PRIMARY KEY,
		"attachment_filename" TEXT NOT NULL,
		"attachment_URL" TEXT NOT NULL,
		"attachment_hash" TEXT NOT NULL
	);`

	createEdits := `CREATE TABLE edits (
		"message_id" TEXT NOT NULL,  
		"edit_time" INTEGER NOT NULL,
		"content" TEXT NOT NULL
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
		"embed_footer" TEXT,
		"embed_author_name" TEXT,
		"embed_author_url" TEXT, 
		"embed_field" TEXT
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

func addMessage(db *sql.DB, m *discordgo.Message, fast_update bool) error{
	timestamp, _ := discordgo.SnowflakeTimestamp(m.ID)
	timestamp_unix := timestamp.Unix()
	id := m.ID
	content := m.Content;
	author_id := m.Author.ID;
	author_username := m.Author.Username;
	channel_id := m.ChannelID;
	guild_id := m.GuildID;
	edited_timestamp, err := m.EditedTimestamp.Parse()
	var reply_to string;
	if m.MessageReference != nil{
		reply_to = (*m).MessageReference.MessageID;
	}
	var media []*discordgo.MessageAttachment;
	if m.Attachments != nil{
		media = m.Attachments;
	}
	log.Printf("Downloading message %s %s %s %s %s %v \n", timestamp, id, content, author_id, author_username, len(m.Embeds));
	stmt := `
	INSERT INTO messages (message_id, channel_id, guild_id, date, content, sender_id, sender_name, reply_to, edited_timestamp)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	//fast update, unique constraint error -> log and skip message(return nil)
	//fast update, non unique constraint error -> return error
	//non fast update, unique constriant error -> log, continue(do not return anything) and download edits, attachments, embeds
	//non fast update non unique constraint error -> return error

	if err != nil{
		//If not edited, insert message in db IF it does not already exist
		_, err := db.Exec(stmt, id, channel_id, guild_id, timestamp_unix, content, author_id, author_username, reply_to, "")

		//check if fast_udpate is true. if true and error is unique constraint error, exit out.
		if (fast_update == true && err != nil){
			//check for unique constraint err. If found, exit program( return nil)
			if sqliteErr, ok := err.(sqlite3.Error); ok {
				if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555{
					log.Println("Message already downloaded. Update complete")
					return UniqueConstraintError
				}else{
					return err
				}
			}
		}else if fast_update ==false && err != nil{
			//check for unique constraint err. If found, exit program( return nil)
			if sqliteErr, ok := err.(sqlite3.Error); ok {
				if !(sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555){
					return err
				}
			}
		}

		err = addEmbed(db, m)
		if err != nil{
			log.Println("ERROR: Could not insert embed")
			return err
		}

	}else{

		//If message is edited, check in edit table if content and edit_time match message content and edit_time
		//If edit does not exist, enter message in DB. 
		log.Println("Edit detected")
		edited_timestamp_unix := edited_timestamp.Unix()

		//Insert edited message in messages table just in case it isn't already archived
		_, err := db.Exec(stmt, id, channel_id, guild_id, timestamp_unix, content, author_id, author_username, reply_to, edited_timestamp_unix)

		//check if fast_udpate is true. if true and error is unique constraint error, exit out.
		if (fast_update == true && err != nil){
			//check for unique constraint err. If found, exit program( return nil)
			if sqliteErr, ok := err.(sqlite3.Error); ok {
				if sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555{
					log.Println("Message already downloaded. Update complete")
					return UniqueConstraintError
				}else{
					return err
				}
			}
		}else if fast_update ==false && err != nil{
			//check for unique constraint err. If found, exit program( return nil)
			if sqliteErr, ok := err.(sqlite3.Error); ok {
				if !(sqliteErr.Code == 19 && sqliteErr.ExtendedCode == 1555){
					return err
				}
			}
		}

		err = addEdit(db, id, edited_timestamp_unix, content)
		if err != nil{
			log.Println("ERROR: Could not insert edit")
			log.Println(err)
			return err;
		}
		err = addEmbed(db, m)
		if err != nil{
			log.Println("ERROR: Could not insert embed")
			log.Println(err)
			return err;
		}
	}

	err = addAttachment(db, m, media)
	if err != nil{
		log.Println("ERROR: Could not insert attachment (already exists?)")
		return err;
	}
	return nil;
}


func addEdit(db *sql.DB, message_id string, edit_time int64, content string) error{
	//query for edit with same time and string. If does not exist, add to table
	var exists bool
	var count int
	sel_stmt :=`
	SELECT COUNT(*)
	FROM edits
	WHERE message_id = $1 AND edit_time = $2 AND content = $3
	`
	err := db.QueryRow(sel_stmt, message_id, edit_time, content).Scan(&count)
	if err != nil{
		log.Fatal(err)
	}

	if count > 0{
		log.Println("Edit found. Skipping...")
		exists = true
	}else{
		log.Println("Edit not found. Adding to DB")
		exists = false
	}

	//If edit does not exist, add to table
	if !exists {
		edit_stmt :=`
		INSERT INTO edits (message_id, edit_time, content)
		VALUES ($1, $2, $3)
		`

		_, err := db.Exec(edit_stmt, message_id, edit_time, content)
		log.Println("Inserting edit")
		if err != nil{
			return err
		}
	}
	return nil
}

func addAttachment(db *sql.DB, m *discordgo.Message, media []*discordgo.MessageAttachment) error{

	if len(media) != 0{

		stmt := `
		INSERT OR IGNORE INTO attachments (message_id, attachment_id, attachment_filename, attachment_URL, attachment_hash)
		VALUES ($1, $2, $3, $4, $5)
		`
		for _, i := range media{
			log.Println("Downloading attachment "+ i.URL)
			//Download attachment
			err, hash := DownloadFile(i.URL, m.ChannelID);
			if err != nil{
				log.Println("Download failed:" + i.URL)
				return err
			}

			_, err = db.Exec(stmt, m.ID, i.ID, i.Filename, i.URL, hash)
			if err != nil{
				return err
			}
		}
	}
	return nil;
}

func addEmbed(db *sql.DB, m *discordgo.Message) error {


	if len(m.Embeds) != 0{

		var URL string 
		var title string
		var description string 
		var timestamp string 
		var thumbnail_url string 
		var image_url string
		var footer_text string
		var author_name string
		var author_url string

		var exists bool
		var count int
	
		sel_stmt :=`
		SELECT COUNT(*)
		FROM embeds
		WHERE message_id = $1 AND 
		embed_url = $2 AND 
		embed_title = $3 AND
		embed_description = $4 AND 
		embed_timestamp = $5 AND 
		embed_thumbnail_url = $6 AND 
		embed_image_url = $7 AND 
		embed_footer = $8 AND
		embed_author_name = $9 AND
		embed_author_url = $10 AND
		embed_field = $11
		`

		stmt := `
		INSERT INTO embeds (
		"message_id",
		"embed_url",
		"embed_title",
		"embed_description",  
		"embed_timestamp",
		"embed_thumbnail_url",
		"embed_thumbnail_hash",    
		"embed_image_url",
		"embed_image_hash",
		"embed_footer",
		"embed_author_name",
		"embed_author_url",
		"embed_field")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`
		for _, i := range m.Embeds{
			var field string

			if len(i.Fields) != 0{
				for _, j := range i.Fields{
					field = field + j.Name + "\n" + j.Value + "\n"
				}
			}

			URL = i.URL
			title = i.Title
			description = i.Description
			timestamp = i.Timestamp

			if i.Thumbnail != nil{
				thumbnail_url = i.Thumbnail.URL
			}
			if i.Image != nil{
				image_url = i.Image.URL
			}
			if i.Footer != nil{
				footer_text = i.Footer.Text
			}
			if i.Author != nil{
				author_name = i.Author.Name
				author_url = i.Author.URL
			}

			err := db.QueryRow(sel_stmt, 
				m.ID, 
				URL, 
				title,
				description, 
				timestamp, 
				thumbnail_url, 
				image_url, 
				footer_text, 
				author_name, 
				author_url,
				field).Scan(&count)

			if err != nil{
				log.Println(err)
				return err
			}
			if count > 0{
				log.Println("Embed found. Skipping...")
				exists = true
			}else{
				log.Println("Embed not found. Adding to DB")
				exists = false
			}

			if !exists{
				var hash_thumbnail string
				var hash_image string

				if i.Image != nil{
					//Download thumbnail/image. File format... hash
					err, hash_image = DownloadFile(i.Image.URL, m.ChannelID);
					if err != nil{
						log.Println("Download failed:" + i.Image.URL)
					}
				}
				if i.Thumbnail != nil{
					err, hash_thumbnail = DownloadFile(i.Thumbnail.URL, m.ChannelID);
					if err != nil{
						log.Println("Download failed:" + i.Thumbnail.URL)
					}
				}


				_, err3 := db.Exec(stmt, 
					m.ID,
					URL,
					title,
					description, 
					timestamp, 
					thumbnail_url, 
					hash_thumbnail,
					image_url,
					hash_image,
					footer_text,
					author_name,
					author_url,
					field,
				)
				if err3 != nil{
					log.Println(err3)
					return err3
				}

			}
		}
	}
	return nil;
}

var UniqueConstraintError =  errors.New("Unique constraint error")

func getMessages(db *sql.DB, guild_id string, channel_id string, last_date int) (error, *Messages){
	var messages Messages;
	//use keyset pagination
	//first page: fetch first 100 messages. get date > curr time. keep track of the date of the last message returned
	//second page: fetch second batch of 100. get date > date of the last message returned in the previous batch
	//query messages -> query edits -> query embeds -> query attachments
	stmt := `SELECT * FROM messages where channel_id = $1 AND guild_id = $2 AND date < $3 ORDER BY date DESC LIMIT 100` 
	rows, err := db.Query(stmt, channel_id, guild_id, last_date)

	if err != nil {
        return err, nil;
    }
    defer rows.Close()

	for rows.Next(){
		var message Message
		err := rows.Scan(
						&message.Message_id,
						&message.Channel_id,
						&message.Guild_id,
						&message.Message_timestamp,
						&message.Content,
						&message.Sender_id,
						&message.Sender_name,
						&message.Reply_to,
						&message.Edit_time); 
		if err != nil {
			log.Println(err)
            return err, nil;
        }
		err, edits := getEdits(db, message.Message_id)
		if err != nil{
			log.Println(err)
			return err, nil;
		}
		err, embeds := getEmbeds(db, message.Message_id)
		if err != nil{
			return err, nil;
		}
		err, attachments := getAttachments(db, message.Message_id)
		if err != nil{
			return err, nil;
		}

		message.Edits = *edits;
		message.Embeds = *embeds;
		message.Attachments = *attachments;
		messages.Messages = append(messages.Messages, message)
	}

	return nil, &messages;
}

func getEdits(db *sql.DB, message_id string) (error, *[]Edit){
	var edits []Edit;
	stmt := `SELECT * FROM edits where message_id = $1` 
	rows, err := db.Query(stmt, message_id)
	if err != nil {
        return err, nil;
    }
    defer rows.Close()
	for rows.Next(){
		var edit Edit
		err := rows.Scan(&edit.Message_id,
						&edit.Content,
						&edit.Edit_time); 
		if err != nil {
            return err, nil;
        }
		edits = append(edits, edit)
	}
	return nil, &edits
}

func getEmbeds(db *sql.DB, message_id string)(error, *[]Embed){
	var embeds []Embed;
	stmt := `SELECT * FROM embeds where message_id = $1` 
	rows, err := db.Query(stmt, message_id)
	if err != nil {
        return err, nil;
    }
    defer rows.Close()

	for rows.Next(){
		var embed Embed
		err := rows.Scan(
			&embed.Message_id, 
			&embed.Embed_url, 
			&embed.Embed_title, 
			&embed.Embed_description, 
			&embed.Embed_timestamp, 
			&embed.Embed_thumbnail_url, 
			&embed.Embed_thumbnail_hash,
			&embed.Embed_image_url, 
			&embed.Embed_image_hash, 
			&embed.Embed_footer, 
			&embed.Embed_author_name, 
			&embed.Embed_author_url, 
			&embed.Embed_field); 
		if err != nil {
            return err, nil;
        }

		embeds = append(embeds, embed)
	}
	return nil, &embeds
}

func getAttachments(db *sql.DB, message_id string) (error, *[]Attachment){
	var attachments []Attachment;
	stmt := `SELECT * FROM attachments where message_id = $1` 
	rows, err := db.Query(stmt, message_id)
	if err != nil {
        return err, nil;
    }
    defer rows.Close()

	for rows.Next(){
		var attachment Attachment
		err := rows.Scan(
						&attachment.Attachment_id,
						&attachment.Message_id,
						&attachment.Attachment_filename,
						&attachment.Attachment_url,
						&attachment.Attachment_hash); 
		if err != nil {
            return err, nil;
        }

		attachments = append(attachments, attachment)
	}
	return nil, &attachments;
}

type Message struct{
	Message_id string
	Channel_id string
	Guild_id string
	Message_timestamp int
	Content string
	Sender_id string
	Sender_name string
	Reply_to string
	Edit_time string
	Edits []Edit
	Embeds []Embed
	Attachments []Attachment
}

type Embed struct{
	Message_id string
	Embed_url string
	Embed_title string
	Embed_description string
	Embed_timestamp string
	Embed_thumbnail_url string
	Embed_thumbnail_hash string
	Embed_image_url string
	Embed_image_hash string
	Embed_footer string
	Embed_author_name string
	Embed_author_url string
	Embed_field string
}

type Attachment struct{
	Attachment_id string
	Message_id string
	Attachment_filename string
	Attachment_url string
	Attachment_hash string
}

type Edit struct{
	Message_id string
	Edit_time string
	Content string
}

type Messages struct {
	Messages []Message
}
