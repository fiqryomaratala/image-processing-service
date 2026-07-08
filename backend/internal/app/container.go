package app

import (
	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/server"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func newContainer(
	cfg *config.Config,
	log *zap.Logger,
	db *gorm.DB,
	rabbitConn *amqp.Connection,
	rabbitCh *amqp.Channel,
	storageClient *minio.Client,
	httpServer *server.Server,
) *App {
	return &App{
		Config:          cfg,
		Logger:          log,
		Database:        db,
		RabbitMQConn:    rabbitConn,
		RabbitMQChannel: rabbitCh,
		Storage:         storageClient,
		HTTPServer:      httpServer,
	}
}
