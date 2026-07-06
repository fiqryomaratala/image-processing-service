package main

import (
	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	ilogger "github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()
	logger := ilogger.MustNew(cfg.App)
	defer func() {
		_ = logger.Sync()
	}()

	if err := run(cfg, logger); err != nil {
		logger.Fatal("failed to run api server", zap.Error(err))
	}
}
