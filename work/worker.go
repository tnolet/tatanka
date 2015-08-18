package work

import (
	"log"
	// "time"
)

const (
	container = "tnolet/scraper:0.1.0"
)

func New(id int, workerQueue chan chan WorkItem, doneChan chan WorkItem) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkItem),
		WorkerQueue: workerQueue,
		DoneChan:    doneChan,
	}
	return worker
}

/*
	Workers take work from the queue, start a job and block till the work is finished.
	All workers can be stopped by w.QuitChan.
*/
func (w Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work
			select {
			case workItem := <-w.Work:

				log.Printf("Worker %d received job request: %v", w.ID, workItem)

				jobDone := make(chan bool)

				go runNoopScraper(jobDone)
				// go runScraper(done)
				<-jobDone
				log.Printf("Worker %d finished job: %v", w.ID, workItem)
				w.DoneChan <- workItem
			}
		}
	}()
}
