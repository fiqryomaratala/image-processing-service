package main

import (
	"log"
	"net/http"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/shared"
	"github.com/gin-gonic/gin"
)

func run(cfg config.Config, logger *log.Logger) error {
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

	logger.Printf("api server listening on %s", cfg.App.Address())

	return router.Run(cfg.App.Address())
}
