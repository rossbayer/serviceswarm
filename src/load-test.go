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
	scenario.Consumers = 2
	//scenario.CreationRate = swarm.NewRate(25, 1000 * time.Millisecond)

	//task := tasks.NewSampleTask(1)
	//scenario.AddTask(task)
	//task2 := tasks.NewSampleTask(2)
	//scenario.AddTask(task2)
	task := &tasks.SampleAggregateTask{}
	scenario.AddTask(task)

	scenario.Run()
}
