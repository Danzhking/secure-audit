package queue

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
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

	q, err := ch.QueueDeclare(
		"security_events",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Error("Queue declare error", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Queue declared", zap.String("queue", q.Name))

	return &Publisher{
		channel: ch,
		queue:   q,
	}, nil
}

func (p *Publisher) Publish(event interface{}) error {
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
		zap.L().Error("RabbitMQ publish error", zap.Error(err))
		return err
	}

	zap.L().Debug("Event published", zap.String("queue", p.queue.Name))
	return nil
}
