package control

import (
	"errors"
	"github.com/tnolet/tatanka/bidder"
	"log"
	"strconv"
	"time"
)

/*
	from = lifetime - bid offset
	till = from + bid window

*/

func (c *Controller) GetValidBidWindow(lifeTime int, offset int, window int) (from time.Time, till time.Time) {

	if offset > lifeTime || lifeTime == 0 || window == 0 {
		log.Fatalf("Invalid lifetime %v, offset %v or window %v", lifeTime, offset, window)
		return from, till
	}

	now := time.Now()
	from = now.Add(time.Duration(lifeTime-offset) * time.Second)
	till = now.Add(time.Duration((lifeTime-offset)+window) * time.Second)
	log.Printf("Bid window ranges from %v till %v.", from.Format(time.RFC3339), till.Format(time.RFC3339))
	return from, till
}

/*

Rules regarding price: just calculate the spotprice + the ration in percentage, i,e:
	- spotprice = 0.5
	- ratio = 35%
	- bidprice = 0.5 * 1.35
however, the resulting bidprice cannot be higher than or equal to the ondemand price
*/
func (c *Controller) GetBidPrice(strategy string) (bidprice string, err error) {

	if strategy == "simple" {
		return c.simpleStrategy()
	}

	if strategy == "advanced" {
		return c.advancedStrategy()
	}

	return bidprice, nil
}

func (c *Controller) simpleStrategy() (bidprice string, err error) {

	ondemandPrice, err := c.bidder.GetOnDemandPrice(c.state.InstanceType, c.state.InstanceSize)
	if err != nil {
		return bidprice, errors.New("Error getting price:" + err.Error())
	}

	log.Println("On Demand price is:", ondemandPrice)

	spotPrice, err := c.bidder.GetSpotPriceHistory(c.state.InstanceSize)
	if err != nil {
		return bidprice, errors.New("Error getting spot price history: " + err.Error())
	}

	log.Println("Spot price is:", spotPrice)

	o, _ := strconv.ParseFloat(ondemandPrice, 64)
	s, _ := strconv.ParseFloat(spotPrice, 64)
	r := float64(c.state.BidRatio)

	bidprice, err = c.CalculateBidPrice(o, s, r)
	if err != nil {
		return bidprice, err
	}
	return bidprice, nil
}

/*
	Advanced strategy:
-	make multiple bidders
- get prices from all regions in parallel

*/
func (c *Controller) advancedStrategy() (string, error) {

	log.Println("Using advanced pricing strategy: overriding standard bidder")

	var allRegionPrices = make(map[string]float64)
	var priceChan = make(chan map[string]float64)

	// loop over all regions and get the bid price in parallel
	for _, region := range c.state.Regions {
		go func(region string) {
			bidder := bidder.New("", region)
			price, _ := bidder.GetSpotPriceHistory(c.state.InstanceSize)

			// if err != nil {
			// 	return price, errors.New("Error getting spot price history: " + err.Error())
			// }
			s, _ := strconv.ParseFloat(price, 64)
			RegionPriceMap := map[string]float64{region: s}
			priceChan <- RegionPriceMap
		}(region)
	}

	// collect all the prices from all regions
	for i := 0; i < len(c.state.Regions); i++ {
		regionPrice := <-priceChan
		for region, price := range regionPrice {
			log.Printf("Price in %s at: %f", region, price)
			allRegionPrices[region] = price
		}
	}

	// get the best price
	var bestPrice float64 = 10000000000
	var bestRegion string
	for region, price := range allRegionPrices {
		if price <= bestPrice {
			bestPrice = price
			bestRegion = region
		}
	}

	log.Printf("Best price in %s at: %f", bestRegion, bestPrice)
	bidPrice, err := c.CalculateBidPrice(123023.0, bestPrice, 0.5)
	if err != nil {
		return bidPrice, err
	}
	return bidPrice, nil
}

func (c *Controller) CalculateBidPrice(ondemandPrice float64, spotPrice float64, ratio float64) (bidprice string, err error) {
	_bp := (spotPrice * (1 + (ratio / 100)))

	if _bp >= ondemandPrice {
		bidprice = strconv.FormatFloat(_bp, 'f', 3, 64)
		return bidprice, errors.New("Bid price higher than on demand price: " + bidprice)
	}

	bidprice = strconv.FormatFloat(_bp, 'f', 3, 64)
	return bidprice, nil
}
