package work

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Collector interface {
	GetWork()
	PutWork(string)
}

type WorkCollector struct {
	svc *sqs.SQS
	url string
}

type WorkPackage struct {
	WorkItems []WorkItem `json:"subjects"`
}

type Worker struct {
	ID          int
	Work        chan WorkItem
	WorkerQueue chan chan WorkItem
	DoneChan    chan WorkItem
}

type WorkItem string
