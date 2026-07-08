package middleware

import (
	"net/http"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/shared"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	log := logger.Get()

	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Error("panic recovered",
			zap.Any("error", recovered),
			zap.String("request_id", getRequestID(c)),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)

		c.AbortWithStatusJSON(http.StatusInternalServerError, shared.HTTPResponse{
			Success: false,
			Message: "Internal server error",
			Data: gin.H{
				"request_id": getRequestID(c),
			},
		})
	})
}
