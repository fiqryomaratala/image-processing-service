package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Consume(queueName string) (<-chan amqp.Delivery, error) {
	ch := GetChannel()

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	deliveries, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume queue %s: %w", queueName, err)
	}

	return deliveries, nil
}
