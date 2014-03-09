package swarm

type Consumer struct {
	ID       int
	Scenario *Scenario
}

func NewConsumer(id int, scenario *Scenario) *Consumer {
	return &Consumer{ID: id, Scenario: scenario}
}
