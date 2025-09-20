package email

import (
	"gopkg.in/gomail.v2"
)

const Email = ""

func SendEmailConfirmation(to string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "your-email@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.example.com", 587, "your-email@example.com", "your-password")
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
