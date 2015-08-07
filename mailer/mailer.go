package mailer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"log"
)

type Mailer struct {
	svc    *ses.SES
	params *ses.SendEmailInput
}

func New(email string, region string) *Mailer {

	var m Mailer
	m.svc = ses.New(&aws.Config{Region: region})
	m.params = &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(email),
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
	return &m
}

func (m *Mailer) Send(body string) {

	m.params.Message.Body.Text.Data = aws.String(body)
	_, err := m.svc.SendEmail(m.params)
	if err != nil {
		log.Println("Error sending email...", err.Error())
	}
}
