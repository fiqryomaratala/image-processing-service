package app

import (
	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/server"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Config          *config.Config
	Logger          *zap.Logger
	Database        *gorm.DB
	RabbitMQConn    *amqp.Connection
	RabbitMQChannel *amqp.Channel
	Storage         *minio.Client
	HTTPServer      *server.Server
}
