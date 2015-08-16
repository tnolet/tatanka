package work

import (
	"log"
)

func StartDispatcher(workChan chan WorkPackage, workers int) {

	WorkerQueue := make(chan chan WorkPackage, workers)

	for i := 0; i < workers; i++ {
		log.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-workChan:
				log.Println("Received work requeust")
				go func() {
					worker := <-WorkerQueue

					log.Println("Dispatching work request")
					worker <- work
				}()
			}
		}
	}()
}
