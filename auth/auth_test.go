package auth

import (
	"Backend/Logger"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"testing"
	"time"
)

var engine Auth
var access string

func TestInitializing(t *testing.T) {
	logger, err := Logger.GetLogger(zap.NewDevelopmentConfig())
	if err != nil {
		t.Error(err)
	}

	engine = Auth{
		SingingMethod:  jwt.SigningMethodHS512,
		Issuer:         "Test FP v1.5",
		Audience:       []string{"No one"},
		AccessExpired:  time.Hour * 4,
		RefreshExpired: time.Hour * 24 * 30,
		ChangeExpires:  time.Minute * 15,
		ChatExpires:    time.Hour * 4,
		RefreshLength:  64,

		Log: logger,
	}
	err = engine.GetKeys("./keys/private.pem", "./keys/public.pem")
	if err != nil {
		t.Error(err)
	}
}

func TestGeneratingTokenPair(t *testing.T) {
	token, err := engine.GenerateToken(34, 1)
	if err != nil {
		t.Error(err)
	}

	access = token.Access.Key
	refresh := token.Refresh.Key

	fmt.Println("Access: ", access, "\t Expires: ", token.Access.Expires)
	fmt.Println("Refresh: ", refresh, "\t Expires: ", token.Refresh.Expires)
}

func TestParsingAccess(t *testing.T) {
	acc, err := engine.GetAccessClaims(access)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(acc)
}
