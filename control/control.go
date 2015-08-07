package control

import (
	"github.com/tnolet/tatanka/bidder"
	"github.com/tnolet/tatanka/compute"
	"github.com/tnolet/tatanka/mailer"
	"github.com/tnolet/tatanka/store"
	"log"
)

type Controller struct {
	state     store.State
	mailer    mailer.Mailer
	ctrlChan  chan Message
	bidder    *bidder.Bidder
	bidRegion string
}

func New(s store.State, m mailer.Mailer, c chan Message) *Controller {
	return &Controller{state: s, mailer: m, ctrlChan: c}
}

func (c *Controller) Start() {

	log.Printf("Starting controller...")

	go func() {
		for {
			select {
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

func (c *Controller) Init() {

	c.state.LastRegion = c.state.CurrentRegion
	c.state.CurrentRegion = compute.GetCurrentRegion()
	c.state.CurrentInstanceID = compute.GetCurrentInstanceID()

	c.bidRegion = compute.GetRandomRegion(c.state.Regions)
	c.bidder = bidder.New(c.state.PriceListUrl, c.bidRegion)
	c.bidder.CancelSpotRequests()

	bidPrice, err := c.GetBidPrice()
	if err != nil {
		log.Println(err.Error())
	}

	amiID := compute.GetLinuxAMIforRegion(c.bidRegion)
	from, till := c.GetValidBidWindow(c.state.LifeTime, c.state.BidOffset, c.state.BidWindow)
	_, err = c.bidder.CreateSpotRequest(bidPrice, c.state.InstanceSize, amiID, from, till)
	if err != nil {
		log.Println("Error creating spot request: ", err.Error())
	}
	c.mailer.Send(InitMail(c.state.CurrentInstanceID, c.state.CurrentRegion, c.state.LifeTime, c.bidRegion))
}

func (c *Controller) Evac() {

	log.Println("Starting evacuation sequence")

	c.bidder.CancelSpotRequests()

	bidPrice, err := c.GetBidPrice()
	if err != nil {
		log.Println(err.Error())
	}

	amiID := compute.GetLinuxAMIforRegion(c.bidRegion)

	// make bid valid from + 1 minute till the bidwindow.
	from, till := c.GetValidBidWindow(60, 0, c.state.BidWindow)
	_, err = c.bidder.CreateSpotRequest(bidPrice, c.state.InstanceSize, amiID, from, till)

	if err != nil {
		log.Println(err.Error())
	}

	c.mailer.Send(EvacuationMail(c.state.CurrentInstanceID, c.state.CurrentRegion, c.bidRegion))

	// terminate
	comp := compute.New(c.state.CurrentRegion)
	_, err = comp.TerminateInstance(c.state.CurrentInstanceID)

	if err != nil {
		log.Println(err.Error())
	}
}
