package Sender

import (
	"bytes"
	"context"
	"fmt"
	"gopkg.in/gomail.v2"
	"text/template"
)

const (
	RegistrationConfirmationAction = iota + 1
	ProfileSecretsChangeConfirmationAction
	DeteteChildAction
)

type ConfirmationsActions int

func (t ConfirmationsActions) GetType() string {
	return []string{"Registration", "User Profile Secrets Change", "Delete Child"}[t-1]
}

type EmailSender struct {
	Email                           string `yaml:"Email"`
	Password                        string `yaml:"Password"`
	SMTPServer                      string `yaml:"SMTPServer"`
	Port                            int    `yaml:"Port"`
	SSL                             bool   `yaml:"SSL"`
	VerificationChangeTemplate      string
	VerificationEmailTemplate       string
	VerificationDeleteChildTemplate string

	Dealier *gomail.Dialer

	ChangeConfirmTemplate *template.Template
	ConfirmEmailTemplate  *template.Template
	DeleteChildTemplate   *template.Template
}

func (e *EmailSender) Config() error {
	e.Dealier = gomail.NewDialer(e.SMTPServer, e.Port,
		e.Email, e.Password)

	var err error
	e.ChangeConfirmTemplate, err = template.ParseFiles(e.VerificationChangeTemplate)
	if err != nil {
		return err
	}
	e.ConfirmEmailTemplate, err = template.ParseFiles(e.VerificationEmailTemplate)
	if err != nil {
		return err
	}
	e.DeleteChildTemplate, err = template.ParseFiles(e.VerificationDeleteChildTemplate)
	if err != nil {
		return err
	}

	err = e.SendEmail(context.Background(), EmailMessage{
		UserName:  "Тест Тестович",
		Type:      RegistrationConfirmationAction,
		Code:      3566,
		Recipient: e.Email,
	})
	if err != nil {
		return err
	}

	err = e.SendEmail(context.Background(), EmailMessage{
		UserName:  "Теста Тестовна",
		Type:      ProfileSecretsChangeConfirmationAction,
		Code:      5432,
		Recipient: e.Email,
	})
	if err != nil {
		return err
	}

	err = e.SendEmail(context.Background(), EmailMessage{
		UserName:  "Родитель ParentОвич",
		Type:      DeteteChildAction,
		Code:      7675,
		Recipient: e.Email,
	})
	if err != nil {
		return err
	}

	return nil
}

type EmailMessage struct {
	UserName string
	Type     int
	Payload  string

	Code int

	Recipient string
}

func (s *EmailSender) SendEmail(ctx context.Context, mess EmailMessage) error {
	err := make(chan error)
	go s.acyncSend(mess, err)

	select {
	case <-ctx.Done():
		return fmt.Errorf("send canseled")

	case inerr := <-err:
		return inerr
	}
}

func (s *EmailSender) acyncSend(mess EmailMessage, send chan error) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.Email)
	msg.SetHeader("To", mess.Recipient)
	if mess.Type == RegistrationConfirmationAction {
		msg.SetHeader("Subject", ConfirmationsActions(mess.Type).GetType())
		var tpl bytes.Buffer
		if err := s.ChangeConfirmTemplate.Execute(&tpl, mess); err != nil {
			send <- err
		}

		msg.SetBody("text/html", tpl.String())
	} else if mess.Type == ProfileSecretsChangeConfirmationAction {
		msg.SetHeader("Subject", ConfirmationsActions(mess.Type).GetType())
		var tpl bytes.Buffer
		if err := s.ConfirmEmailTemplate.Execute(&tpl, mess); err != nil {
			send <- err
		}

		msg.SetBody("text/html", tpl.String())
	} else if mess.Type == DeteteChildAction {
		msg.SetHeader("Subject", ConfirmationsActions(mess.Type).GetType())
		var tpl bytes.Buffer
		if err := s.DeleteChildTemplate.Execute(&tpl, mess); err != nil {
			send <- err
		}

		msg.SetBody("text/html", tpl.String())
	} else {
		send <- fmt.Errorf("wrong message type")
	}

	if err := s.Dealier.DialAndSend(msg); err != nil {
		send <- err
	}

	send <- nil
}
