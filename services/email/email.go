package email

import (
	"fmt"
	"loan-service/utils/errs"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService interface {
	SendMail(subject, body string, from mail.Email, to []mail.Email) error
}

func NewEmailService(apiKey string) EmailService {
	return &emailService{
		client: sendgrid.NewSendClient(apiKey),
	}
}

type emailService struct {
	client *sendgrid.Client
}

func (e *emailService) SendMail(subject, body string, from mail.Email, to []mail.Email) error {
	for _, t := range to {
		fmt.Printf("[email sending] from:%s to:%s subject:%s\n", from.Address, t.Address, subject)
		message := mail.NewSingleEmail(&from, subject, &t, body, "")
		response, err := e.client.Send(message)
		if err != nil {
			return errs.Wrap(ErrEmailNotSent)
		}

		fmt.Printf("[email sent] %d - %#v - %#v\n", response.StatusCode, response.Body, response.Headers)
	}

	return nil
}
