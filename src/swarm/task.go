package swarm

type Task interface {
	Name() string
	Exec()
}

type WithSetup interface {
	Setup()
}

type WithTeardown interface {
	Teardown()
}
