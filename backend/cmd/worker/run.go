package main

import (
	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"go.uber.org/zap"
)

func run(cfg *config.Config, logger *zap.Logger) {
	logger.Info("image worker started", zap.String("environment", cfg.App.Env))
}
