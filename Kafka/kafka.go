package Kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
)

var null = Kafka{}

type Kafka struct {
	Config     Config
	Connection *kafka.Conn

	WriteTimeout time.Duration

	Log *zap.Logger

	CloserFunc func() error
}
type Config struct {
	Partition int

	Host           string
	Port           int
	UDP            bool
	WriteTimeOut   time.Duration
	ConnectTimeOut time.Duration

	Log *zap.Logger
}

func GetKafka(config Config, topic string) (Kafka, error) {
	ctx, _ := context.WithTimeout(context.Background(), config.ConnectTimeOut)
	protocol := "tcp"
	if config.UDP {
		protocol = "udp"
	}
	conn, err := kafka.DialLeader(ctx, protocol, fmt.Sprintf("%s:%d", config.Host, config.Port), topic, config.Partition)
	if err != nil {
		return null, err
	}

	var res Kafka
	res.Connection = conn
	res.Config = config
	res.CloserFunc = conn.Close
	res.WriteTimeout = config.WriteTimeOut
	res.Log = config.Log
	return res, nil
}
