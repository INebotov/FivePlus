package auth

//
//import (
//	"github.com/golang-jwt/jwt/v4"
//	"time"
//)
//
//type ChatToken struct {
//	Issuer    string           `json:"iss,omitempty"`
//	Subject   uint             `json:"sub,omitempty"`
//	Audience  jwt.ClaimStrings `json:"aud,omitempty"`
//	ExpiresAt *jwt.NumericDate `json:"exp,omitempty"`
//	NotBefore *jwt.NumericDate `json:"nbf,omitempty"`
//	IssuedAt  *jwt.NumericDate `json:"iat,omitempty"`
//	ChatID    string           `json:"cid,omitempty"`
//	jwt.Claims
//}
//
//func (a Auth) GenChatToken(userid uint, chatid string) (string, error) {
//	expires := time.Now().Add(a.AccessExpired)
//
//	token := jwt.NewWithClaims(a.SingingMethod, ChatToken{
//		Issuer:    a.Issuer,
//		Subject:   userid,
//		Audience:  a.Audience,
//		ExpiresAt: jwt.NewNumericDate(expires),
//		NotBefore: jwt.NewNumericDate(time.Now()),
//		IssuedAt:  jwt.NewNumericDate(time.Now()),
//		ChatID:    chatid,
//	})
//
//	key, err := token.SignedString(a.PrivateKey)
//	if err != nil {
//		return "", err
//	}
//
//	return key, nil
//}
