package work

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

func NewWorkCollector(url string) *WorkCollector {
	log.Println("Initializing work collector for queue: ", url)

	svc := sqs.New(nil)
	return &WorkCollector{svc: svc, url: url}

}

func (w *WorkCollector) GetWork() (workPackages []*WorkPackage, err error) {

	var maxMessages int64 = 1
	var visTimeout int64 = 10
	var waitTimeout int64 = 15

	params := &sqs.ReceiveMessageInput{
		QueueURL:            &w.url,
		MaxNumberOfMessages: &maxMessages,
		VisibilityTimeout:   &visTimeout,
		WaitTimeSeconds:     &waitTimeout,
	}

	resp, err := w.svc.ReceiveMessage(params)
	if err != nil {
		return workPackages, errors.New("Error getting message from collector queue: " + err.Error())
	}

	// compile work package and delete message from queue
	for _, msg := range resp.Messages {
		log.Println("Got work message with id:", *msg.MessageID)
		workPackages = append(workPackages, parseMessage(*msg.Body))

		if err := w.DeleteMessage(*msg.ReceiptHandle); err != nil {
			log.Println(err.Error())
		}
	}

	return workPackages, nil
}

func (w *WorkCollector) PutWork(workPackages []*WorkPackage) (err error) {

	_msg, err := json.Marshal(workPackages)
	if err != nil {
		return err
	}

	msg := string(_msg)
	log.Println("saving:", msg)

	params := &sqs.SendMessageInput{
		MessageBody: &msg,
		QueueURL:    &w.url,
	}

	resp, err := w.svc.SendMessage(params)
	if err != nil {
		return errors.New("Error putting message to collector queue: " + err.Error())
	}

	log.Println("Message ID for created message is:", *resp.MessageID)
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
