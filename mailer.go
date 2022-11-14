package main

import (
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mailer struct {
	Email    string
	Password string
	Host     string
}

func NewMailer(email, password, host string) *Mailer {
	return &Mailer{
		Email:    email,
		Password: password,
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

	for _, v := range to {
		email.AddTo(v.Email)
	}
	// TODO: ADD Template support

	email.SetBody(mail.TextHTML, body)
	return email.Send(smtpClient)
}
