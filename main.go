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

	var archiver = archiver.Archiver{}
	jobArgs, args := archiver.InitCli()
	archiver.Args = args
	db, err := db.Init_db(archiver.Args.Output)
	if err != nil {
		panic(err.Error())
	}

	errDg, dg := archiver.CreateConnection()
	if errDg != nil {
		panic(errDg.Error())
	}
	archiver.Dg = dg

	archiver.Db = *db
	archiver.Queue = job.NewJobQueue(&archiver)

	err = archiver.ParseArgs()
	if err != nil {
		panic(err.Error())
	}
	if jobArgs.Mode != models.NONE {
		//Wait until job is complete and then exit

		// log.Println(archiver.Queue.Queue["asdf"].Category)
		// go func() {
		// 	for i := 0; i < 10; i++ {
		// 		archiver.Queue.Enqueue(job.NewJob(jobArgs))
		// 	}
		// }()

		// go func() {
		// 	for i := 0; i < 10; i++ {
		// 		ja := models.JobArgs{Mode: models.CHANNEL, Channel: "asdf123"}
		// 		archiver.Queue.Enqueue(job.NewJob(ja))
		// 	}
		// }()
		// time.Sleep(5 * time.Second)
		// log.Println("waiting for job to finish")
		archiver.Queue.Enqueue(job.NewJob(jobArgs))
		archiver.Queue.Wg.Wait()
		log.Println("jobs have finished")
		// log.Println(archiver.Queue.Jobs)

	} else {
		//If no job, run forever and wait for jobs
		fmt.Println("Bot is now running.  Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	}

}
