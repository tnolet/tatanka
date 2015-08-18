package control

import (
	"github.com/tnolet/tatanka/work"
	"log"
)

func (c *Controller) StartWork() {

	var err error

	// Dispatcher
	doneChannel := make(chan work.WorkItem, c.state.WorkerAmount)
	WorkerQueue := make(chan chan work.WorkItem)

	for i := 0; i < c.state.WorkerAmount; i++ {

		log.Println("Starting worker", i+1)

		worker := work.New(i+1, WorkerQueue, doneChannel)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-c.workChan:
				worker := <-WorkerQueue
				worker <- work
			}
		}
	}()

	// Collector
	log.Println("Starting work dispatcher and collecting work...")
	c.collector = work.NewWorkCollector(c.state.QueueUrl)

	workPackages, err := c.collector.GetWork()

	if err != nil {
		log.Println(err.Error())
	}

	// create a map of all work items and mark then as new
	for _, pkg := range workPackages {
		for _, item := range pkg.WorkItems {
			c.workMap[item] = "todo"
		}
	}

	// shovel work into the work channel
	go func() {
		for k, _ := range c.workMap {
			if c.moreWork {
				c.workChan <- k
			} else {
				close(c.workChan)
			}
		}
	}()

	// mark as done only when job is finished
	go func() {
		for {
			select {
			case item := <-doneChannel:
				c.workMap[item] = "done"
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

	var items []work.WorkItem
	var pkgs []*work.WorkPackage

	c.moreWork = false

	log.Println("Saving unfinished work...")

	// filter out all the work that is still todo
	for k, v := range c.workMap {
		if v == "todo" {
			items = append(items, k)
		}
	}

	pkg := &work.WorkPackage{WorkItems: items}
	pkgs = append(pkgs, pkg)

	if err := c.collector.PutWork(pkgs); err != nil {
		log.Println("Error putting work:", err.Error())
	}

	c.ctrlChan <- &StartEvac{}

}
