package main

import (
	iapp "github.com/fiqryomaratala/image-processing-service/backend/internal/app"
	"go.uber.org/zap"
)

func main() {
	application, err := iapp.BootstrapWorker()
	if err != nil {
		panic(err)
	}

	logger := application.Logger
	defer func() {
		_ = logger.Sync()
	}()

	run(application)

	if err := iapp.Run(application); err != nil {
		logger.Fatal("failed to run worker lifecycle", zap.Error(err))
	}
}
