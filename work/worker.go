package work

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

const (
	container = "tnolet/scraper:0.1.0"
)

func New(id int, workerQueue chan chan WorkItem) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkItem),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}

	return worker
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work
			select {
			case workItem := <-w.Work:
				log.Printf("worker%d: Received work request: %v", w.ID, workItem)
				runScraper(workItem)
				time.Sleep(4 * time.Second)
			case <-w.QuitChan:
				log.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func runScraper(item WorkItem) error {
	// scrapy crawl abc -a subject="My Subject" -a query="My+Subject"

	crawl := "crawl"
	site := "abc"
	arg := "-a"
	subject := "subject='" + string(item) + "'"
	query := "query=" + strings.Replace(string(item), " ", "-", -1)
	cmd := exec.Command(
		"/usr/local/bin/docker",
		"run",
		container,
		crawl,
		site,
		arg,
		subject,
		arg,
		query)

	err := cmd.Run()
	if err != nil {
		log.Println("command unsuccessful:", err.Error())
		return err
	}
	return nil
}
