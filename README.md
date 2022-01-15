# discord-dl

discord-dl is a CLI tool and Discord Bot that can be used to archive discord channels in realtime.  

It's safe to say that forums, the medium that has acted as the internet's primary knowledge center, are dead.  Communities are moving en-masse to walled off platforms like Discord that cannot be indexed by search engines and are completely inaccesible by those who are not registered users.  As a result, it is imperative to have a way of preserving these communities so its content can be openly accessible. 

## Installation
```git clone https://github.com/Yakabuff/discord-dl.git```

```cd discord-dl```

```make build```

```cd bin```

```go install discord-dl``` or ```./discord-dl {flags}```

## Features

- Archives messages, edits, embeds and attachments
- Storing message history in an SQLite database
- BOT capabilities so it can be installed by server administrators to archive messages in real time
- Simple web API/frontend to display and query messages

## To-do

- Ability to 'rebuild' the server in the event of deletion
- Commands to start tasks without having to restart the Bot
- More supported databases
- Cloud storage support

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
    - `--dms=true`
        - Downloads all DM channels
    - `--progress`
        - Output your progress. (Not implemented)
    - `--fast-update`
        - Fetches messages until it reaches an already archived message
    - `--listen`
        - Listens for new messages (BOT only). Can be used in conjunction with other modes
    - `--download_media`
        - Enable this flag to download attachments/files (Includes embedded files) (Not implemented)
    - `--deploy`
        - Mode to start the web server.  Can be used by itself or in conjunction with other modes.

## Contributing
Pull requests are welcome. Please open an issue first to discuss what you would like to change.

## License
GNU AGPL v3