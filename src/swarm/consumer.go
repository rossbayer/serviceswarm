package swarm

import (
	"math/rand"
	"time"
)

type Consumer struct {
	ID       int
	Scenario *Scenario
	executor *Executor
	Finished bool
}

func NewConsumer(id int, scenario *Scenario) *Consumer {
	return &Consumer{ID: id, Scenario: scenario}
}

func (c *Consumer) Start(tasks []Task) {
	copiedTasks := make([]Task, len(tasks))
	for i, task := range tasks {
		copiedTask := task.Copy()
		copiedTasks[i] = copiedTask

		// Perform setup on copied task if necessary
		withSetup, ok := copiedTask.(WithSetup)
		if ok {
			withSetup.Setup(c)
		}
	}

	c.executor = NewExecutor(copiedTasks, c.Scenario.MinWait, c.Scenario.MaxWait)

	go func() {
		c.executor.Exec()
		c.Finished = true
	}()
}

func (c *Consumer) Executor() *Executor {
	return c.executor
}

type Executor struct {
	tasks []Task
	rand *rand.Rand
	minWait time.Duration
	maxWait time.Duration
}

func NewExecutor(tasks []Task, minWait time.Duration, maxWait time.Duration) *Executor {
	return &Executor{
		tasks: tasks,
		rand: rand.New(rand.NewSource(time.Now().Unix())),
		minWait: minWait,
		maxWait: maxWait,
	}
}

func (e *Executor) AddTask(task Task) {
	copiedTask := task.Copy()
	e.tasks = append(e.tasks, copiedTask)
}

func (e *Executor) Exec() {
	for e.ExecNext() {}
}

func (e *Executor) ExecNext() bool {
	var task Task
	// Select a task based on its weight
	sum := int32(0)
	for _, task = range e.tasks {
		sum += int32(task.Weight())
	}

	random_n := int(e.rand.Int31n(sum))
	for _, task = range e.tasks {
		if random_n < task.Weight() {
			break
		}

		random_n -= task.Weight()
	}

	// Execute the task
	var result bool
	aggregateTask, ok := task.(AggregateTask)
	if ok {
		childExecutor := NewExecutor(make([]Task, 0, 10), e.minWait, e.maxWait)
		result = aggregateTask.ExecTasks(childExecutor)
	} else {
		result = task.Exec()
	}

	if result {
		e.wait()
	}

	// Teardown the task if it is complete
	withTeardown, ok := task.(WithTeardown)
	if !result && ok {
		withTeardown.Teardown()
	}

	return result
}

func (e *Executor) wait() {
	sleepTime := time.Duration(e.rand.Int63n(int64(e.maxWait)))
	if sleepTime < e.minWait {
		sleepTime = e.minWait
	}
	time.Sleep(sleepTime)
}
