package middleware

import (
	"time"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	log := logger.Get()

	return func(c *gin.Context) {
		startedAt := time.Now()

		c.Next()

		fields := []zap.Field{
			zap.String("request_id", getRequestID(c)),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status_code", c.Writer.Status()),
			zap.Duration("response_time", time.Since(startedAt)),
			zap.String("client_ip", c.ClientIP()),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		log.Info("http request completed", fields...)
	}
}

func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if value, ok := requestID.(string); ok {
			return value
		}
	}

	return ""
}
