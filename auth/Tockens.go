package auth

import (
	"BackendSimple/db"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type TokenPair struct {
	Access        string    `json:"access"`
	AccessExpires time.Time `json:"access_expires"`

	Refresh        string    `json:"refresh"`
	RefreshExpires time.Time `json:"refresh_expires"`
}

func (a Auth) GenerateToken(userid uint, userlevel int, client string) (TokenPair, error) {
	expires := time.Now().Add(a.AccessExpired)
	id := a.NewTockenID()

	token := jwt.NewWithClaims(a.SingingMethod, AccessToken{
		Issuer:        a.Issuer,
		Subject:       userid,
		Audience:      a.Audience,
		ExpiresAt:     jwt.NewNumericDate(expires),
		NotBefore:     jwt.NewNumericDate(time.Now()),
		IssuedAt:      jwt.NewNumericDate(time.Now()),
		ID:            id,
		LevelOfAccess: userlevel,
	})

	refresh := a.GenerateRefresh()
	refreshExp := time.Now().Add(a.RefreshExpired)

	key, err := token.SignedString(a.PrivateKey)
	if err != nil {
		return TokenPair{}, err
	}

	tock := db.RefreshToken{
		Refresh:   refresh,
		ExpiresIn: refreshExp.Unix(),
		UserID:    userid,
		Client:    client,
	}

	err = a.DB.SetToken(&tock)

	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		Access:        key,
		AccessExpires: expires,

		Refresh:        refresh,
		RefreshExpires: refreshExp,
	}, nil
}

func (a Auth) GetTokenClaims(token string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if token.Method != a.SingingMethod {
			return nil, BadSigningMethodError
		}
		return a.PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims, nil
	} else {
		return nil, TokenIsIvalidError
	}
}
