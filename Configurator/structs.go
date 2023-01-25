package Configurator

import (
	"Backend/db"
	"go.uber.org/zap"
	"time"
)

type StringOrTime interface {
	string | time.Duration
}
type HandlersConfig[T StringOrTime] struct {
	EmailConfirmationExpired T
	Port                     int
	Host                     string
}
type Exit[E StringOrTime] struct {
	Timeout E
	WhaitWS bool
}
type KafkaConfig[E StringOrTime] struct {
	Topic     string
	Partition int

	Host           string
	Port           int
	UDP            bool
	WriteTimeOut   E
	ConnectTimeOut E
}
type Config[E StringOrTime] struct {
	App string

	Logger zap.Config
	Exit   Exit[E]
	Kafka  KafkaConfig[E]

	DataBase db.DBParams
	Auth     AuthParams[E]
	Handlers HandlersConfig[E]
	Chat     ChatParams[E]
}
type ChatParams[D StringOrTime] struct {
	WriteWait      D     // Max wait time when writing message to peer
	PongWait       D     //Max time till next pong from peer
	MaxMessageSize int64 // Maximum message size allowed from peer.
}
type AuthParams[D StringOrTime] struct {
	Audience       []string
	AccessExpired  D
	RefreshExpired D
	ChangeExpires  D
	ChatExpires    D
	RefreshLength  uint
	Keys           AuthKeys
}
type AuthKeys struct {
	Private string
	Public  string
}
