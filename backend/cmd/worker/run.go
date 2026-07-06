package main

import (
	"log"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
)

func run(cfg config.Config, logger *log.Logger) {
	logger.Printf("Image Worker Started | env=%s", cfg.App.Env)
}
