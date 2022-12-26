package router

import (
	"BackendSimple/Chat"
	"BackendSimple/auth"
	"BackendSimple/db"
	"time"
)

type Handlers struct {
	Auth             auth.Auth
	DB               db.DB
	EmailConfExpired time.Duration
	PhoneConfExpired time.Duration
	Chat             *Chat.Chat
}

type StandartResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func contains[T comparable, A []T](s A, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
