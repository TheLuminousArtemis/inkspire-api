package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"time"

	"text/template"

	gomail "gopkg.in/mail.v2"
)

type mailtrapClient struct {
	username  string
	password  string
	fromEmail string
}

func NewMailTrapClient(username, fromEmail, password string) (mailtrapClient, error) {
	if username == "" || password == "" {
		return mailtrapClient{}, errors.New("username and password must be provided")
	}

	return mailtrapClient{
		username:  username,
		password:  password,
		fromEmail: fromEmail,
	}, nil
}

func (m mailtrapClient) Send(templateFile, username, email string, data any, isProdEnv bool) (int, error) {
	// Template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())

	message.AddAlternative("text/html", body.String())
	err = m.run(message, isProdEnv)
	if err != nil {
		return -1, err
	}
	return 200, nil
}

func (m mailtrapClient) run(message *gomail.Message, isProdEnv bool) error {
	if isProdEnv {
		dialer := gomail.NewDialer("sandbox.smtp.mailtrap.io", 2525, m.username, m.password)
		for i := 0; i < maxRetries; i++ {
			if err := dialer.DialAndSend(message); err != nil {
				log.Printf("Retry %d: Failed to send email: %v", i+1, err)
				time.Sleep(time.Second * time.Duration(i+1))
				continue
			}
			return nil
		}
		return fmt.Errorf("failed to send email")
	}
	return nil
}
