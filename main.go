package main

import (
	"BackendSimple/Chat"
	"BackendSimple/auth"
	"BackendSimple/db"
	"BackendSimple/router"
	"flag"
	"github.com/gofiber/adaptor/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"time"
)

func main() {
	database, err := db.GetDB(db.DBParams{
		Name: "DB.sqlite",
		Type: db.TypeSQLLite,
	})
	if err != nil {
		panic(err)
	}

	autntification := auth.Auth{
		SingingMethod:  jwt.SigningMethodHS512,
		Issuer:         "FP Backend v0.1",
		Audience:       []string{"FP Backend v0.*"},
		AccessExpired:  time.Hour,
		RefreshExpired: time.Hour * 24 * 30,
		DB:             database,
		ChangeExpires:  time.Minute * 15,
		ChatExpires:    time.Hour * 4,
		RefreshLength:  64,
	}
	err = autntification.GetKeys("./private.pem", "./public.pem")
	if err != nil {
		panic(err)
	}
	rout, err := router.GetRouter()
	if err != nil {
		panic(err)
	}

	chat := Chat.NewChat(database, autntification)

	handlers := router.Handlers{
		Auth:             autntification,
		DB:               database,
		EmailConfExpired: time.Minute * 30,
		PhoneConfExpired: time.Minute * 20,
		Chat:             &chat,
	}
	middleware := router.Middleware{
		Auth: autntification,
		DB:   database,
	}
	err = rout.ApplyHandlers(handlers, middleware)
	if err != nil {
		panic(err)
	}

	flag.Parse()

	http.HandleFunc("/chat", chat.ServeWs)
	http.Handle("/", adaptor.FiberApp(rout.App))

	log.Println("Listening on :9090!")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
