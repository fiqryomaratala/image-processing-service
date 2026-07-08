package app

import (
	"fmt"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/database"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/handler"
	ilogger "github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/queue"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/server"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/storage"
	"go.uber.org/zap"
)

func BootstrapAPI() (*App, error) {
	return bootstrap(true)
}

func BootstrapWorker() (*App, error) {
	return bootstrap(false)
}

func bootstrap(withHTTPServer bool) (*App, error) {
	if err := config.Load(); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	if err := ilogger.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	cfg := config.Get()
	log := ilogger.Get()

	log.Info("Bootstrapping application...")
	log.Info("Configuration loaded", zap.String("environment", cfg.App.Env))
	log.Info("Logger initialized", zap.String("level", cfg.Logger.Level))

	if err := database.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize PostgreSQL: %w", err)
	}
	if err := database.Health(); err != nil {
		return nil, fmt.Errorf("database health check failed: %w", err)
	}
	log.Info("PostgreSQL connected")

	if err := queue.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize RabbitMQ: %w", err)
	}
	if err := queue.Health(); err != nil {
		return nil, fmt.Errorf("rabbitmq health check failed: %w", err)
	}
	log.Info("RabbitMQ connected")

	if err := storage.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO: %w", err)
	}
	if err := storage.EnsureBucket(); err != nil {
		return nil, fmt.Errorf("failed to ensure MinIO bucket: %w", err)
	}
	if err := storage.Health(); err != nil {
		return nil, fmt.Errorf("storage health check failed: %w", err)
	}
	log.Info("MinIO connected", zap.String("bucket", cfg.MinIO.BucketName))

	var httpServer *server.Server
	if withHTTPServer {
		healthHandler := handler.NewHealthHandler()
		httpServer = server.New(cfg.App, cfg.CORS, log, healthHandler)
		log.Info("HTTP Server initialized", zap.String("address", cfg.App.Address()))
	}

	container := newContainer(
		cfg,
		log,
		database.Get(),
		queue.GetConnection(),
		queue.GetChannel(),
		storage.GetClient(),
		httpServer,
	)

	log.Info("Application bootstrapped successfully")

	return container, nil
}
