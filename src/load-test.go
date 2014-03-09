package main

import (
	"swarm"
	"tasks"
	"time"
)

func main() {
	scenario := swarm.NewScenario()
	scenario.MinWait = 1 * time.Second
	scenario.MaxWait = 3 * time.Second
	//scenario.Consumers = 100
	//scenario.CreationRate = swarm.NewRate(25, 1000 * time.Millisecond)

	task := &tasks.SampleTask{}
	scenario.AddTask(task)

	scenario.Run()

	time.Sleep(3000 * time.Millisecond)

	scenario.Stop()
}
