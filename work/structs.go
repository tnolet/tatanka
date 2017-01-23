package work

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Collector interface {
	GetWork()
	PutWork(string)
}

type WorkCollector struct {
	svc        *sqs.SQS
	url        string
	wrkItmChan chan WorkItem
}

type Worker struct {
	ID          int
	Work        chan WorkItem
	WorkerQueue chan chan WorkItem
	DoneChan    chan WorkItem
}

type WorkItem struct {
	Setup struct {
		Commands [][]string `json:"commands"`
	} `json:"setup"`

	Work struct {
		Commands [][]string `json:"commands"`
	} `json:"work"`

	Teardown struct {
		Commands [][]string `json:"commands"`
	}

	Variables map[string]string `json:"variables"`
}
