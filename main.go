package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yakabuff/discord-dl/archiver"
	"github.com/yakabuff/discord-dl/db"
	"github.com/yakabuff/discord-dl/job"
	"github.com/yakabuff/discord-dl/models"
)

func main() {

	jobArgs, args := archiver.InitCli()

	//Parse config file if specified
	if args.Input != "" {
		err := archiver.ParseConfigFile(args.Input, &args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	if !archiver.ValidFlags(jobArgs, args) {
		fmt.Fprintln(os.Stderr, "Invalid flags")
		os.Exit(1)
	}

	var theArchiver = archiver.Archiver{Args: args}

	theArchiver.InitLogger()
	if args.Output != "" {
		db, err := db.Init_db(theArchiver.Args.Output)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		theArchiver.Db = *db
	}

	if args.Token != "" {
		errDg, dg := theArchiver.CreateConnection()
		if errDg != nil {
			// log.Println(theArchiver.Args.Token)
			fmt.Fprintln(os.Stderr, errDg.Error())
			os.Exit(1)
		}

		theArchiver.Dg = dg

		theArchiver.InitListener()
	}

	err := theArchiver.InitWeb()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	theArchiver.Queue = job.NewJobQueue(&theArchiver, theArchiver.Args.Logging)

	if jobArgs.Mode != models.NONE {
		//Wait until job is complete and then exit
		theArchiver.Queue.Enqueue(job.NewJob(jobArgs))
		theArchiver.Queue.Wg.Wait()
		theArchiver.Queue.Progress.Wait()
		log.Println("Job has finished")

	} else {
		//If no job, run forever and wait for jobs
		fmt.Println("discord-dl is now running.  Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	}

}
