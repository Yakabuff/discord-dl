package archiver

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/yakabuff/discord-dl/db"
	"github.com/yakabuff/discord-dl/models"
	"github.com/yakabuff/discord-dl/web"
)

type Archiver struct {
	Db   db.Db
	Args models.Args
	Dg   *discordgo.Session
	Web  web.Web
}

func (a Archiver) ParseCliArgs() error {
	if a.Args.Mode == models.INPUT {
		fmt.Println("Selected input mode")
		err := a.parseConfig(a.Args.Input)
		if err != nil {
			return err
		}
	}

	//Either listening or web deploy
	if a.Args.Listen == true {
		log.Println("Listening for changes...")
		a.addHandlers()
	}
	if a.Args.Deploy == true {
		log.Println("Starting webview...")
		a.Web = web.NewWeb(a.Db, a.Args.DeployPort, a.Args.MediaLocation)
		a.Web.Deploy(a.Db)
	}

	if a.Args.Mode != models.NONE {
		switch a.Args.Mode {
		case models.GUILD:
			fmt.Println("Archiving guild")
			err := a.GuildDownload(a.Args.Guild)
			if err != nil {
				return err
			}
		case models.CHANNEL:
			fmt.Println("Selected channel mode")
			err := a.ChannelDownload(a.Args.Channel)
			if err != nil {
				return err
			}
		}
	} else {
		//add slash command listener
		fmt.Println("Bot is now running.  Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	}
	return nil
}

func (a Archiver) parseConfig(fileName string) error {
	//read file
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	//parse each element of json as an args type
	decoder := json.NewDecoder(file)
	config := models.Args{}
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}
	a.Args = config
	mode := checkFlagMode("", config.Guild, config.Channel)
	if mode == models.INVALID {
		return errors.New("Invalid mode")
	}
	a.Args.Mode = mode
	return nil
}

func checkFlagMode(input string, guild string, channel string) models.Mode {
	var count int
	var mode models.Mode
	if input != "" {
		count++
		mode = models.INPUT
	}
	if guild != "" {
		count++
		mode = models.GUILD
	}
	if channel != "" {
		count++
		mode = models.CHANNEL
	}
	if count == 1 {
		return mode
	} else if count > 1 {
		return models.INVALID
	} else {
		return models.NONE
	}
}

func (a Archiver) InitCli() models.Args {
	progress := flag.Bool("progress", false, "Displays progress of task. Enabling this will output verbose logging")
	before := flag.String("before", "", "Retrieves all messages before this date or message id")
	after := flag.String("after", "", "Retrieves all messages after this date or message id")
	fastUpdate := flag.Bool("fast-update", false, "Retrieves all message after the last downloaded message")
	downloadMedia := flag.Bool("download-media", false, "downloads embedded images and files from message")
	token := flag.String("t", "", "Sets user or bot token")
	output := flag.String("o", "", "Sets output db path")
	input := flag.String("i", "", "Input mode. Gets config from input file")
	guild := flag.String("guild", "", "Guild mode. Retrieves messages and channels from selected guild")
	channel := flag.String("channel", "", "Retrieves messages from selected channel")
	listen := flag.Bool("listen", false, "Listens for new messages/events and archives in real time.  Can only be used with a bot account")
	deploy := flag.Bool("deploy", false, "Deploys webapp")
	deployPort := flag.Int("deploy_port", 8080, "Set webview port")
	blacklistedChannels := flag.String("blacklisted-channels", "", "Sets list of blacklisted channel IDs as a string delimited by a space. Can only be used with guild")
	mediaLocation := flag.String("media-location", "media", "Set location to store attachments and media")
	flag.Parse()

	mode := checkFlagMode(*input, *guild, *channel)

	if mode == models.INVALID {
		fmt.Fprintln(os.Stderr, "Invalid flags")
		os.Exit(1)
	}

	args := models.Args{
		Progress:            *progress,
		DownloadMedia:       *downloadMedia,
		MediaLocation:       *mediaLocation,
		Before:              *before,
		After:               *after,
		FastUpdate:          *fastUpdate,
		Token:               *token,
		Output:              *output,
		Input:               *input,
		Guild:               *guild,
		Channel:             *channel,
		Listen:              *listen,
		Deploy:              *deploy,
		DeployPort:          *deployPort,
		BlacklistedChannels: strings.Split(*blacklistedChannels, " "),
		Mode:                mode}

	if *input != "" && len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "Option --i cannot be used in conjunction with other flags")
		os.Exit(1)
	}

	if *guild != "" && *channel != "" {
		fmt.Fprintln(os.Stderr, "Cannot use --guild and --channel together")
		os.Exit(1)
	}

	if (*before != "" || *after != "") && *fastUpdate != false {
		fmt.Fprintln(os.Stderr, "Cannot have before/after flags with fast-update")
		os.Exit(1)
	}

	if *before != "" && *after != "" {
		if strings.Contains(*before, "-") && !strings.Contains(*after, "-") || !strings.Contains(*before, "-") && strings.Contains(*after, "-") {
			fmt.Fprintln(os.Stderr, "Before and after flags must be in the same format")
			os.Exit(1)
		}

		if !strings.Contains(*before, "-") && !strings.Contains(*after, "-") {
			fmt.Fprintln(os.Stderr, "Invalid date format")
			os.Exit(1)
		}
	}

	return args
}
