package work

import (
	"log"
)

type Queue interface {
	GetMessage() string
	PutMessage(string) error
}

type WorkQueue struct {
	url string
}

func (w *WorkQueue) GetMessage() string {
	return "Got new work from queue"
}

func (w *WorkQueue) PutMessage(work string) error {
	log.Printf("putting work %v back in queue", work)
	return nil
}

func NewWorkQueue() *WorkQueue {
	return &WorkQueue{url: "http://lalala"}
}
