package mail

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"log"
)

type Mailer struct {
	svc    *ses.SES
	params *ses.SendEmailInput
	mailChan  chan string
}

func New(toAddress string, region string, mailChan chan string) *Mailer {

	var m Mailer
	m.svc = ses.New(&aws.Config{Region: region})
	m.params = &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(toAddress),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data:    aws.String("Hi, I'm Tatanka"),
					Charset: aws.String("UTF-8"),
				},
			},
			Subject: &ses.Content{
				Data:    aws.String("Tatanka"),
				Charset: aws.String("UTF-8"),
			},
		},
		Source: aws.String("tim@magnetic.io"),
		ReplyToAddresses: []*string{
			aws.String("tim@magnetic.io"),
		},
		ReturnPath: aws.String("tim@magnetic.io"),
	}
	m.mailChan = mailChan
	return &m
}

func (m *Mailer) Start() {

	log.Printf("Initializing mailer...")

	go func() {
		for {
			select {
			case mail := <-m.mailChan:
				m.send(mail)
			}
		}
	}()
}

func (m *Mailer) send(body string) {

	log.Println("Sending mail...")

	m.params.Message.Body.Text.Data = aws.String(body)
	_, err := m.svc.SendEmail(m.params)
	if err != nil {
		log.Println("Error sending email...", err.Error())
	}
}
