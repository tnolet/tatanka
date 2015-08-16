package control

import (
	"github.com/tnolet/tatanka/compute"
	"log"
	"os"
)

func (c *Controller) Evac() {

	log.Println("Starting evacuation sequence")

	/*  Before evac, check if the spot request created during INIT has already been fulfilled.
	    If this is not the case, cancel it and create a brand new one.
	*/

	if !c.noop {
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
	} else {
		os.Exit(0)
	}
}
