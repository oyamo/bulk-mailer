package main

import (
	"errors"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mailer struct {
	Email    string
	Password string
	Name     string
	Host     string
}

func NewMailer(email, password, host, name string) *Mailer {
	return &Mailer{
		Email:    email,
		Password: password,
		Name:     name,
		Host:     host,
	}
}

func (m *Mailer) SendMail(to []Recipient, subject, body string) error {
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = 587
	server.Username = m.Email
	server.Password = m.Password
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(m.Email)
	email.SetSubject(subject)

	if len(to) == 0 {
		return errors.New("no recipients")
	}

	email.AddTo(to[0].String())

	if len(to) > 1 {
		for _, recipient := range to[1:] {
			email.AddBcc(recipient.String())
		}
	}
	// TODO: ADD Template support

	email.SetBody(mail.TextHTML, body)
	return email.Send(smtpClient)
}
