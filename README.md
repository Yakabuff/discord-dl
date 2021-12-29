# discord-dl

discord-dl is a tool to archive discord channels.  

I think it's safe to say that forums, the medium that has held countless communities and the internet together for decades and has acted as its primary knowledge center, are dead.  Communities are moving en-masse to walled off platforms like discord that cannot be indexed by search engines and are completely inaccesible by those who are not users of these platforms.  As a result, it is imperative to have a way of preserving these communities such that its content can be openly accesible and preserved. 

## Installation

```make build```

```cd bin```

```go install discord-dl```

## Features

- Archives messages, edits, embeds and attachments
- SQLite database storing message history

## To-do

- Web api and frontend to query messages
- BOT capabilities so it can be installed by server administrators to archive messages in real time
- Ability to 'rebuild' the server in the event of deletion
- More supported databases
- Cloud storage support

## Instructions

- ## Flags: 
    - `--t=token` 
        - Specify the token of your user agent.  Prepend `Bot ` to the token for bot tokens. 
        - eg: `--t='Bot <token>'`
    - `--channel=channel_id`
        - Specify the ID of the channel to archive
    - `--guild=guild_id`
        - Specify the ID of the guild to archive
    - `--before=YYYY-MM-DD` or `--before=message_id`
        - Specify the date or message ID to get messages before
    - `--after=YYYY-MM-DD` or `--after=message_id`
        - Specify the date or message ID to get messages after
    - `--dms=user_id`
        - Specify user_id of the DM to archive. Not implemented
    - `--progress`
        - Output your progress
    - `--fast-update`
        - Fetches messages until it reaches an already archived message.
    - `--listen`
        - Listens for new messages (BOT only). Not implemented
    - `--download_media`
        - Enable this flag to download attachments/files (Includes embedded files)

## Contributing
Pull requests are welcome. Please open an issue first to discuss what you would like to change.

## License
GNU AGPL v3