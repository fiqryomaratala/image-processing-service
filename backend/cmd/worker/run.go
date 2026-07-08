package main

import (
	iapp "github.com/fiqryomaratala/image-processing-service/backend/internal/app"
	"go.uber.org/zap"
)

func run(application *iapp.App) {
	application.Logger.Info(
		"Worker started",
		zap.String("service", application.Config.App.Name),
		zap.String("environment", application.Config.App.Env),
	)
}
