package main

import (
	"os"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/bootstrap"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	ilogger "github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"go.uber.org/zap"
)

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	cfg := config.Get()
	logger := ilogger.MustNew(cfg.Logger, cfg.App)
	defer func() {
		_ = logger.Sync()
	}()

	app, err := bootstrap.NewApp(cfg, logger)
	if err != nil {
		logger.Fatal("failed to bootstrap application", zap.Error(err))
	}
	defer app.Close()

	if err := run(app); err != nil {
		logger.Fatal("failed to run api server", zap.Error(err))
	}

	os.Exit(0)
}
