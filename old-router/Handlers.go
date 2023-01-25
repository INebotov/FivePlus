package router

//
//import (
//	"Backend/Logger"
//	"Backend/auth"
//	"Backend/db"
//	"time"
//)
//
//type Handlers struct {
//	Auth auth.Auth
//	DB   db.DB
//	Log  Logger.Logger
//
//	EmailConfExpired time.Duration
//	PhoneConfExpired time.Duration
//}
//
//type StandartResponse struct {
//	Code    int    `json:"code"`
//	Message string `json:"message"`
//}
//
//func contains[T comparable, A []T](s A, e T) bool {
//	for _, a := range s {
//		if a == e {
//			return true
//		}
//	}
//	return false
//}
