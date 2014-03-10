package tasks

import (
	"log"
	swarmlog "swarm/log"
	"swarm"
	"swarm/http"
	"strconv"
)

type SampleAggregateTask struct {
	swarm.AggregateBase
	log *log.Logger
}

func (a *SampleAggregateTask) Name() string {
	return "aggregate-task"
}

func (a *SampleAggregateTask) Weight() int {
	return 100
}

func (a *SampleAggregateTask) Setup(consumer *swarm.Consumer) {
	a.log = swarmlog.NewLogger("SampleAggregateTask")
}

func (a *SampleAggregateTask) ExecTasks(exec *swarm.Executor) bool {
	a.log.Println("Executing subtasks")
	sampleTask := NewSampleTask(10)
	exec.AddTask(sampleTask)
	exec.Exec()
	return false
}

func (a *SampleAggregateTask) Copy() swarm.Task {
	return &SampleAggregateTask{log: a.log}
}

type SampleTask struct {
	log *log.Logger
	client *http.HTTPClient
	count int
}

func NewSampleTask(id int) *SampleTask {
	return &SampleTask{log: swarmlog.NewLogger("SampleTask " + strconv.Itoa(id))}
}

func (s *SampleTask) Name() string {
	return "SampleTask"
}

func (s *SampleTask) Weight() int {
	return 5
}

func (s *SampleTask) Copy() swarm.Task {
	copy := &SampleTask{}
	copy.log = s.log
	copy.client = s.client

	return copy
}

func (s *SampleTask) Setup(consumer *swarm.Consumer) {
	s.log = swarmlog.NewLogger("sample-task")
	s.log.Println("Doing setup...")
	s.client = http.NewClient(consumer)
}

func (s *SampleTask) Exec() bool {
   	s.log.Println("Exec...")
	s.count++
	return s.count < 5
}

func (s *SampleTask) Teardown() {
	s.log.Println("Tearing down...")
}
