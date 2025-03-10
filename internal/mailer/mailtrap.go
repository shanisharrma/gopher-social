package mailer

import (
	"bytes"
	"errors"
	"html/template"

	"gopkg.in/gomail.v2"
)

type MailtrapClient struct {
	FromEmail string
	ApiKey    string
	Username  string
}

func NewMailtrap(apiKey, username, fromEmail string) (*MailtrapClient, error) {
	if apiKey == "" {
		return &MailtrapClient{}, errors.New("api key is required")
	}

	return &MailtrapClient{
		FromEmail: fromEmail,
		Username:  username,
		ApiKey:    apiKey,
	}, nil
}

func (m *MailtrapClient) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	// Template parsing and building
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
	message.SetHeader("From", m.FromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())

	message.AddAlternative("text/html", body.String())

	dialer := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, m.Username, m.ApiKey)
	if err := dialer.DialAndSend(message); err != nil {
		return -1, err
	}

	return 200, nil
}
