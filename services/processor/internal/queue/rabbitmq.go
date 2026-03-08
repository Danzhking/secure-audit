package queue

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ(url string) *amqp.Connection {
	var conn *amqp.Connection
	var err error

	for i := range 10 {
		conn, err = amqp.Dial(url)
		if err == nil {
			log.Println("RabbitMQ connected")
			return conn
		}
		log.Printf("RabbitMQ not ready (attempt %d/10): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("RabbitMQ connection failed after 10 attempts:", err)
	return nil
}
