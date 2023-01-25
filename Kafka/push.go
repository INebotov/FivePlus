package Kafka

import (
	"github.com/segmentio/kafka-go"
	"time"
)

func (k *Kafka) Push(messages ...kafka.Message) error {
	err := k.Connection.SetWriteDeadline(time.Now().Add(k.WriteTimeout))
	if err != nil {
		return err
	}
	_, err = k.Connection.WriteMessages(messages...)
	return err
}
