package control

import (
	"github.com/tnolet/tatanka/compute"
	"log"
	"time"
)

/*
Termination checker loops and checks for termination every 5 seconds.
It also takes in the total lifetime.
If the lifetime is reached, or the instance is scheduled to be terminated,
a true value will be send over the evac channel.
*/
func (c *Controller) StartDeathWatch() {

	now := time.Now()
	evacTime := now.Add(time.Duration(c.state.LifeTime) * time.Second)
	ticker := time.NewTicker(5 * time.Second)
	evacNotice := false

	log.Printf("Starting death watch with life time till %v", evacTime)

	go func() {
		for evacNotice == false {
			select {
			case <-ticker.C:
				if compute.InstanceToBeTerminated(c.state.TerminationUrl) || time.Now().After(evacTime) {
					evacNotice = true
					c.ctrlChan <- &StartEvac{}
				} else {
					log.Println("Deathwatch ticker: OK")
				}
			}
		}
	}()
	return
}
