package db

import (
	"database/sql"
	"log"
)

const CHANNEL_NAMES = "channel_names"
const CHANNEL_TOPICS = "channel_topics"
const GUILD_NAMES = "guild_names"
const GUILD_ICONS = "guild_icons"
const GUILD_BANNER = "guild_banners"

func createTable(db *sql.DB) {
	createChannels := `CREATE TABLE channels(
		"channel_id" TEXT NOT NULL PRIMARY KEY
	);`

	createChannelsMetadata := `CREATE TABLE channels_meta(
		"channel_id" TEXT NOT NULL PRIMARY KEY,
		"name" TEXT NOT NULL,
		"topic" TEXT,
		"guild_id" TEXT NOT NULL,
		"is_thread" INTEGER NOT NULL CHECK(is_thread >= 0 and is_thread <= 1),
		FOREIGN KEY (guild_id) REFERENCES guilds(guild_id),
		FOREIGN KEY (channel_id) REFERENCES channels(channel_id)
	);`

	//insert new name for channel if most recent is not the same as new name
	createChannelHistoricalNames := `CREATE TABLE channel_names(
		"channel_id" TEXT NOT NULL,
		"date_renamed" INTEGER NOT NULL,
		"channel_name" TEXT NOT NULL,
		FOREIGN KEY (channel_id) REFERENCES channels(channel_id),
		PRIMARY KEY(channel_id, date_renamed)
	);`

	createChannelHistoricalTopic := `CREATE TABLE channel_topics(
		"channel_id" TEXT NOT NULL,
		"date_renamed" INTEGER NOT NULL,
		"channel_topic" TEXT NOT NULL,
		FOREIGN KEY (channel_id) REFERENCES channels(channel_id),
		PRIMARY KEY(channel_id, date_renamed)
	);`

	createGuilds := `CREATE TABLE guilds(
		"guild_id" TEXT NOT NULL PRIMARY KEY
	);`
	createGuildsMetadata := `CREATE TABLE guilds_meta(
		"guild_id" TEXT NOT NULL PRIMARY KEY,
		"icon" TEXT,
		"banner" TEXT,
		"name" TEXT NOT NULL
	);`

	createGuildHistoricalNames := `CREATE TABLE guild_names(
		"guild_id" TEXT NOT NULL,
		"date_renamed" INTEGER NOT NULL,
		"guild_name" TEXT NOT NULL,
		FOREIGN KEY (guild_id) REFERENCES guilds(guild_id),
		PRIMARY KEY(guild_id, date_renamed)
	);`

	createGuildHistoricalIcon := `CREATE TABLE guild_icons(
		"guild_id" TEXT NOT NULL,
		"date_renamed" INTEGER NOT NULL,
		"guild_icon_hash" TEXT NOT NULL,
		FOREIGN KEY (guild_id) REFERENCES guilds(guild_id),
		PRIMARY KEY(guild_id, date_renamed)
	);`

	createGuildHistoricalBanner := `CREATE TABLE guild_banners(
		"guild_id" TEXT NOT NULL,
		"date_renamed" INTEGER NOT NULL,
		"guild_banner_hash" TEXT NOT NULL NOT NULL,
		FOREIGN KEY (guild_id) REFERENCES guilds(guild_id),
		PRIMARY KEY(guild_id, date_renamed)
	);`

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
		"thread_id" TEXT,
		FOREIGN KEY (channel_id) REFERENCES channels(channel_id)
	);`

	createAttachments := `CREATE TABLE attachments (
		"message_id" TEXT NOT NULL,  
		"attachment_id" TEXT NOT NULL PRIMARY KEY,
		"attachment_filename" TEXT NOT NULL,
		"attachment_URL" TEXT NOT NULL,
		"attachment_hash" TEXT NOT NULL,
		FOREIGN KEY (message_id) REFERENCES messages(message_id)
	);`

	createEdits := `CREATE TABLE edits (
		"message_id" TEXT,  
		"edit_time" INTEGER,
		"content" TEXT,
		PRIMARY KEY(message_id, edit_time, content),
		FOREIGN KEY (message_id) REFERENCES messages(message_id)
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
		),
		FOREIGN KEY (message_id) REFERENCES messages(message_id)
	);`

	log.Println("Create channels table...")
	statement_channel, err := db.Prepare(createChannels) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement_channel.Exec() // Execute SQL Statements
	log.Println("Channels table created")

	log.Println("Create channels metatable...")
	statement_channel_meta, err := db.Prepare(createChannelsMetadata) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement_channel_meta.Exec() // Execute SQL Statements
	log.Println("Channels meta table created")

	log.Println("Create channel names table...")
	channel_names, err := db.Prepare(createChannelHistoricalNames) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	channel_names.Exec() // Execute SQL Statements
	log.Println("channel names table created")

	log.Println("Create channel names table...")
	channel_topics, err := db.Prepare(createChannelHistoricalTopic) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	channel_topics.Exec() // Execute SQL Statements
	log.Println("channel names table created")

	log.Println("Create guild historical icon table...")
	guild_icon, err := db.Prepare(createGuildHistoricalIcon) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	guild_icon.Exec() // Execute SQL Statements
	log.Println("guild historical icon table created")

	log.Println("Create guild historical banner table...")
	guild_banner, err := db.Prepare(createGuildHistoricalBanner) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	guild_banner.Exec() // Execute SQL Statements
	log.Println("guild historical banner created")

	log.Println("Create guild table...")
	statement_guilds, err := db.Prepare(createGuilds) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement_guilds.Exec() // Execute SQL Statements
	log.Println("guilds table created")

	log.Println("Create guild meta table...")
	statement_guilds_meta, err := db.Prepare(createGuildsMetadata) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement_guilds_meta.Exec() // Execute SQL Statements
	log.Println("guilds meta table created")

	log.Println("Create guild names table...")
	guild_names, err := db.Prepare(createGuildHistoricalNames) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	guild_names.Exec() // Execute SQL Statements
	log.Println("guilds names table created")

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
