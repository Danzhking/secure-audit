package queue

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func ConnectRabbitMQ(url string) *amqp.Connection {
	var conn *amqp.Connection
	var err error

	for i := range 15 {
		conn, err = amqp.Dial(url)
		if err == nil {
			zap.L().Info("RabbitMQ connected")
			return conn
		}
		zap.L().Warn("RabbitMQ not ready",
			zap.Int("attempt", i+1),
			zap.Error(err),
		)
		time.Sleep(3 * time.Second)
	}

	zap.L().Fatal("RabbitMQ connection failed after 15 attempts", zap.Error(err))
	return nil
}
