package handler

import (
	"net/http"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/shared"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, shared.HTTPResponse{
		Success: true,
		Message: "Image Processing Service API is running",
		Data: gin.H{
			"status": "healthy",
		},
	})
}
