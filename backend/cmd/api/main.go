package main

import (
	iapp "github.com/fiqryomaratala/image-processing-service/backend/internal/app"
	_ "github.com/fiqryomaratala/image-processing-service/backend/docs"
	"go.uber.org/zap"
)

// @title Image Processing Service API
// @version 1.0
// @description REST API documentation for the Image Processing Service project.
// @contact.name Image Processing Service Team
// @contact.email dev@image-processing.local
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /
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
