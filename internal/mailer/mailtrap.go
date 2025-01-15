package mailer

import (
	"bytes"
	"errors"
	"log"
	"text/template"

	gomail "gopkg.in/mail.v2"
)

type mailTrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(apiKey string, fromEmail string) (mailTrapClient, error) {
	if apiKey == "" {
		return mailTrapClient{}, errors.New("api key required")
	}
	return mailTrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m mailTrapClient) Send(templateFile, username, email string, data any, isSandBox bool) (int, error) {
	// TODO move this to another place for abstraction
	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return -1, err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())
	message.AddAlternative("text/html", body.String())

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey)

	// TODO apply retries and error handling
	if err := dialer.DialAndSend(message); err != nil {
		log.Printf("failed to send email: %s", err.Error())
		return -1, err
	}

	return 200, nil
}
