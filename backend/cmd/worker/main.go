package main

import (
	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	ilogger "github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	logger := ilogger.MustNew(cfg.App)
	defer func() {
		_ = logger.Sync()
	}()

	run(cfg, logger)
}
