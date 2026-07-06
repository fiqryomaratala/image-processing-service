package main

import (
	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	ilogger "github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
)

func main() {
	cfg := config.Load()
	logger := ilogger.New()

	run(cfg, logger)
}
