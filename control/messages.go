package control

type Message interface {
	Get() string
}

type Init struct{}

func (Init) Get() string { return "INIT" }

type StartDeathWatch struct{}

func (StartDeathWatch) Get() string { return "START_DEATH_WATCH" }

type StartWork struct{}

func (StartWork) Get() string { return "START_WORK" }

type StartEvac struct{}

func (StartEvac) Get() string { return "START_EVAC" }
