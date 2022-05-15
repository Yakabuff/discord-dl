package archiver

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/yakabuff/discord-dl/db"
	"github.com/yakabuff/discord-dl/job"
	"github.com/yakabuff/discord-dl/models"
	"github.com/yakabuff/discord-dl/web"
)

type Archiver struct {
	Db    db.Db
	Args  models.ArchiverArgs
	Dg    *discordgo.Session
	Web   web.Web
	Queue job.JobQueue
	Wg    sync.WaitGroup
}

func (a Archiver) ParseArgs() error {
	// if a.Args.Mode == models.INPUT {
	// 	fmt.Println("Selected input mode")
	// 	err := a.parseConfig(a.Args.Input)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	//Either listening or web deploy
	if a.Args.Listen == true && strings.HasPrefix(a.Args.Token, "Bot") {
		log.Println("Listening for changes...")
		a.addHandlers()
	}
	if a.Args.Deploy == true {
		log.Println("Starting webview...")

		a.Web = web.NewWeb(a.Db, a.Args.DeployPort, a.Args.MediaLocation, &a.Queue)
		a.Web.Deploy(a.Db)
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
	config := models.ArchiverArgs{}
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}
	a.Args = config
	return nil
}

func checkFlagMode(input string, guild string, channel string) models.Mode {
	var count int
	var mode models.Mode
	// if input != "" {
	// 	count++
	// 	mode = models.INPUT
	// }
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

func (a Archiver) InitCli() (models.JobArgs, models.ArchiverArgs) {
	// progress := flag.Bool("progress", false, "Displays progress of task. Enabling this will output verbose logging")
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

	job := models.JobArgs{
		Mode:       mode,
		Before:     *before,
		After:      *after,
		FastUpdate: *fastUpdate,
		Guild:      *guild,
		Channel:    *channel,
	}

	args := models.ArchiverArgs{
		Mode:                mode,
		DownloadMedia:       *downloadMedia,
		MediaLocation:       *mediaLocation,
		Token:               *token,
		Output:              *output,
		Input:               *input,
		Listen:              *listen,
		Deploy:              *deploy,
		DeployPort:          *deployPort,
		BlacklistedChannels: strings.Split(*blacklistedChannels, " "),
	}

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

	return job, args
}
