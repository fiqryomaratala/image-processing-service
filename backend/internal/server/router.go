package server

import (
	"github.com/fiqryomaratala/image-processing-service/backend/internal/handler"
	"github.com/gin-gonic/gin"
)

const apiV1Prefix = "/api/v1"

func registerRoutes(router *gin.Engine, healthHandler *handler.HealthHandler) {
	v1 := router.Group(apiV1Prefix)
	{
		v1.GET("/health", healthHandler.Get)
	}
}
