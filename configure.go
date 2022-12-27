package main

import (
	"BackendSimple/db"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func getTime(multiplayer time.Duration, amount string) (time.Duration, error) {
	timeInt, err := strconv.Atoi(amount)
	if err != nil {
		return 0 * time.Millisecond, err
	}
	return multiplayer * time.Duration(timeInt), nil
}

func getTimeDuration(t string) (time.Duration, error) {
	var res time.Duration
	var err error

	const null = 0 * time.Millisecond

	switch true {
	case strings.Contains(t, "d"):
		if res, err = getTime(time.Hour*24, t[:len(t)-1]); err != nil {
			return null, err
		}
	case strings.Contains(t, "h"):
		if res, err = getTime(time.Hour, t[:len(t)-1]); err != nil {
			return null, err
		}
	case strings.Contains(t, "m"):
		if res, err = getTime(time.Minute, t[:len(t)-1]); err != nil {
			return null, err
		}
	case strings.Contains(t, "s"):
		if res, err = getTime(time.Second, t[:len(t)-1]); err != nil {
			return null, err
		}
	case strings.Contains(t, "ms"):
		if res, err = getTime(time.Millisecond, t[:len(t)-2]); err != nil {
			return null, err
		}
	case strings.Contains(t, "ms"):
		if res, err = getTime(time.Millisecond, t[:len(t)-2]); err != nil {
			return null, err
		}
	default:
		return null, fmt.Errorf("wrong time format")
	}

	return res, nil

}

func Parce(c Config[string]) (Config[time.Duration], error) {
	var null = Config[time.Duration]{}
	pass, err := os.ReadFile(c.DataBase.Password)
	if err != nil {
		return null, err
	}
	EmailPass, err := os.ReadFile(c.Sender.Email.Password)
	if err != nil {
		return null, err
	}

	AccessExpired, err := getTimeDuration(c.Auth.AccessExpired)
	if err != nil {
		return null, err
	}
	RefreshExpired, err := getTimeDuration(c.Auth.RefreshExpired)
	if err != nil {
		return null, err
	}
	ChangeExpires, err := getTimeDuration(c.Auth.ChangeExpires)
	if err != nil {
		return null, err
	}
	ChatExpires, err := getTimeDuration(c.Auth.ChatExpires)
	if err != nil {
		return null, err
	}

	EmailConfirmationExpired, err := getTimeDuration(c.Handlers.EmailConfirmationExpired)
	if err != nil {
		return null, err
	}

	WriteWait, err := getTimeDuration(c.Chat.WriteWait)
	if err != nil {
		return null, err
	}
	PongWait, err := getTimeDuration(c.Chat.PongWait)
	if err != nil {
		return null, err
	}

	res := Config[time.Duration]{
		App: c.App,
		DataBase: db.DBParams{
			Name:     c.DataBase.Name,
			Type:     c.DataBase.Type,
			Host:     c.DataBase.Host,
			User:     c.DataBase.User,
			Password: string(pass),
			TimeZone: c.DataBase.TimeZone,
			Port:     c.DataBase.Port,
			SslMode:  c.DataBase.SslMode,
		},
		Auth: AuthParams[time.Duration]{
			Audience:       c.Auth.Audience,
			AccessExpired:  AccessExpired,
			RefreshExpired: RefreshExpired,
			ChangeExpires:  ChangeExpires,
			ChatExpires:    ChatExpires,
			RefreshLength:  c.Auth.RefreshLength,
			Keys:           c.Auth.Keys,
		},
		Handlers: struct {
			EmailConfirmationExpired time.Duration
		}{
			EmailConfirmationExpired: EmailConfirmationExpired},
		Sender: struct {
			Email EmailSenderParams
		}{
			Email: EmailSenderParams{
				SenderType:                 c.Sender.Email.SenderType,
				Email:                      c.Sender.Email.Email,
				Password:                   string(EmailPass),
				SMTPServer:                 c.Sender.Email.SMTPServer,
				Port:                       c.Sender.Email.Port,
				SSL:                        c.Sender.Email.SSL,
				VerificationEmailTemplate:  c.Sender.Email.VerificationEmailTemplate,
				VerificationChangeTemplate: c.Sender.Email.VerificationChangeTemplate,
				DeleteChildTemplate:        c.Sender.Email.DeleteChildTemplate,
			},
		},
		Chat: ChatParams[time.Duration]{
			WriteWait:      WriteWait,
			PongWait:       PongWait,
			MaxMessageSize: c.Chat.MaxMessageSize,
		},
	}
	return res, nil
}

func Check(c Config[time.Duration]) error {
	d := c.DataBase
	a := c.Auth
	if d.Name == "" ||
		(d.Type == 2 && (d.Host == "" || d.User == "" || d.Password == "" ||
			d.Port == 0)) || a.Keys.Private == "" || a.Keys.Public == "" { // TODO more
		return fmt.Errorf("wrong config format")
	}
	return nil
}

func Configure(path string) Config[time.Duration] {
	var FirstConfig Config[string]

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error reading configuration: %v", err)
	}

	err = yaml.Unmarshal(data, &FirstConfig)
	if err != nil {
		log.Fatalf("error parcing configuration: %v", err)
	}

	config, err := Parce(FirstConfig)
	if err != nil {
		log.Fatalf("error parcing to real conf: %v", err)
	}
	err = Check(config)
	if err != nil {
		log.Fatalf("error parcing to real conf: %v", err)
	}
	return config
}

type StringOrTime interface {
	string | time.Duration
}

type Config[E StringOrTime] struct {
	App string

	DataBase db.DBParams

	Auth AuthParams[E]

	Handlers struct {
		EmailConfirmationExpired E
	}

	Sender struct {
		Email EmailSenderParams
	}

	Chat ChatParams[E]
}
type EmailSenderParams struct {
	SenderType                 int
	Email                      string
	Password                   string
	SMTPServer                 string
	Port                       int
	SSL                        bool
	VerificationChangeTemplate string
	VerificationEmailTemplate  string
	DeleteChildTemplate        string
}
type ChatParams[D StringOrTime] struct {
	WriteWait      D     // Max wait time when writing message to peer
	PongWait       D     //Max time till next pong from peer
	MaxMessageSize int64 // Maximum message size allowed from peer.
}
type AuthParams[D StringOrTime] struct {
	Audience       []string
	AccessExpired  D
	RefreshExpired D
	ChangeExpires  D
	ChatExpires    D
	RefreshLength  uint
	Keys           AuthKeys
}
type AuthKeys struct {
	Private string
	Public  string
}
