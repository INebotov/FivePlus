package Configurator

import (
	"Backend/db"
	"go.uber.org/zap"
	"time"
)

type Config[E StringOrTime] struct {
	App       string `yaml:"App"`
	Namespace string `yaml:"Namespace"`

	Logger zap.Config     `yaml:"Logger"`
	Exit   Exit[E]        `yaml:"Exit"`
	Kafka  KafkaConfig[E] `yaml:"Kafka"`

	DataBase db.DBParams       `yaml:"DataBase"`
	Auth     AuthParams[E]     `yaml:"Auth"`
	Handlers HandlersConfig[E] `yaml:"Handlers"`
	Chat     ChatParams[E]     `yaml:"Chat"`
}

type StringOrTime interface {
	string | time.Duration
}
type HandlersConfig[T StringOrTime] struct {
	EmailConfirmationExpired T      `yaml:"EmailConfirmationExpired"`
	PhoneConfExpired         T      `yaml:"PhoneConfExpired"`
	Port                     int    `yaml:"Port"`
	Host                     string `yaml:"Host"`
}
type Exit[E StringOrTime] struct {
	Timeout E    `yaml:"Timeout"`
	WhaitWS bool `yaml:"WhaitWS"`
}
type KafkaConfig[E StringOrTime] struct {
	Connect   bool `yaml:"Connect"`
	Partition int  `yaml:"Partition"`

	Host           string `yaml:"Host"`
	Port           int    `yaml:"Port"`
	UDP            bool   `yaml:"UDP"`
	WriteTimeOut   E      `yaml:"WriteTimeOut"`
	ConnectTimeOut E      `yaml:"ConnectTimeOut"`
}
type ChatParams[D StringOrTime] struct {
	WriteWait      D     `yaml:"WriteWait"`      // Max wait time when writing message to peer
	PongWait       D     `yaml:"PongWait"`       //Max time till next pong from peer
	MaxMessageSize int64 `yaml:"MaxMessageSize"` // Maximum message size allowed from peer.
}
type AuthParams[D StringOrTime] struct {
	Audience       []string `yaml:"Audience"`
	AccessExpired  D        `yaml:"AccessExpired"`
	RefreshExpired D        `yaml:"RefreshExpired"`
	ChangeExpires  D        `yaml:"ChangeExpires"`
	ChatExpires    D        `yaml:"ChatExpires"`
	RefreshLength  uint     `yaml:"RefreshLength"`
	Keys           AuthKeys `yaml:"Keys"`
}
type AuthKeys struct {
	Private string `yaml:"Private"`
	Public  string `yaml:"Public"`
}
