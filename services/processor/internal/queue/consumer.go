package queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Consumer struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewConsumer(conn *amqp.Connection) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := ch.Qos(1, 0, false); err != nil {
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
		return nil, err
	}

	zap.L().Info("Очередь потребителя объявлена", zap.String("queue", q.Name))

	return &Consumer{
		channel: ch,
		queue:   q,
	}, nil
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.Consume(
		c.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	zap.L().Info("Ожидание сообщений", zap.String("queue", c.queue.Name))
	return msgs, nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
}
