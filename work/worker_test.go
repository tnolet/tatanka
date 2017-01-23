package work

import (
	"testing"
	// "time"
)

// test creating multiple workers
func TestStart(t *testing.T) {

	doneChan := make(chan bool)
	WorkerQueue := make(chan chan WorkItem)

	worker1 := New(1, WorkerQueue, doneChan)
	worker2 := New(2, WorkerQueue, doneChan)

	worker1.Start()
	worker2.Start()

	// construct a work item
	var workItm WorkItem
	setup := [][]string{{"docker", "pull", "tnolet/scraper:0.1.0"}}
	work := [][]string{{"scrapy", "crawl", "abc", "-a", "subject=$fullName", "-a", "query=$fullNameWithTransform"}}

	workItm.Setup.Commands = setup
	workItm.Work.Commands = work

	worker := <-WorkerQueue
}
