package auth

import (
	"Backend/Errors"
	"github.com/golang-jwt/jwt/v4"
	"reflect"
	"time"
)

type ActionToken struct {
	Issuer    string           `json:"iss,omitempty"`
	Subject   uint             `json:"sub,omitempty"`
	ExpiresAt *jwt.NumericDate `json:"exp,omitempty"`
	Actions   string           `json:"act,omitempty"`
	jwt.Claims
}

func (a *Auth) GenerateChangeToken(userid uint, actions string) (Token, error) {
	expires := time.Now().Add(a.ChangeExpires)

	token := jwt.NewWithClaims(a.SingingMethod, ActionToken{
		Issuer:    a.Issuer,
		Subject:   userid,
		Actions:   actions,
		ExpiresAt: jwt.NewNumericDate(expires),
	})

	key, err := token.SignedString(a.PrivateKey)
	if err != nil {
		return Token{}, err
	}

	return Token{
		Key:     key,
		Expires: expires,
	}, nil
}

func (a *Auth) GetActionClaims(token string) (ActionToken, error) {
	null := ActionToken{}
	claims, err := a.getTokenClaims(token)
	if err != nil {
		return null, err
	}
	res := ActionToken{}

	var ok bool
	var val interface{}

	val, ok = claims["iss"]
	if !ok || reflect.TypeOf(val).Kind() != reflect.String {
		return null, Errors.TokenIsIvalidError
	}
	res.Issuer = val.(string)

	val, ok = claims["sub"]
	if !ok || reflect.TypeOf(val).Kind() != reflect.Float64 {
		return null, Errors.TokenIsIvalidError
	}
	res.Subject = uint(val.(float64))

	val, ok = claims["act"]
	if !ok || reflect.TypeOf(val).Kind() != reflect.String {
		return null, Errors.TokenIsIvalidError
	}
	res.Actions = val.(string)

	val, ok = claims["exp"]
	if !ok || reflect.TypeOf(val).Kind() != reflect.Float64 {
		return null, Errors.TokenIsIvalidError
	}
	res.ExpiresAt = jwt.NewNumericDate(time.Unix(int64(val.(float64)), 0))

	return res, nil
}
