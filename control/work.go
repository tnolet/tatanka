package control

import (
	"github.com/tnolet/tatanka/work"
	"log"
)

func (c *Controller) StartWork() {

	defer func() {
		recover()
	}()

	// Dispatcher
	doneChannel := make(chan bool)
	WorkerQueue := make(chan chan work.WorkItem)

	for i := 0; i < c.state.WorkerAmount; i++ {

		log.Println("Starting worker", i+1)

		worker := work.New(i+1, WorkerQueue, doneChannel)
		worker.Start()
	}

	/*
		1. A work item comes in from the collector
		2. A worker is grabbed from the workers queue
		3. The work is send to worker channel
	*/
	go func() {
		for {
			select {
			case work := <-c.workChan:
				workerChan := <-WorkerQueue
				workerChan <- work
			}
		}
	}()

	wrkItmChan := make(chan work.WorkItem, 100)
	c.collector = work.NewWorkCollector(c.state.QueueUrl, wrkItmChan)
	c.collector.Start()

	// grab work items
	go func() {
		for {
			select {
			case itm := <-wrkItmChan:
				c.workChan <- itm
			}
		}
	}()

	// mark as done only when job is finished
	go func() {
		for {
			select {
			case _ = <-doneChannel:
				log.Println("job done")
			}
		}
	}()

}

/*
	Orders all workers to quit, and then wait for it to actually happen.
  If all workers call Done(), store the left over work and then proceed
  with the evacuation.
*/
func (c *Controller) StopWork() {

	log.Println("Stopping workers...")

	// var items []work.WorkItem

	// log.Println("Saving unfinished work...")

	// // filter out all the work that is still todo
	// for k, v := range c.workMap {
	// 	if v == "todo" {
	// 		items = append(items, k)
	// 	}
	// }

	// // only save it if there is anything to save
	// if len(items) > 0 {
	// 	pkg := &work.WorkPackage{WorkItems: items}
	// 	pkgs = append(pkgs, pkg)

	// 	if err := c.collector.PutWork(pkgs); err != nil {
	// 		log.Println("Error putting work:", err.Error())
	// 	}
	// }

	c.ctrlChan <- &StartEvac{}

}
