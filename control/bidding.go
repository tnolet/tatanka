package control

import (
	"errors"
	"log"
	"strconv"
	"time"
)

/*
	from = lifetime - bid offset
	till = from + bid window

*/

func (c *Controller) GetValidBidWindow(lifeTime int, offset int, window int) (from time.Time, till time.Time) {

	if (offset < lifeTime || lifeTime == 0 || window == 0) {
		log.Fatalf("Invalid lifetime %v, offset %v or window %v", lifeTime, offset, window)
		return from,till
	}

	now := time.Now()
	from = now.Add(time.Duration(lifeTime-offset) * time.Second)
	till = now.Add(time.Duration((lifeTime-offset)+window) * time.Second)
	log.Printf("Bid window ranges from %v till %v.", from.Format(time.RFC3339), till.Format(time.RFC3339))
	return from, till
}

/*

Rules regarding price: just calculate the spotprice + the ration in percentage, i,e:

spotprice = 0.5
ratio = 35%

bidprice = 0.5 * 1.35

however, the resulting bidprice cannot be higher than or equal to the ondemand price

*/
func (c *Controller) GetBidPrice() (bidprice string, err error) {

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

func (c *Controller) CalculateBidPrice(ondemandPrice float64, spotPrice float64, ratio float64) (bidprice string, err error) {
	_bp := (spotPrice * (1 + (ratio / 100)))

	if _bp >= ondemandPrice {
		bidprice = strconv.FormatFloat(_bp, 'f', 3, 64)
		return bidprice, errors.New("Bid price higher than on demand price: " + bidprice)
	}

	bidprice = strconv.FormatFloat(_bp, 'f', 3, 64)
	return bidprice, nil
}
