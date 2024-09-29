package email

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"loan-service/utils/errs"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService interface {
	DefaultSenderName() string
	DefaultSenderAddress() string
	SendMail(ctx context.Context, subject, body string, from mail.Email, to mail.Email, attachment *AttachmentOpts) error
}

func NewEmailService(apiKey string, defaultSenderAddress, defaultSenderName string) EmailService {
	return &emailService{
		client:               sendgrid.NewSendClient(apiKey),
		defaultSenderAddress: defaultSenderAddress,
		defaultSenderName:    defaultSenderName,
	}
}

type emailService struct {
	client               *sendgrid.Client
	defaultSenderAddress string
	defaultSenderName    string
}

// DefaultSenderAddress implements EmailService.
func (e *emailService) DefaultSenderAddress() string {
	return e.defaultSenderAddress
}

// DefaultSenderName implements EmailService.
func (e *emailService) DefaultSenderName() string {
	return e.defaultSenderName
}

type AttachmentType string

const (
	AttachmentTypeJPG = "image/jpeg"
	AttachmentTypePNG = "image/png"
	AttachmentTypePDF = "application/pdf"
)

type AttachmentOpts struct {
	File        io.Reader
	ContentType AttachmentType
	Filename    string
}

func (e *emailService) SendMail(ctx context.Context, subject, body string, from mail.Email, to mail.Email, attachmentOpts *AttachmentOpts) error {
	fmt.Printf("[email sending] from:%s to:%s subject:%s\n", from.Address, to.Address, subject)

	// create new mail with origin address and body
	newMail := mail.NewV3Mail()
	newMail.SetFrom(&from)
	newMail.AddContent(mail.NewContent("text/html", body))

	// add destination address and subject
	personalization := mail.NewPersonalization()
	personalization.AddTos(&to)
	personalization.Subject = subject

	newMail.AddPersonalizations(personalization)

	if attachmentOpts != nil {
		// read and attach attachment file
		attachment := mail.NewAttachment()
		attachmentBytes, err := io.ReadAll(attachmentOpts.File)
		if err != nil {
			fmt.Println("cannot read provided attachment")
			return errs.Wrap(ErrEmailNotSent)
		}

		// encode attachment
		encodedAttachment := base64.StdEncoding.EncodeToString([]byte(attachmentBytes))
		attachment.SetContent(encodedAttachment)
		attachment.SetType(string(attachmentOpts.ContentType))
		attachment.SetFilename(attachmentOpts.Filename)
		attachment.SetDisposition("attachment")

		newMail.AddAttachment(attachment)
	}

	// message := mail.NewSingleEmail(&from, subject, &t, body, "")
	response, err := e.client.SendWithContext(ctx, newMail)
	if err != nil {
		fmt.Println("failed to send email")
		return errs.Wrap(ErrEmailNotSent)
	}

	fmt.Printf("[email sent] %d - %#v - %#v\n", response.StatusCode, response.Body, response.Headers)
	fmt.Printf("[email response] %#v", *response)

	return nil
}
