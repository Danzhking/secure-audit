package queue

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewPublisher(conn *amqp.Connection) (*Publisher, error) {

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	log.Println("Creating queue security_events")

	q, err := ch.QueueDeclare(
		"security_events",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Println("Queue declare error:", err)
		return nil, err
	}

	log.Println("Queue created:", q.Name)

	return &Publisher{
		channel: ch,
		queue:   q,
	}, nil
}

func (p *Publisher) Publish(event interface{}) error {

	log.Println("Publishing event...")

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println("RabbitMQ publish error:", err)
		return err
	}

	log.Println("Event successfully published")

	return nil
}
