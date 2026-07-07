package main

import (
	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	ilogger "github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
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

	run(cfg, logger)
}
