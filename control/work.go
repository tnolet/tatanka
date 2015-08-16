package control

import (
	"github.com/tnolet/tatanka/work"
	"log"
)

func (c *Controller) StartWork() {

	work.StartDispatcher(c.workChan, 4)
	collector := work.NewWorkCollector(c.state.QueueUrl)

	workPackages, err := collector.GetWork()
	if err != nil {
		log.Println(err.Error())
	}

	for i := 0; i < 100; i++ {
		log.Printf("sending work package %d", i)
		c.workChan <- *workPackages[0]
	}

}
