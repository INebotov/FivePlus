package Configurator

import (
	"fmt"
	"gopkg.in/yaml.v3"
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
	res := Copy(c)

	if c.DataBase.Password != "" {
		pass, err := os.ReadFile(c.DataBase.Password)
		if err != nil {
			return null, err
		}
		res.DataBase.Password = string(pass)
	}

	var err error
	res.Exit.Timeout, err = getTimeDuration(c.Exit.Timeout)
	if err != nil {
		return null, err
	}

	res.Kafka.ConnectTimeOut, err = getTimeDuration(c.Kafka.ConnectTimeOut)
	if err != nil {
		return null, err
	}
	res.Kafka.WriteTimeOut, err = getTimeDuration(c.Kafka.WriteTimeOut)

	res.Auth.AccessExpired, err = getTimeDuration(c.Auth.AccessExpired)
	if err != nil {
		return null, err
	}
	res.Auth.RefreshExpired, err = getTimeDuration(c.Auth.RefreshExpired)
	if err != nil {
		return null, err
	}
	res.Auth.ChangeExpires, err = getTimeDuration(c.Auth.ChangeExpires)
	if err != nil {
		return null, err
	}
	res.Auth.ChatExpires, err = getTimeDuration(c.Auth.ChatExpires)
	if err != nil {
		return null, err
	}

	res.Handlers.EmailConfirmationExpired, err = getTimeDuration(c.Handlers.EmailConfirmationExpired)
	if err != nil {
		return null, err
	}

	res.Chat.WriteWait, err = getTimeDuration(c.Chat.WriteWait)
	if err != nil {
		return null, err
	}
	res.Chat.PongWait, err = getTimeDuration(c.Chat.PongWait)
	if err != nil {
		return null, err
	}

	res.Handlers.PhoneConfExpired, err = getTimeDuration(c.Handlers.PhoneConfExpired)
	if err != nil {
		return null, err
	}
	res.Handlers.EmailConfirmationExpired, err = getTimeDuration(c.Handlers.EmailConfirmationExpired)
	if err != nil {
		return null, err
	}

	return res, nil
}

func Copy(c Config[string]) Config[time.Duration] {
	return Config[time.Duration]{
		App:       c.App,
		Namespace: c.Namespace,
		DataBase:  c.DataBase,
		Logger:    c.Logger,
		Auth: AuthParams[time.Duration]{
			Audience:      c.Auth.Audience,
			RefreshLength: c.Auth.RefreshLength,
			Keys:          c.Auth.Keys,
		},
		Handlers: HandlersConfig[time.Duration]{
			Port: c.Handlers.Port,
			Host: c.Handlers.Host,
		},
		Chat: ChatParams[time.Duration]{
			MaxMessageSize: c.Chat.MaxMessageSize,
		},
	}
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

func Configure(path string) (Config[time.Duration], error) {
	var FirstConfig Config[string]
	null := Config[time.Duration]{}

	data, err := os.ReadFile(path)
	if err != nil {
		return null, err
	}

	err = yaml.Unmarshal(data, &FirstConfig)
	if err != nil {
		return null, err
	}

	config, err := Parce(FirstConfig)
	if err != nil {
		return null, err
	}
	err = Check(config)
	if err != nil {
		return null, err
	}
	return config, nil
}
