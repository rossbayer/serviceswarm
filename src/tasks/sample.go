package tasks

import (
	"log"
	swarmlog "swarm/log"
	"swarm"
	"swarm/http"
)

type SampleTask struct {
	log *log.Logger
	client *http.HTTPClient
}

func (s *SampleTask) Name() string {
	return "SampleTask"
}

func (s *SampleTask) Setup(consumer *swarm.Consumer) {
	s.log = swarmlog.NewLogger("sample-task")
	s.log.Println("Doing setup...")
	s.client = http.NewClient(consumer)
}

func (s *SampleTask) Exec() {
   	s.log.Println("Exec...")
}

func (s *SampleTask) Teardown() {
	s.log.Println("Tearing down...")
}
