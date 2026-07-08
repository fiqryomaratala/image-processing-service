package main

import (
	iapp "github.com/fiqryomaratala/image-processing-service/backend/internal/app"
	"go.uber.org/zap"
)

func main() {
	application, err := iapp.BootstrapAPI()
	if err != nil {
		panic(err)
	}

	logger := application.Logger
	defer func() {
		_ = logger.Sync()
	}()

	if err := iapp.Run(application); err != nil {
		logger.Fatal("failed to run api server", zap.Error(err))
	}
}
