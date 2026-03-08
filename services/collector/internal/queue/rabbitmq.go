package queue

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ(url string) *amqp.Connection {

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("RabbitMQ connection error:", err)
	}

	log.Println("RabbitMQ connected")

	return conn
}
