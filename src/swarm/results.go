package swarm

import (
	"log"
	swarmlog "swarm/log"
	"time"
)

type TaskStatus string

const (
	Error   TaskStatus = "error"
	Success TaskStatus = "success"
	Fatal   TaskStatus = "fatal"
)

type Result struct {
	Status    TaskStatus
	StartTime time.Time
	Duration  time.Duration
	TaskName  string
	Consumer  *Consumer
	Details   interface{}
}

type Results struct {
	resultsChan chan *Result
	stopChan    chan bool
	log         *log.Logger
}

func NewResults() *Results {
	results := &Results{
		resultsChan: make(chan *Result, 1024),
		stopChan:    make(chan bool),
		log:         swarmlog.NewLogger("swarm.results"),
	}

	go results.collect()

	return results
}

func (r *Results) Add(result *Result) {
	r.resultsChan <- result
}

func (r *Results) Stop() {
	r.log.Println("Stopping results collection")
	r.stopChan <- true
	r.log.Println("Waiting for collection to complete")
	<-r.stopChan
}

func (r *Results) collect() {
	r.log.Println("Starting collection of results")

	var (
		stop   bool
		result *Result
	)

	for {
		select {
		case stop = <-r.stopChan:
			r.log.Println("Received stop signal")
			continue
		case result = <-r.resultsChan:
			r.log.Printf("Received result: %d", result.StartTime)
		default:
			if stop {
				r.stopChan <- true
				return
			}
		}
	}
}
