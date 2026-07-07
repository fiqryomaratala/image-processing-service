package main

import (
	"net/http"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/bootstrap"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/shared"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func run(app *bootstrap.App) error {
	cfg := app.Config
	logger := app.Logger

	if cfg.App.Env == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, shared.HTTPResponse{
			Status:  "success",
			Message: "Image Processing Service API is running",
		})
	})

	logger.Info("api server listening", zap.String("address", cfg.App.Address()))

	return router.Run(cfg.App.Address())
}