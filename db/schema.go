package db

var schema string = `
PRAGMA user_version = 1;
CREATE TABLE channels("channel_id" TEXT NOT NULL PRIMARY KEY);

CREATE TABLE guilds(
	"guild_id" TEXT NOT NULL PRIMARY KEY
);

CREATE TABLE channel_topics(
	"channel_id" TEXT NOT NULL,
	"date_renamed" INTEGER NOT NULL,
	"channel_topic" TEXT NOT NULL,
	FOREIGN KEY (channel_id) REFERENCES channels(channel_id),
	PRIMARY KEY(channel_id, date_renamed)
);

CREATE TABLE guilds_meta(
	"guild_id" TEXT NOT NULL PRIMARY KEY,
	"icon" TEXT,
	"banner" TEXT,
	"name" TEXT NOT NULL
);

CREATE TABLE guild_names(
	"guild_id" TEXT NOT NULL,
	"date_renamed" INTEGER NOT NULL,
	"guild_name" TEXT NOT NULL,
	FOREIGN KEY (guild_id) REFERENCES guilds(guild_id)
	PRIMARY KEY(guild_id, date_renamed)
);

CREATE TABLE guild_icons(
	"guild_id" TEXT NOT NULL,
	"date_renamed" INTEGER NOT NULL,
	"guild_icon_hash" TEXT NOT NULL,
	FOREIGN KEY (guild_id) REFERENCES guilds(guild_id),
	PRIMARY KEY(guild_id, date_renamed)
);

CREATE TABLE guild_banners(
	"guild_id" TEXT NOT NULL,
	"date_renamed" INTEGER NOT NULL,
	"guild_banner_hash" TEXT NOT NULL NOT NULL,
	FOREIGN KEY (guild_id) REFERENCES guilds(guild_id),
	PRIMARY KEY(guild_id, date_renamed)
);

CREATE TABLE messages (
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
);

CREATE TABLE attachments (
	"message_id" TEXT NOT NULL,
	"attachment_id" TEXT NOT NULL PRIMARY KEY,
	"attachment_filename" TEXT NOT NULL,
	"attachment_URL" TEXT NOT NULL,
	"attachment_hash" TEXT NOT NULL,
	FOREIGN KEY (message_id) REFERENCES messages(message_id)
);

CREATE TABLE edits (
	"message_id" TEXT,
	"edit_time" INTEGER,
	"content" TEXT,
	PRIMARY KEY(message_id, edit_time, content),
	FOREIGN KEY (message_id) REFERENCES messages(message_id)
);

CREATE TABLE embeds (
	"message_id" TEXT NOT NULL,
	"embed_date_retrieved" TEXT NOT NULL,
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
		embed_date_retrieved,
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
);

CREATE TABLE channels_meta(
	"channel_id" TEXT NOT NULL PRIMARY KEY,
	"name" TEXT NOT NULL,
	"topic" TEXT,
	"guild_id" TEXT NOT NULL,
	"is_thread" INTEGER NOT NULL CHECK(is_thread >= 0 and is_thread <= 1),
	FOREIGN KEY (guild_id) REFERENCES guilds(guild_id),
	FOREIGN KEY (channel_id) REFERENCES channels(channel_id)
);

CREATE TABLE channel_names(
	"channel_id" TEXT NOT NULL,
	"date_renamed" INTEGER NOT NULL,
	"channel_name" TEXT NOT NULL,
	FOREIGN KEY (channel_id) REFERENCES channels(channel_id),
	PRIMARY KEY(channel_id, date_renamed)
);

CREATE TRIGGER guildNameTrigger BEFORE INSERT ON guild_names
BEGIN
SELECT
	CASE
		WHEN EXISTS (SELECT * from (select * from guild_names where guild_id = NEW.guild_id order by date_renamed DESC LIMIT 1) t WHERE t.guild_name IS NEW.guild_name order by date_renamed DESC)
			THEN RAISE (ABORT, 'guildNameTrigger violated')
	END;
END;

CREATE TRIGGER guildIconTrigger BEFORE INSERT ON guild_icons
BEGIN
SELECT
	CASE
		WHEN EXISTS (SELECT * from (select * from guild_icons where guild_id = NEW.guild_id order by date_renamed DESC LIMIT 1) t WHERE t.guild_icon_hash IS NEW.guild_icon_hash order by date_renamed DESC)
			THEN RAISE (ABORT, 'guildIconTrigger violated')
	END;
END;

CREATE TRIGGER guildBannerTrigger BEFORE INSERT ON guild_banners
BEGIN
SELECT
	CASE
		WHEN EXISTS (SELECT * from (select * from guild_banners where guild_id = NEW.guild_id order by date_renamed DESC LIMIT 1) t WHERE t.guild_banner_hash IS NEW.guild_banner_hash order by date_renamed DESC)
			THEN RAISE (ABORT, 'guildBannerTrigger violated')
	END;
END;

CREATE TRIGGER channelTopicTrigger BEFORE INSERT ON channel_topics
BEGIN
SELECT
	CASE
		WHEN EXISTS (SELECT * from (select * from channel_topics where channel_id = NEW.channel_id order by date_renamed DESC LIMIT 1) t WHERE t.channel_topic IS NEW.channel_topic order by date_renamed DESC)
			THEN RAISE (ABORT, 'channelTopicTrigger violated')
	END;
END;

CREATE TRIGGER channelNameTrigger BEFORE INSERT ON channel_names
BEGIN
SELECT
	CASE
		WHEN EXISTS (SELECT * from (select * from channel_names where channel_id = NEW.channel_id order by date_renamed DESC LIMIT 1) t WHERE t.channel_name IS NEW.channel_name order by date_renamed DESC)
			THEN RAISE (ABORT, 'channelNameTrigger violated')
	END;
END;

`
