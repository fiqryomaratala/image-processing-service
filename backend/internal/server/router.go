package server

import (
	"github.com/fiqryomaratala/image-processing-service/backend/internal/handler"
	imagehandler "github.com/fiqryomaratala/image-processing-service/backend/internal/image/handler"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const apiV1Prefix = "/api/v1"

func registerRoutes(router *gin.Engine, healthHandler *handler.HealthHandler, imageHandler *imagehandler.Handler) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group(apiV1Prefix)
	{
		v1.GET("/health", healthHandler.Get)
		v1.POST("/images/upload", imageHandler.Upload)
	}
}
