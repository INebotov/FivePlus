package main

import (
	"Backend/Closer"
	"Backend/Configurator"
	"Backend/Kafka"
	"Backend/Logger"
	"Backend/Metrics"
	"Backend/auth"
	"Backend/db"
	"Backend/web"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"reflect"
	"runtime"
	"time"
)

func main() {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path = "./config.yaml"
	}

	// Internal: config logger closer data authorize emailLM metrics old-router
	I := GetInternal(path)
	if reflect.DeepEqual(Internal{}, I) {
		runtime.Goexit()
	}

	I.Router.ApplyHandlers()
	I.Router.Start()
}

type Internal struct {
	Config *Configurator.Config[time.Duration]
	Logger *zap.Logger

	Metrics *Metrics.Metrics

	Closer   *Closer.Closer
	DataBase *db.DB
	Auth     *auth.Auth

	Kafka struct {
		Email *Kafka.Kafka
		SMS   *Kafka.Kafka
		Media *Kafka.Kafka
	}

	Router *web.Router
}

func GetInternal(path string) Internal {
	null := Internal{}
	fmt.Println("Configuring ....")
	var res Internal

	config, err := Configurator.Configure(path)
	if err != nil {
		fmt.Printf("Error Configuring: %s\n", err)
		return null
	}
	fmt.Println("Done configuring.\nGetting Logger...")
	res.Config = &config

	logger, err := Logger.GetLogger(config.Logger)
	if err != nil {
		fmt.Printf("Error getting logger: %s", err)
		return null
	}
	res.Logger = logger

	closer := Closer.Closer{
		CloseTimeout: config.Exit.Timeout,
		Logger:       logger,
	}
	logger.Debug("Logger is started successfully!")
	closer.Add("Logger Sync function", logger.Sync)
	res.Closer = &closer

	if config.Kafka.Connect {
		logger.Info("Connecting to email Kafka...")
		emailLM, err := Kafka.GetKafka(Kafka.Config{
			Partition:      config.Kafka.Partition,
			Host:           config.Kafka.Host,
			Port:           config.Kafka.Port,
			UDP:            config.Kafka.UDP,
			WriteTimeOut:   config.Kafka.WriteTimeOut,
			ConnectTimeOut: config.Kafka.ConnectTimeOut,

			Log: logger,
		}, "email")
		if err != nil {
			logger.Fatal("Error connecting to email kafka!", zap.Field{
				Key:    "Error",
				Type:   zapcore.ErrorType,
				String: err.Error(),
			})
		}
		closer.Add("Email Kafka Cloase Function", emailLM.CloserFunc)
		logger.Info("Successfully connected to email kafka!")
		res.Kafka.Email = &emailLM
	}

	logger.Info("Connecting to DataBase....")
	data := db.DB{
		Params: config.DataBase,
		Log:    logger,
	}
	err = data.Connect()
	if err != nil {
		logger.Fatal("Error connecting to database!", zap.Field{
			Key:    "Error",
			Type:   zapcore.ErrorType,
			String: err.Error(),
		})
	}
	closer.Add("Database Exit Func", data.CloserFunc)
	logger.Info("Successfully connected to database!")
	res.DataBase = &data

	logger.Info("Initializing auth Service....")
	authorize := auth.Auth{
		SingingMethod:  jwt.SigningMethodHS512,
		Issuer:         config.App,
		Audience:       config.Auth.Audience,
		AccessExpired:  config.Auth.AccessExpired,
		RefreshExpired: config.Auth.RefreshExpired,
		ChangeExpires:  config.Auth.ChangeExpires,
		ChatExpires:    config.Auth.ChatExpires,
		RefreshLength:  config.Auth.RefreshLength,

		Log: logger,
	}
	err = authorize.GetKeys(config.Auth.Keys.Private, config.Auth.Keys.Public)
	if err != nil {
		logger.Fatal("Error getting keys for auth Service!", zap.Field{
			Key:    "Error",
			Type:   zapcore.ErrorType,
			String: err.Error(),
		})
	}
	res.Auth = &authorize

	metric := Metrics.Metrics{Namespace: config.Namespace}
	metric.Init()
	err = initMetric(&metric)
	if err != nil {
		logger.Fatal("Error initializing metrics!", zap.Field{
			Key:    "Error",
			Type:   zapcore.ErrorType,
			String: err.Error(),
		})
	}
	res.Metrics = &metric

	router := web.Router{
		DataBase:         &data,
		Auth:             &authorize,
		Metrics:          &metric,
		EmailConfExpired: config.Handlers.EmailConfirmationExpired,
		PhoneConfExpired: config.Handlers.PhoneConfExpired,

		Host: config.Handlers.Host,
		Port: config.Handlers.Port,

		Log: logger,
	}
	router.Kafka.Email = res.Kafka.Email
	err = router.Init()
	if err != nil {
		logger.Fatal("Error initializing old-router!", zap.Field{
			Key:    "Error",
			Type:   zapcore.ErrorType,
			String: err.Error(),
		})
	}
	res.Router = &router
	closer.Add("Router stop func", router.CloseFunc)

	go closer.Listen()

	return res
}

func initMetric(m *Metrics.Metrics) error {
	err := m.AddCounter(Metrics.Counter{
		Name:   "requests",
		Help:   "Общее количество запроссов к службе.",
		Labels: nil,
		Action: "any_request",
	})
	if err != nil {
		return err
	}

	err = m.AddCounter(Metrics.Counter{
		Name:   "ok_requests",
		Help:   "Количество успрешных запросов.",
		Labels: nil,
		Action: "success_request",
	})
	if err != nil {
		return err
	}

	err = m.AddCounter(Metrics.Counter{
		Name:   "auth_bad_requests",
		Help:   "Количество неуспешных запросов по причине отсудствия прав доступа.",
		Labels: nil,
		Action: "auth_failure",
	})
	if err != nil {
		return err
	}

	err = m.AddCounter(Metrics.Counter{
		Name:   "bad_request",
		Help:   "Количество неуспешных запросов с неправильным телом запроса.",
		Labels: nil,
		Action: "bad_request_failure",
	})
	if err != nil {
		return err
	}

	err = m.AddCounter(Metrics.Counter{
		Name:   "errors",
		Help:   "Количество внутренних некритичных ошибок.",
		Labels: nil,
		Action: "error",
	})
	if err != nil {
		return err
	}
	return nil
}
