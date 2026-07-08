package main

import (
	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/database"
	ilogger "github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/queue"
	"go.uber.org/zap"
)

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	if err := ilogger.Initialize(); err != nil {
		panic(err)
	}

	cfg := config.Get()
	logger := ilogger.Get()
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("Configuration loaded", zap.String("environment", cfg.App.Env))
	logger.Info("Logger initialized", zap.String("level", cfg.Logger.Level))

	if err := database.Initialize(); err != nil {
		logger.Fatal("failed to initialize PostgreSQL", zap.Error(err))
	}

	if err := database.Health(); err != nil {
		logger.Fatal("database health check failed", zap.Error(err))
	}

	logger.Info("Database health check passed")

	if err := queue.Initialize(); err != nil {
		logger.Fatal("failed to initialize RabbitMQ", zap.Error(err))
	}

	if err := queue.Health(); err != nil {
		logger.Fatal("rabbitmq health check failed", zap.Error(err))
	}

	logger.Info("RabbitMQ health check passed")
	logger.Info("Worker started", zap.String("service", cfg.App.Name))

	run(cfg, logger)
}
