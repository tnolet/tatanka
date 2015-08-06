package control

import (
	"errors"
	"github.com/tnolet/tatanka/bidder"
	"github.com/tnolet/tatanka/instance"
	"github.com/tnolet/tatanka/mailer"
	"github.com/tnolet/tatanka/store"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Controller struct {
}

func (c *Controller) Start(s store.State, m mailer.Mailer) {

	/*
		Controller start up sequence:
		1. Choose random region
		2. Clear all spot requests.
		3. Get on demand price.
		4. Get spot price history.
		5. Calculate our bid price.
			a. If bid price is too high, exit the whole routine
		6. Create spot request
		  a. Get AMI for that region
		  b. Issue request for time window
		7. Schedule killer for lifetime.
		8. Start evacuation watcher
		9. Send mail update

	*/

	region := instance.GetRandomRegion(s.Regions)

	log.Println("Initializing bidder for region: ", region)
	bidder := bidder.New(s.PriceListUrl, region)

	log.Println("Cancelling outstanding spot requests")
	bidder.CancelSpotRequests()

	ondemandPrice, err := bidder.GetOnDemandPrice(s.InstanceType, s.InstanceSize)
	if err != nil {
		log.Println("Error getting price: ", err.Error())
	}

	log.Println("On Demand price is: ", ondemandPrice)

	spotPrice, err := bidder.GetSpotPriceHistory(s.InstanceSize)
	if err != nil {
		log.Println("Error getting spot price history: ", err.Error())
	}

	log.Println("Spot price is: ", spotPrice)

	bidPrice, err := c.CalculateBidPrice(ondemandPrice, spotPrice, s.BidRatio)
	if err != nil {
		log.Println("Error calculating bid price: ", err.Error())
	}

	log.Println("Bid price is: ", bidPrice)

	amiID := instance.GetLinuxAMIforRegion(region)
	from, till := c.GetValidBidWindow(s.LifeTime, s.BidWindow)
	_, err = bidder.CreateSpotRequest(bidPrice, s.InstanceSize, amiID, from, till)

	if err != nil {
		log.Println("Error creating spot request: ", err.Error())
	}
}

func (c *Controller) GetValidBidWindow(lifeTime int, window int) (from time.Time, till time.Time) {
	now := time.Now()
	from = now.Add(time.Duration(lifeTime-window) * time.Second)
	till = now.Add(time.Duration(lifeTime) * time.Second)
	return from, till
}

/*

Rules regarding price: just calculate the spotprice + the ration in percentage, i,e:

spotprice = 0.5
ratio = 35%

bidprice = 0.5 * 1.35

however, the resulting bidprice cannot be higher than or equal to the ondemand price

*/
func (c *Controller) CalculateBidPrice(ondemandPrice string, spotPrice string, ratio int) (bidprice string, err error) {

	o, _ := strconv.ParseFloat(ondemandPrice, 64)
	s, _ := strconv.ParseFloat(spotPrice, 64)
	r := float64(ratio)
	bp := (s * (1 + (r / 100)))

	if bp >= o {
		bidprice = strconv.FormatFloat(bp, 'f', 3, 64)
		return bidprice, errors.New("Bid price higher than on demand price: " + bidprice)
	}

	bidprice = strconv.FormatFloat(bp, 'f', 3, 64)
	return bidprice, nil
}

/*
Termination checker loops and checks for termination every 5 seconds.
It also takes in the total lifetime.
If the lifetime is reached, or the instance is scheduled to be terminated,
a true value will be send over the evac channel.
*/
func (c *Controller) StartTerminationChecker(terminationUrl string, lifeTime int, evacChan chan bool) {

	now := time.Now()
	evacTime := now.Add(time.Duration(lifeTime) * time.Second)
	ticker := time.NewTicker(5 * time.Second)
	evacNotice := false

	log.Printf("Start death watch with life time till %v", evacTime)

	go func() {
		for evacNotice == false {
			select {
			case <-ticker.C:
				if c.instanceToBeTermninated(terminationUrl) || time.Now().After(evacTime) {
					evacNotice = true
					evacChan <- evacNotice
				} else {
					evacChan <- evacNotice
				}
			}
		}
	}()
}

// return true if AWS has scheduled the instance to be terminated
func (c *Controller) instanceToBeTermninated(terminationUrl string) bool {

	res, err := http.Get(terminationUrl)
	if err != nil {
		return false
	}
	if res.StatusCode == 200 {
		return true
	} else {
		return false
	}
}
