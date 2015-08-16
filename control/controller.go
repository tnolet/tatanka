package control

import (
	"github.com/tnolet/tatanka/bidder"
	"github.com/tnolet/tatanka/compute"
	"github.com/tnolet/tatanka/store"
	"log"
	"time"
)

type Controller struct {
	state     store.State
	mailChan  chan string
	ctrlChan  chan Message
	stateChan chan store.State
	bidder    *bidder.Bidder
	noop      bool
}

func New(m chan string, c chan Message, s chan store.State, noop bool) *Controller {
	return &Controller{mailChan: m, ctrlChan: c, stateChan: s, noop: noop}
}

func (c *Controller) Start() {

	log.Printf("Initializing controller...")

	// only proceed once the state has loaded
	c.state = <-c.stateChan

	go func() {
		for {
			select {
			case state := <-c.stateChan:
				c.state = state
			case msg := <-c.ctrlChan:
				switch msg.Get() {
				case "INIT":
					c.Init()
				case "START_DEATH_WATCH":
					c.StartDeathWatch()
				case "START_WORK":
					c.StartWork()
				case "START_EVAC":
					c.Evac()
				}
			}
		}
	}()
}

func (c *Controller) State() *store.State {
	return &c.state
}

func (c *Controller) StartWork() {
	log.Println("start work")
}

// simple helper to encapsulate all boring init functions
func (c *Controller) prepLocalStateOnInit() {
	c.state.LastRegion = c.state.CurrentRegion
	c.state.CurrentRegion = compute.GetCurrentRegion()
	c.state.CurrentInstanceID = compute.GetCurrentInstanceID()
	c.state.LastLifeTimeTarget = c.state.LifeTime
	c.state.StartTime = time.Now()

	c.state.LastBidRegion = c.state.CurrentBidRegion
	c.state.CurrentBidRegion = compute.GetRandomRegion(c.state.Regions)
}

func (c *Controller) prepLocalStateOnEvac() {
	c.state.LastLifeTimeActual = int(time.Duration.Seconds(time.Since(c.state.StartTime)))
}
