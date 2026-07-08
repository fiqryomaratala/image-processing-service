package queue

import (
	"fmt"
	"sync"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var (
	connection *amqp.Connection
	channel    *amqp.Channel
	mu         sync.Mutex
)

func Initialize() error {
	mu.Lock()
	defer mu.Unlock()

	if connection != nil && channel != nil && !connection.IsClosed() && !channel.IsClosed() {
		return nil
	}

	cfg := config.Get()
	log := logger.Get()

	log.Info("Connecting to RabbitMQ...", zap.String("host", cfg.RabbitMQ.Host), zap.String("port", cfg.RabbitMQ.Port))

	conn, err := amqp.Dial(buildURL(cfg.RabbitMQ))
	if err != nil {
		log.Error("Failed to connect RabbitMQ", zap.Error(err))
		return fmt.Errorf("failed to connect rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		log.Error("Failed to connect RabbitMQ", zap.Error(err))
		return fmt.Errorf("failed to open rabbitmq channel: %w", err)
	}

	connection = conn
	channel = ch

	log.Info("RabbitMQ connected successfully")

	return nil
}

func GetConnection() *amqp.Connection {
	if connection == nil {
		panic("rabbitmq is not initialized: call queue.Initialize() before queue.GetConnection()")
	}

	return connection
}

func GetChannel() *amqp.Channel {
	if channel == nil {
		panic("rabbitmq channel is not initialized: call queue.Initialize() before queue.GetChannel()")
	}

	return channel
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if err := closeChannelLocked(); err != nil {
		if closeConnErr := closeConnectionLocked(); closeConnErr != nil {
			return fmt.Errorf("%v; %w", err, closeConnErr)
		}

		return err
	}

	return closeConnectionLocked()
}

func CloseChannel() error {
	mu.Lock()
	defer mu.Unlock()

	return closeChannelLocked()
}

func CloseConnection() error {
	mu.Lock()
	defer mu.Unlock()

	return closeConnectionLocked()
}

func buildURL(cfg config.RabbitMQConfig) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
}

func closeChannelLocked() error {
	if channel == nil {
		return nil
	}

	if channel.IsClosed() {
		channel = nil
		return nil
	}

	if err := channel.Close(); err != nil {
		return fmt.Errorf("failed to close rabbitmq channel: %w", err)
	}

	channel = nil

	return nil
}

func closeConnectionLocked() error {
	if connection == nil {
		return nil
	}

	if connection.IsClosed() {
		connection = nil
		return nil
	}

	if err := connection.Close(); err != nil {
		return fmt.Errorf("failed to close rabbitmq connection: %w", err)
	}

	connection = nil

	return nil
}
