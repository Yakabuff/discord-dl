package job

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/yakabuff/discord-dl/models"
)

type Archiver interface {
	ChannelDownload(channel string, fastUpdate bool, after string, before string, state JobState) error
	GetChannelsGuild(guildID string) ([]string, error)
	IndexGuild(guild string) error
	IndexChannel(channel string) error
}
type Jobs struct {
	Jobs []JobOut
}

type JobOut struct {
	Id        string
	Snowflake string
	Progress  int
	Status    string
	Error     string
}

type Job struct {
	Id        string
	Snowflake string
	Progress  int
	Args      models.JobArgs
	Status    Status
	Error     error
}
type JobQueue struct {
	*sync.Mutex
	Queue    map[string]*Worker
	MaxSize  int
	Wg       sync.WaitGroup
	Archiver Archiver
	Jobs     map[string]*Job
}

// type JobRecord struct {
// 	*sync.Mutex
// 	Jobs map[string]*Job
// }

type Worker struct {
	Channel  *chan *Job
	Category string
	MaxSize  int
	CurrJob  *Job
}

type JobState struct {
	Progress *int
	Error    *error
	Status   *Status
}

type Status int

const (
	PENDING Status = iota
	RUNNING
	CANCELLED
	ERROR
	FINISHED
)

func NewJob(args models.JobArgs) Job {
	id := uuid.New().String()
	var snowflake string
	switch args.Mode {
	case models.CHANNEL:
		snowflake = args.Channel
	case models.GUILD:
		snowflake = args.Channel
	}

	job := Job{Id: id, Args: args, Snowflake: snowflake, Status: PENDING}
	return job
}

func NewJobQueue(a Archiver) JobQueue {
	m := sync.Mutex{}
	q := make(map[string]*Worker)
	jr := make(map[string]*Job)
	return JobQueue{
		&m,
		q,
		60,
		sync.WaitGroup{},
		a,
		jr,
	}
}

// func NewJobRecord() JobRecord {
// 	m := sync.Mutex{}
// 	q := make(map[string]*Job)
// 	return JobRecord{
// 		&m,
// 		q,
// 	}
// }

func (a *JobQueue) AddJobRecord(job *Job) {

	a.Jobs[job.Id] = job

}

func (a *JobQueue) Enqueue(job Job) error {

	a.Lock()
	defer a.Unlock()

	if _, ok := a.Queue[job.Snowflake]; ok {
		log.Println("Worker already exists")
		//If group exists, add job to preexisting worker
		var w *chan *Job = a.Queue[job.Snowflake].Channel
		select {
		case *w <- &job:
			log.Println("Worker already exists, job added succesfully")
			a.AddJobRecord(&job)
		default:
			log.Println("Channel is full. Job was not added")
			return errors.New("channel is full. Job was not added")
		}

	} else {
		c := make(chan *Job, 100)
		a.Queue[job.Snowflake] = &Worker{Channel: &c, Category: job.Snowflake}
		// log.Println("State of map before adding to a supposedly non existent key")
		// log.Println(a.Queue.Queue)

		var w *chan *Job = a.Queue[job.Snowflake].Channel
		select {
		case *w <- &job:
			//Startup a worker and start processing jobs
			a.Wg.Add(1)
			log.Printf("Job ID %s added for newly created worker. Task commencing", job.Id)
			a.AddJobRecord(&job)
			go a.StartWorker(a.Queue[job.Snowflake])
			// time.Sleep(10 * time.Second)
		default:
			log.Println("Channel is full. Job was not added")
			return errors.New("channel is full. Job was not added")
		}

	}
	return nil
}

func (a *JobQueue) StartWorker(w *Worker) {
	log.Println("executing worker")
	for {
		select {
		case task := <-*w.Channel:
			w.CurrJob = task
			// if task != nil {
			// 	//Execute job struct
			// 	log.Println("executing job")
			// }
			task.Status = RUNNING
			w.Category = task.Snowflake
			log.Printf("executing job %s %s", w.Category, task.Id)
			a.ExecJobArgs(task.Args, task)
			log.Printf("Finished executing job %s %s", w.Category, task.Id)
		default:
			log.Println("Last job processed in " + w.Category)
			delete(a.Queue, w.Category)
			a.Wg.Done()
			return
		}
	}
}

func (a *JobQueue) ExecJobArgs(j models.JobArgs, job *Job) {
	// var err error
	if j.Mode != models.NONE {
		switch j.Mode {
		case models.GUILD:
			fmt.Println("Archiving guild")
			// call function in archiver that returns list of channels in guild
			// index guild metadata
			// queue channels from list
			job.Status = RUNNING
			err := a.Archiver.IndexGuild(job.Snowflake)
			if err != nil {
				job.Error = err
				job.Status = ERROR
				return
			} else {
				job.Status = FINISHED
			}

			guilds, err := a.Archiver.GetChannelsGuild(j.Guild)
			if err != nil {
				job.Error = err
				job.Status = ERROR
				return
			} else {
				job.Status = FINISHED
			}

			for _, val := range guilds {
				ja := models.JobArgs{Mode: models.GUILD, Before: job.Args.Before, After: job.Args.After, FastUpdate: job.Args.FastUpdate, Guild: "", Channel: val}
				jobtmp := NewJob(ja)
				a.Enqueue(jobtmp)
			}
			if err != nil {
				job.Status = ERROR
				job.Error = err
			} else {
				job.Status = FINISHED
			}

			log.Println("finished archiving guild......")
		case models.CHANNEL:
			fmt.Println("Selected channel mode")
			job.Status = RUNNING
			var state JobState = JobState{Progress: &job.Progress, Error: &job.Error, Status: &job.Status}
			err := a.Archiver.IndexChannel(j.Channel)
			if err != nil {
				job.Error = err
				job.Status = ERROR
				return
			} else {
				job.Status = FINISHED
			}

			err = a.Archiver.ChannelDownload(j.Channel, j.FastUpdate, j.After, j.Before, state)

			job.Progress = 100
			if err != nil {
				job.Error = err
				job.Status = ERROR
				if errors.Is(err, models.FastUpdateError) {
					job.Status = FINISHED
				}
			} else {
				job.Status = FINISHED
			}
		}
	}
}

func (a *JobQueue) CancelJob(id string) error {
	a.Lock()
	defer a.Unlock()
	if _, ok := a.Jobs[id]; ok {
		a.Jobs[id].Status = CANCELLED
	} else {
		return errors.New("Invalid ID")
	}
	return nil
}

func (a *JobQueue) GetAllJobs() Jobs {
	var res Jobs
	for key := range a.Jobs {
		var s string
		var err string
		switch a.Jobs[key].Status {
		case PENDING:
			s = "PENDING"
		case CANCELLED:
			s = "CANCELLED"
		case FINISHED:
			s = "FINISHED"
		case RUNNING:
			s = "RUNNING"
		case ERROR:
			s = "ERROR"
		}
		if a.Jobs[key].Error != nil {
			err = a.Jobs[key].Error.Error()
		}
		j := JobOut{Id: a.Jobs[key].Id, Snowflake: a.Jobs[key].Snowflake, Progress: a.Jobs[key].Progress, Status: s, Error: err}
		res.Jobs = append(res.Jobs, j)
	}
	return res
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

//{channel1: [channel1job, channel1job], channel2:[channel2job, channel2job]}
