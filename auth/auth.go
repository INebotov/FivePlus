package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

type Auth struct {
	PrivateKey []byte
	PublicKey  []byte

	SingingMethod jwt.SigningMethod

	Issuer string

	Audience []string

	AccessExpired  time.Duration
	RefreshExpired time.Duration
	ChangeExpires  time.Duration
	ChatExpires    time.Duration
	RefreshLength  uint

	Log *zap.Logger
}

func (a *Auth) GetKeys(privKeyPath, pubKeyPath string) (err error) {
	a.PrivateKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		return err
	}
	a.PublicKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return err
	}
	return nil
}
