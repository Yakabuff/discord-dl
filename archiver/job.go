package archiver

import (
	"github.com/yakabuff/discord-dl/models"
)

type Job struct {
	Id       int
	Progress int
	Args     models.Args
}

func Enqueue(job Job, jobChan <-chan Job) {

}

//note on jobs:
//archiver has a set of jobs in a datastructure or a channel.
//archiver will assign a worker to each job.(non blocking). Max of 10 jobs?
//Each worker will be running on a new goroutine
//each job will get same DB and DG but different set of args

//each worker can be cancelled. can get progress of each job
