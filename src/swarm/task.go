package swarm

type Task interface {
	// The name of the task
	Name() string

	// Number representing the weight of the task vs. other tasks
	Weight() int

	// Execute the task and return bool indicating whether this task is
	// complete
	Exec() bool

	// Generate a duplicate copy of this task for use by an individual
	// Consumer.  A new instance of a task will be created for each
	// consumer that executes it.
	Copy() Task
}

type WithSetup interface {
	Setup(*Consumer)
}

type WithTeardown interface {
	Teardown()
}

type AggregateTask interface {
	ExecTasks(exec *Executor) bool
}

type AggregateBase struct {}

func (a *AggregateBase) Exec() bool {
	return false
}
