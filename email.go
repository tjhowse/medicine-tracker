package main

import (
	"fmt"
	"net/smtp"
)

type email struct {
	settings *ServerSettings
}

func GetEmail(settings *ServerSettings) *email {
	return &email{
		settings: settings,
	}
}

func (e *email) SendEmail(to, subject, body string) error {

	auth := smtp.PlainAuth("", e.settings.EmailUsername, e.settings.EmailPassword, e.settings.EmailHost)

	message := fmt.Sprintf("From: %s\r\n", e.settings.EmailFrom)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += fmt.Sprintf("\r\n%s\r\n", body)

	// Sending email.
	err := smtp.SendMail(fmt.Sprintf("%s:%d", e.settings.EmailHost, e.settings.EmailPort), auth, e.settings.EmailFrom, []string{to}, []byte(message))
	if err != nil {
		return err
	}
	return nil
}
