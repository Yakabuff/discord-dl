package archiver

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/yakabuff/discord-dl/common"
	"github.com/yakabuff/discord-dl/db"
	"github.com/yakabuff/discord-dl/job"
	"github.com/yakabuff/discord-dl/models"
	"github.com/yakabuff/discord-dl/web"
)

const VERSION = "3.0.0-alpha"

type Archiver struct {
	Db    db.Db
	Args  models.ArchiverArgs
	Dg    *discordgo.Session
	Web   web.Web
	Queue job.JobQueue
	Wg    sync.WaitGroup
}

var log *logrus.Logger

func (a Archiver) InitLogger() {
	if a.Args.Logging {
		l, err := common.NewErrLogger()
		if err != nil {
			logrus.New().Fatal(err)
		}
		log = l
		log.SetReportCaller(true)
	} else {
		logrus.SetOutput(ioutil.Discard)
	}
}
func (a Archiver) ParseArgs() error {
	//Either listening or web deploy
	if a.Args.Listen == true && strings.HasPrefix(a.Args.Token, "Bot") {
		log.Info("Listening for changes...")
		a.addHandlers()
	}
	if a.Args.Deploy == true {
		log.Info("Starting webview...")

		a.Web = web.NewWeb(a.Db, a.Args.DeployPort, a.Args.MediaLocation, &a.Queue, a.Args.Logging)
		a.Web.Deploy(a.Db)
	}

	return nil
}

func ParseConfigFile(fileName string, args *models.ArchiverArgs) error {

	_, err := toml.DecodeFile(fileName, &args)

	if err != nil {
		return err
	}
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

func InitCli() (models.JobArgs, models.ArchiverArgs) {

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
	listenChannels := flag.String("listen-channels", "", "Sets list of channels you wish to listen to")
	listenGuilds := flag.String("listen-guilds", "", "Sets list of guilds you wish to listen to")
	deploy := flag.Bool("deploy", false, "Deploys webapp")
	deployPort := flag.Int("deploy-port", 8080, "Set webview port")
	blacklistedChannels := flag.String("blacklisted-channels", "", "Sets list of blacklisted channel IDs as a string delimited by a comma. Can only be used with guild")
	mediaLocation := flag.String("media-location", "media", "Set location to store attachments and media")
	version := flag.Bool("version", false, "Checks version")
	logging := flag.Bool("log", true, "Verbose logging to file")
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

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
		DownloadMedia:       *downloadMedia,
		MediaLocation:       *mediaLocation,
		Token:               *token,
		Output:              *output,
		Input:               *input,
		Listen:              *listen,
		Deploy:              *deploy,
		DeployPort:          *deployPort,
		BlacklistedChannels: strings.Split(*blacklistedChannels, ","),
		ListenChannels:      strings.Split(*listenChannels, ","),
		ListenGuilds:        strings.Split(*listenGuilds, ","),
		Logging:             *logging,
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
