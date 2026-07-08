package bootstrap

import (
	"fmt"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/database"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/queue"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/storage"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Config       *config.Config
	Logger       *zap.Logger
	Postgres     *gorm.DB
	RabbitMQConn *amqp.Connection
	RabbitMQChan *amqp.Channel
	MinIO        *minio.Client
}

func NewApp(cfg *config.Config, logger *zap.Logger) (*App, error) {
	db := database.Get()
	rabbitConn := queue.GetConnection()
	rabbitChan := queue.GetChannel()

	minioClient, err := storage.NewMinIO(cfg.MinIO)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio: %w", err)
	}

	app := &App{
		Config:       cfg,
		Logger:       logger,
		Postgres:     db,
		RabbitMQConn: rabbitConn,
		RabbitMQChan: rabbitChan,
		MinIO:        minioClient,
	}

	return app, nil
}

func (a *App) Close() {
	if a.RabbitMQChan != nil {
		_ = queue.Close()
	}

	if a.Postgres != nil {
		_ = database.Close()
	}
}
