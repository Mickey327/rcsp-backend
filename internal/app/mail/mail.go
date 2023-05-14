package mail

import (
	"log"

	"github.com/Mickey327/rcsp-backend/internal/app/config"
	"gopkg.in/gomail.v2"
)

type Mail struct {
	From    string
	To      string
	Subject string
	Message string
}

func New(from, to, subject, message string) *Mail {
	return &Mail{
		from,
		to,
		subject,
		message,
	}
}

func (m *Mail) SendMail() {
	message := gomail.NewMessage()
	message.SetHeader("From", m.From)
	message.SetHeader("To", m.To)
	message.SetHeader("Subject", m.Subject)
	message.SetBody("text/html", m.Message)

	cfg := config.GetConfig()
	log.Println(cfg.MailHost, cfg.MailPort, cfg.Email, cfg.MailPassword)
	d := gomail.NewDialer(cfg.MailHost, cfg.MailPort, cfg.Email, cfg.MailPassword)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(message); err != nil {
		log.Println(err.Error())
		log.Println("Unable to send message")
	}
}
