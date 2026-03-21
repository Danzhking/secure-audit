package queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func ConnectRabbitMQ(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	if err != nil {
		zap.L().Fatal("RabbitMQ connection error", zap.Error(err))
	}

	zap.L().Info("RabbitMQ connected")
	return conn
}
