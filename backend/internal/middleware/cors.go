package middleware

import (
	"net/http"

	"github.com/fiqryomaratala/image-processing-service/backend/internal/config"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	cfg := config.Get()

	return func(c *gin.Context) {
		headers := c.Writer.Header()
		headers.Set("Access-Control-Allow-Origin", cfg.CORS.AllowedOrigins)
		headers.Set("Access-Control-Allow-Methods", cfg.CORS.AllowedMethods)
		headers.Set("Access-Control-Allow-Headers", cfg.CORS.AllowedHeaders)
		headers.Set("Access-Control-Expose-Headers", RequestIDHeader)
		headers.Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
