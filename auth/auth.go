package auth

import (
	"BackendSimple/db"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"io/ioutil"
	"math/rand"
	"time"
)

type Auth struct {
	PrivateKey []byte
	PublicKey  []byte

	SingingMethod jwt.SigningMethod

	Issuer   string
	Audience []string

	AccessExpired  time.Duration
	RefreshExpired time.Duration
	ChangeExpires  time.Duration
	ChatExpires    time.Duration
	RefreshLength  uint

	DB db.DB
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
	jwt.Claims
}

func (a Auth) NewTockenID() string {
	return uuid.New().String() // TODO
}

func (a Auth) GetKeys(privKeyPath, pubKeyPath string) (err error) {
	a.PrivateKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		return err
	}
	a.PublicKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return err
	}
	//a.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	//if err != nil {
	//	return err
	//}
	//a.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKey)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (a Auth) GenerateRefresh() string {
	rand.Seed(time.Now().UnixNano())

	allowed := []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")
	allowedLength := len(allowed)

	res := make([]rune, a.RefreshLength)
	for i := uint(0); i < a.RefreshLength; i++ {
		res[i] = allowed[rand.Intn(allowedLength)]
	}
	return string(res)
}
