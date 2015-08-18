package work

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"time"
)

func NewWorkCollector(url string, pkgChan chan WorkPackage) *WorkCollector {

	svc := sqs.New(nil)
	return &WorkCollector{svc: svc, url: url, pkgChan: pkgChan}

}

func (w *WorkCollector) Start() {

	log.Println("Initializing work collector for queue: ", w.url)
	go w.getWork()

}

func (w *WorkCollector) getWork() {

	var maxMessages int64 = 1
	var visTimeout int64 = 10
	var waitTimeout int64 = 15

	for {

		params := &sqs.ReceiveMessageInput{
			QueueUrl:            &w.url,
			MaxNumberOfMessages: &maxMessages,
			VisibilityTimeout:   &visTimeout,
			WaitTimeSeconds:     &waitTimeout,
		}

		resp, err := w.svc.ReceiveMessage(params)
		if err != nil {
			log.Println("Error getting message from collector queue: " + err.Error())
		}

		// compile work package, delete message from queue then send out the package
		for _, msg := range resp.Messages {
			log.Println("Got work message with id:", *msg.MessageId)

			pkg := parseMessage(*msg.Body)

			if err := w.DeleteMessage(*msg.ReceiptHandle); err != nil {
				log.Println(err.Error())
			}

			w.pkgChan <- *pkg

			time.Sleep(15 * time.Second)

		}
	}
}

func (w *WorkCollector) PutWork(workPackages []*WorkPackage) (err error) {

	for _, pkg := range workPackages {
		_msg, err := json.Marshal(pkg)
		if err != nil {
			return err
		}
		msg := string(_msg)

		params := &sqs.SendMessageInput{
			MessageBody: &msg,
			QueueUrl:    &w.url,
		}

		resp, err := w.svc.SendMessage(params)
		if err != nil {
			return errors.New("Error putting message to collector queue: " + err.Error())
		}

		log.Println("Message ID for created message is:", *resp.MessageId)
	}
	return nil
}

func parseMessage(msg string) *WorkPackage {

	var workPackage WorkPackage
	if err := json.Unmarshal([]byte(msg), &workPackage); err != nil {
		log.Println("Error parsing message to work package", err.Error())
	}

	return &workPackage
}

func (w *WorkCollector) DeleteMessage(receiptHandle string) error {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      &w.url,
		ReceiptHandle: &receiptHandle,
	}
	_, err := w.svc.DeleteMessage(params)
	if err != nil {
		return errors.New("Error deleting message from collector queue: " + err.Error())
	}
	return nil
}
