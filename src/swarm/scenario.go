package swarm

import (
	"log"
	"math/rand"
	swarmlog "swarm/log"
	"time"
)

type Rate struct {
	Count  int
	Period time.Duration
}

func NewRate(count int, period time.Duration) *Rate {
	return &Rate{Count: count, Period: period}
}

type Scenario struct {
	MinWait      time.Duration
	MaxWait      time.Duration
	Consumers    int
	CreationRate *Rate
	Results      *Results

	log            *log.Logger
	rand           *rand.Rand
	lastConsumerID int
	consumers      []*Consumer
	tasks          []Task
	stopping       bool
}

func NewScenario() *Scenario {
	return &Scenario{
		MinWait:        1 * time.Second,
		MaxWait:        5 * time.Second,
		Consumers:      1,
		rand:           rand.New(rand.NewSource(time.Now().Unix())),
		lastConsumerID: 0,
		consumers:      make([]*Consumer, 0, 1000),
		tasks:          make([]Task, 0, 1000),
	}
}

func (s *Scenario) AddTask(task Task) {
	s.tasks = append(s.tasks, task)
}

func (s *Scenario) Run() {
	s.log = swarmlog.NewLogger("swarm.scenario")
	s.log.Println("Starting the scenario...")

	// Verify configuration values are valid
	if(s.Consumers < 1) {
		s.log.Fatalln("Number of consumers must be 1 or more.")
		return
	}

	// If no creation rate is specified, then make it so all consumers
	// are generated at once.
	if(s.CreationRate == nil) {
		s.CreationRate = &Rate{Count: s.Consumers, Period: 1}
	}

	if(s.Consumers < s.CreationRate.Count) {
		s.log.Fatalln("Number of consumers must be greater than or equal to " +
			"number defined in creation rate.")
		return
	}

	// Start results collection
	s.Results = NewResults()

	// Start consumer generation
	go s.generateConsumers()
}

func (s *Scenario) Stop() {
	s.log.Println("Stopping scenario")
	s.stopping = true
	s.Results.Stop()
	s.log.Println("Scenario complete")
}

func (s *Scenario) Wait() {
	sleepTime := time.Duration(s.rand.Int63n(int64(s.MaxWait)))
	if sleepTime < s.MinWait {
		sleepTime = s.MinWait
	}
	time.Sleep(sleepTime)
}

func (s *Scenario) generateConsumers() {
	s.log.Println("Starting consumer generation")
	for len(s.consumers) < s.Consumers && !s.stopping {
		s.consumers = append(s.consumers, s.newConsumerGroup()...)
		time.Sleep(s.CreationRate.Period)
	}

	s.log.Println("Finished creating consumers")
}

func (s *Scenario) newConsumer() *Consumer {
	c := NewConsumer(s.lastConsumerID, s)
	s.lastConsumerID++
	return c
}

func (s *Scenario) newConsumerGroup() []*Consumer {
	totalToCreate := s.CreationRate.Count
	currentCount := len(s.consumers)
	if currentCount + totalToCreate > s.Consumers {
		totalToCreate = s.Consumers - currentCount
	}

	s.log.Printf("Creating %d consumers", totalToCreate)
	consumers := make([]*Consumer, totalToCreate)
	for i := 0; i < totalToCreate; i++ {
		consumers[i] = s.newConsumer()
	}

	return consumers
}
