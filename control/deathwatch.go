package control

import (
	"github.com/tnolet/tatanka/compute"
	"log"
	"time"
)

/*
Termination checker loops and checks for termination every 5 seconds.

It also takes into account the total lifetime. If the lifetime is reached, or the instance is
scheduled to be terminated, a StartEvac message is send over the control channel.

A lifetime of 0 or less is equal to an infinite lifetime.

Every 12 checks (1 minute) it will log its state.
Every 720 check (1 hour) it will send a short update email

*/
func (c *Controller) StartDeathWatch() {

	now := time.Now()
	evacTime := now.Add(time.Duration(c.state.LifeTime) * time.Second)
	ticker := time.NewTicker(5 * time.Second)
	evacNotice := false
	counter := 0

	log.Printf("Starting death watch with life time till %v", evacTime)

	go func() {
		for evacNotice == false {
			select {
			case <-ticker.C:
				counter += 1
				if compute.InstanceToBeTerminated(c.state.TerminationUrl) || time.Now().After(evacTime) {
					evacNotice = true
					c.ctrlChan <- &StartEvac{}
				} else {
					if counter%12 == 0 {
						log.Println("Deathwatch ticker: OK")
					}
					if counter%720 == 0 {
						uptime := int(time.Duration.Seconds(time.Since(c.state.StartTime)))
						c.mailChan <- CasualUpdateMail(uptime, c.state.CurrentRegion)
					}
				}
			}
		}
	}()
	return
}
