package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type ChangeToken struct {
	Issuer    string           `json:"iss,omitempty"`
	Subject   uint             `json:"sub,omitempty"`
	ExpiresAt *jwt.NumericDate `json:"exp,omitempty"`
	jwt.Claims
}

type TockToSend struct {
	Key     string
	Expires time.Time
}

func (a Auth) GenerateChangeToken(userid uint) (TockToSend, error) {
	expires := time.Now().Add(a.ChangeExpires)

	token := jwt.NewWithClaims(a.SingingMethod, ChangeToken{
		Issuer:    a.Issuer,
		Subject:   userid,
		ExpiresAt: jwt.NewNumericDate(expires),
	})

	key, err := token.SignedString(a.PrivateKey)
	if err != nil {
		return TockToSend{}, err
	}

	return TockToSend{
		Key:     key,
		Expires: expires,
	}, nil
}
