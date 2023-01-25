package auth

import (
	"Backend/Errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"math/rand"
	"reflect"
	"time"
)

type Token struct {
	Key     string
	Expires time.Time
}

type TokenPair struct {
	Access  Token `json:"access"`
	Refresh Token `json:"refresh"`
}

type AccessToken struct {
	Issuer        string           `json:"iss,omitempty"`
	Subject       uint             `json:"sub,omitempty"`
	Audience      jwt.ClaimStrings `json:"aud,omitempty"`
	ExpiresAt     *jwt.NumericDate `json:"exp,omitempty"`
	NotBefore     *jwt.NumericDate `json:"nbf,omitempty"`
	IssuedAt      *jwt.NumericDate `json:"iat,omitempty"`
	ID            string           `json:"jti,omitempty"`
	LevelOfAccess int              `json:"alv,omitempty"`
	jwt.Claims                     // `json:"-"`
}

func (a *Auth) GenerateToken(userid uint, userlevel int) (TokenPair, error) {
	expires := time.Now().Add(a.AccessExpired)
	id := uuid.New().String()

	acc := AccessToken{
		Issuer:        a.Issuer,
		Subject:       userid,
		Audience:      a.Audience,
		ExpiresAt:     jwt.NewNumericDate(expires),
		NotBefore:     jwt.NewNumericDate(time.Now()),
		IssuedAt:      jwt.NewNumericDate(time.Now()),
		ID:            id,
		LevelOfAccess: userlevel,
	}

	token := jwt.NewWithClaims(a.SingingMethod, acc)

	rand.Seed(time.Now().UnixNano())

	allowed := []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")
	allowedLength := len(allowed)

	refresh := make([]rune, a.RefreshLength)
	for i := uint(0); i < a.RefreshLength; i++ {
		refresh[i] = allowed[rand.Intn(allowedLength)]
	}
	refreshExp := time.Now().Add(a.RefreshExpired)

	key, err := token.SignedString(a.PrivateKey)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		Access:  Token{key, expires},
		Refresh: Token{string(refresh), refreshExp},
	}, nil
}
func (a *Auth) GetAccessClaims(token string) (AccessToken, error) {
	null := AccessToken{}
	claims, err := a.getTokenClaims(token)
	if err != nil {
		return null, err
	}
	res := AccessToken{}

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

	val, ok = claims["aud"]
	if !ok || reflect.TypeOf(val) != reflect.TypeOf([]interface{}{""}) { // Danger panic
		return null, Errors.TokenIsIvalidError
	}
	tmp := val.([]interface{})

	res.Audience = make([]string, len(tmp))
	for i := 0; i < len(tmp); i++ {
		res.Audience[i] = tmp[i].(string)
	}

	val, ok = claims["exp"]
	if !ok || reflect.TypeOf(val).Kind() != reflect.Float64 {
		return null, Errors.TokenIsIvalidError
	}
	res.ExpiresAt = jwt.NewNumericDate(time.Unix(int64(val.(float64)), 0))

	val, ok = claims["nbf"]
	if !ok || reflect.TypeOf(val).Kind() != reflect.Float64 {
		return null, Errors.TokenIsIvalidError
	}
	res.NotBefore = jwt.NewNumericDate(time.Unix(int64(val.(float64)), 0))

	val, ok = claims["iat"]
	if !ok || reflect.TypeOf(val).Kind() != reflect.Float64 {
		return null, Errors.TokenIsIvalidError
	}

	res.IssuedAt = jwt.NewNumericDate(time.Unix(int64(val.(float64)), 0))
	val, ok = claims["jti"]
	if !ok || reflect.TypeOf(val).Kind() != reflect.String {
		return null, Errors.TokenIsIvalidError
	}
	res.ID = val.(string)

	val, ok = claims["alv"]
	if !ok || reflect.TypeOf(val).Kind() != reflect.Float64 {
		return null, Errors.TokenIsIvalidError
	}
	res.LevelOfAccess = int(val.(float64))

	return res, nil
}

func (a *Auth) getTokenClaims(token string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if token.Method != a.SingingMethod {
			return nil, Errors.BadSigningMethodError
		}
		return a.PrivateKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims, nil
	} else {
		return nil, Errors.TokenIsIvalidError
	}
}
