package control

import (
	"github.com/tnolet/tatanka/bidder"
	"github.com/tnolet/tatanka/compute"
	"github.com/tnolet/tatanka/store"
	"github.com/tnolet/tatanka/work"
	"log"
	"time"
)

type Controller struct {
	state     store.State
	mailChan  chan string
	ctrlChan  chan Message
	stateChan chan store.State
	workChan  chan work.WorkItem       // use to send work items
	workMap   map[work.WorkItem]string // used to store the reservoir of work
	moreWork  bool                     // used as a general switch to start/stop work
	collector *work.WorkCollector
	bidder    *bidder.Bidder
	noop      bool
}

func New(mail chan string, c chan Message, s chan store.State, noop bool) *Controller {

	w := make(chan work.WorkItem)
	m := make(map[work.WorkItem]string)

	return &Controller{
		mailChan:  mail,
		ctrlChan:  c,
		stateChan: s,
		workChan:  w,
		workMap:   m,
		moreWork:  true,
		noop:      noop}
}

func (c *Controller) Start() {

	log.Printf("Initializing controller...")

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
				case "STOP_WORK":
					c.StopWork()
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
