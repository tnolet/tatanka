package work

import (
	"log"
)

func StartDispatcher(workChan chan WorkItem, workers int) {

	WorkerQueue := make(chan chan WorkItem, workers)

	for i := 0; i < workers; i++ {
		log.Println("Starting worker", i+1)
		worker := New(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-workChan:
				go func() {
					worker := <-WorkerQueue
					worker <- work
				}()
			}
		}
	}()
}
