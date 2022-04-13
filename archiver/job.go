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

func (j Job) Exec() {
	//start new go routine -> run parseCLI args
}

//note on jobs:
//archiver has a set of jobs in a datastructure or a channel.
//archiver will assign a worker to each job.(non blocking). Max of 10 jobs?
//Each worker will be running on a new goroutine
//each job will get same DB and DG but different set of args

//each worker can be cancelled. can get progress of each job

//job queue: I'm thinking instead of a queue of channels that fire concurrently.. it should be like this:

//queue channel1 -> in progress, queue channel1 -> in queue, channel 1 is in process of archiving, queue channel2 -> in progress , queue channel 3 -> in progress
//queue channel 2 -> in queue, channel 2(1) completed, channel 2(2) -> in progress
//{[channel1, channel1], [channel2, channel2], [guild1, guild1]}
