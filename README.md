
![logo](ddl.png)

**Discord-dl** is an utility that can be used to archive Discord channels and guilds.  

***

It's safe to say that forums, the medium that has once acted as the internet's primary knowledge center, are dead.  Communities are moving en-masse to information blackholes like Discord that cannot be indexed by search engines and are completely inaccesible by those who are not registered users.  It is vital to have a way of preserving these communities so its content can remain openly accessible. 

***

**Note:** While I recognize that there is no alternative to selfbots especially when it comes to archiving, it is unfortunately against Discord TOS.  Use at your own discretion.

**Master branch is experimental and unstable. Please use a stable release.  The stable release is also highly experimental, unfortunately :^)**

## Installation

1) Install Go

2) Make a Bot account: https://discord.com/developers/applications

3) ```git clone https://github.com/Yakabuff/discord-dl.git```

    ```cd discord-dl```

    ```make build```

    ```cd bin```

    ```go install discord-dl```

4) Invite the bot to the server: https://discordpy.readthedocs.io/en/stable/discord.html

5) Daemonize the bot (optional) or run the bot

## Features

- Archives messages, edits, embeds, threads and attachments
- Storing message history in an SQLite database
- Listens for messages to archive in real time
- Simple web API/frontend to display and query messages
- Job queue for improved concurrency and the ability to pause/stop jobs/see progress

## To-do

- Ability to 'rebuild' the server in the event of deletion
- Cloud object storage support
- Postgres support

## Instructions

- ## Flags: 
    - `--t=token` 
        - Specify the token of your user agent.  Prepend `Bot ` to the token for bot tokens. 
        - eg: `--t='Bot <token>'` or `--t=<token>`
    - `--channel=channel_id`
        - Mode that archives the specified channel ID
    - `--guild=guild_id`
        - Mode that archives the specified guild ID
    - `--before=YYYY-MM-DD` or `--before=message_id`
        - Specify the date or message ID to get messages before
    - `--after=YYYY-MM-DD` or `--after=message_id`
        - Specify the date or message ID to get messages after
    - `--progress`
        - Output your progress. (Not implemented)
    - `--fast-update`
        - Fetches messages until it reaches an already archived message
    - `--listen`
        - Listens for new messages (BOT only) in selected channels and guilds. Can be used in conjunction with other modes
        - Usage: `--listen`
    - `--listen-guilds`
        - Only listen for messages from specific guilds
        - Usage: `--listen-guilds=guild1,guild2,guild3`
    - `--listen-channels`
        - Only listen for messages from specific channels
        - Usage: `--listen-channels=channelid1,channelid2`
    - `--download-media`
        - Enable this flag to download attachments/files (Includes embedded files)
    - `--deploy`
        - Mode to start the web server.  Can be used by itself or in conjunction with other modes.
    - `--input="config_path"`
        - Specify path to config file.  Values from the input file (TOML) will precedence over values from CLI flags.
    - `--output="path"`
        - Specify path to database.  Will fallback to default value (archive.db) if empty.
    - `--media-location`
        - Specify location where you want images, videos and attachments to be saved
    - `--deployPort`
        - Specify port you want the webserver to run on. Will default to 8080 if empty
    -  `--blacklisted-channels`
        - Specify channels you which to exclude.  Delimit the channels by a comma. 
        - ex: `--blacklisted-channels=channel1id,channel2id,channel3id`

- ## Examples:
    - Downloading messages and media from a channel after Jan 1st 2018:
        - `discord-dl --channel="12371093817" --after="2018-01-01" --download-media=true --token="asdfasdfasdfasdf" --o="test.db" --media-location="~/pics"`
## Contributing
Pull requests are welcome. Please open an issue first to discuss what you would like to change.

## License
GNU AGPL v3