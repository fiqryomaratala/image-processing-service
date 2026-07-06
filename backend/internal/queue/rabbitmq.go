package queue

import (
	"fmt"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ(cfg config.RabbitMQConfig) (*amqp.Connection, *amqp.Channel, error) {
	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to open rabbitmq channel: %w", err)
	}

	return conn, ch, nil
}