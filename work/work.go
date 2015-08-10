package work

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

type Worker struct {
	initScript string
	workQueue *WorkQueue
}

func(w *Worker) Start() {

}

