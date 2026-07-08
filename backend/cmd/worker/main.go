package main

import (
	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	ilogger "github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
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
	logger.Info("Worker starting...", zap.String("service", cfg.App.Name))

	run(cfg, logger)
}
