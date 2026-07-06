package main

import (
	"log"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	ilogger "github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
)

func main() {
	cfg := config.Load()
	logger := ilogger.New()

	if err := run(cfg, logger); err != nil {
		log.Fatal(err)
	}
}
