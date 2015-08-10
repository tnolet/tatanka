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
}

func New(m chan string, c chan Message, s chan store.State) *Controller {
	return &Controller{mailChan: m, ctrlChan: c, stateChan: s}
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

func (c *Controller) Init() {

	log.Println("Starting control initialization sequence...")

	c.prepLocalStateOnInit()
	c.bidder = bidder.New(c.state.PriceListUrl, c.state.CurrentBidRegion)
	c.bidder.CancelSpotRequests()

	bidPrice, err := c.GetBidPrice()
	if err != nil {
		c.mailChan <- FatalErrorMail("Error getting a bid price: " + err.Error())
		log.Fatal(err.Error())
	}

	amiID := compute.GetLinuxAMIforRegion(c.state.CurrentBidRegion)
	from, till := c.GetValidBidWindow(c.state.LifeTime, c.state.BidOffset, c.state.BidWindow)
	reqs, err := c.bidder.CreateSpotRequest(bidPrice, c.state.InstanceSize, amiID, from, till)

	if err != nil {
		c.mailChan <- FatalErrorMail("Error creating a spot request: " + err.Error())
		log.Fatal("Error creating spot request: ", err.Error())
	}

	c.state.LastReqID = c.state.CurrentReqID
	c.state.CurrentReqID = reqs[0].Id

	c.mailChan <- InitMail(c.state.CurrentInstanceID,
		c.state.CurrentRegion,
		c.state.LifeTime,
		c.state.CurrentBidRegion,
		c.state.CurrentReqID)
}

func (c *Controller) Evac() {

	log.Println("Starting evacuation sequence")

	/*	Before evac, check if the spot request created during INIT has already been fulfilled.
		If this is not the case, cancel it and create a brand new one.
	*/
	if c.bidder.SpotInstanceActive(c.state.CurrentReqID) != true {
		c.bidder.CancelSpotRequests()

		bidPrice, err := c.GetBidPrice()
		if err != nil {
			c.mailChan <- FatalErrorMail("Error getting a bid price: " + err.Error())
			log.Fatal(err.Error())
		}
		amiID := compute.GetLinuxAMIforRegion(c.state.CurrentBidRegion)

		// make bid valid from + 1 minute till the bidwindow.
		from, till := c.GetValidBidWindow(60, 0, c.state.BidWindow)
		reqs, err := c.bidder.CreateSpotRequest(bidPrice, c.state.InstanceSize, amiID, from, till)

		if err != nil {
			c.mailChan <- FatalErrorMail("Error creating a spot request before evac: " + err.Error())
			log.Fatal("Error creating spot request: ", err.Error())
		}

		c.state.LastReqID = c.state.CurrentReqID
		c.state.CurrentReqID = reqs[0].Id
	}

	c.prepLocalStateOnEvac()
	c.mailChan <- EvacuationMail(c.state.CurrentInstanceID, c.state.CurrentRegion, c.state.CurrentBidRegion)
	c.stateChan <- c.state

	// terminate
	comp := compute.New(c.state.CurrentRegion)
	_, err := comp.TerminateInstance(c.state.CurrentInstanceID)

	if err != nil {
		log.Println(err.Error())
	}
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
