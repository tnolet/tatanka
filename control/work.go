package control

import (
	"github.com/tnolet/tatanka/work"
	"log"
)

func (c *Controller) StartWork() {

	work.StartDispatcher(c.workChan, 5)
	collector := work.NewWorkCollector(c.state.QueueUrl)

	workPackages, err := collector.GetWork()
	if err != nil {
		log.Println(err.Error())
	}

	/*
	    loop over received work packages and pops of
		   work items into the work channel.
	*/
	for _, pkg := range workPackages {
		count := 0
		for {
			item := pkg.Shift()
			count += 1
			if item == "" || count == 5 {
				break
			}
			c.workChan <- item
		}

	}
}
