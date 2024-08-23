package email

import (
	"gopkg.in/gomail.v2"
)

const Email = ""

func SendEmailConfirmation(to string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "your-email@example.com")
	// Устанавливаем получателя
	m.SetHeader("To", to)
	// Тема письма
	m.SetHeader("Subject", subject)
	// Тело письма
	m.SetBody("text/plain", body)

	// Настройки SMTP-сервера
	d := gomail.NewDialer("smtp.example.com", 587, "your-email@example.com", "your-password")

	// Отправка письма
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

//to := "recipient@example.com"
//subject := "Подтверждение регистрации"
//body := "Пожалуйста, подтвердите ваш email, перейдя по следующей ссылке: https://example.com/confirm?token=abc123"
