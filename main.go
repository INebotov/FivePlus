package main

import (
	"BackendSimple/Chat"
	"BackendSimple/Sender"
	"BackendSimple/auth"
	"BackendSimple/db"
	"BackendSimple/router"
	"flag"
	"github.com/gofiber/adaptor/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"os"
)

func main() {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path = "./config.yaml"
	}

	config := Configure(path)

	sender := Sender.EmailSender{
		Email:                           config.Sender.Email.Email,
		Password:                        config.Sender.Email.Password,
		SMTPServer:                      config.Sender.Email.SMTPServer,
		Port:                            config.Sender.Email.Port,
		SSL:                             config.Sender.Email.SSL,
		VerificationChangeTemplate:      config.Sender.Email.VerificationChangeTemplate,
		VerificationEmailTemplate:       config.Sender.Email.VerificationEmailTemplate,
		VerificationDeleteChildTemplate: config.Sender.Email.DeleteChildTemplate,
	}

	err := sender.Config()
	if err != nil {
		panic(err)
	}

	database, err := db.GetDB(config.DataBase, sender)
	if err != nil {
		panic(err)
	}

	autntification := auth.Auth{
		SingingMethod:  jwt.SigningMethodHS512,
		Issuer:         config.App,
		Audience:       config.Auth.Audience,
		AccessExpired:  config.Auth.AccessExpired,
		RefreshExpired: config.Auth.RefreshExpired,
		DB:             database,
		ChangeExpires:  config.Auth.ChangeExpires,
		ChatExpires:    config.Auth.ChatExpires,
		RefreshLength:  config.Auth.RefreshLength,
	}
	err = autntification.GetKeys(config.Auth.Keys.Private, config.Auth.Keys.Public)
	if err != nil {
		panic(err)
	}
	rout, err := router.GetRouter()
	if err != nil {
		panic(err)
	}

	chat := Chat.NewChat(database, autntification, Chat.ClientParams{
		WriteWait:      config.Chat.WriteWait,
		PongWait:       config.Chat.PongWait,
		MaxMessageSize: config.Chat.MaxMessageSize,
	})

	handlers := router.Handlers{
		Auth:             autntification,
		DB:               database,
		EmailConfExpired: config.Handlers.EmailConfirmationExpired,
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
