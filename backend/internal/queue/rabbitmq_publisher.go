package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	channel   *amqp.Channel
	queueName string
}

var _ Publisher = (*RabbitMQPublisher)(nil)

func NewRabbitMQPublisher(channel *amqp.Channel) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		channel:   channel,
		queueName: ImageProcessingQueue,
	}
}

func (p *RabbitMQPublisher) PublishImageJob(ctx context.Context, job ImageJob) error {
	if p.channel == nil {
		return QueueUnavailable("rabbitmq channel is not initialized")
	}

	body, err := json.Marshal(job)
	if err != nil {
		return PublishFailed(fmt.Sprintf("failed to marshal image job: %v", err))
	}

	queue, err := p.channel.QueueDeclare(
		p.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return QueueUnavailable(fmt.Sprintf("failed to declare queue %s: %v", p.queueName, err))
	}

	if err := p.channel.PublishWithContext(ctx, "", queue.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now().UTC(),
		Body:         body,
	}); err != nil {
		return PublishFailed(fmt.Sprintf("failed to publish job to queue %s: %v", p.queueName, err))
	}

	return nil
}
