package control

import (
	"github.com/tnolet/tatanka/bidder"
	"github.com/tnolet/tatanka/compute"
	"log"
)

// Initializes state and the bidder. Starts the Work routine at the end.
func (c *Controller) Init() {

	log.Println("Starting control initialization sequence...")

	if !c.noop {

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

	c.ctrlChan <- &StartWork{}
}
