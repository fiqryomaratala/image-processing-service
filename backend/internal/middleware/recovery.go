package middleware

import (
	"fmt"

	apperrors "github.com/fiqryomaratala/image-processing-service/backend/internal/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	_ = log

	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				panicErr := fmt.Errorf("panic recovered: %v", recovered)

				_ = c.Error(
					apperrors.Internal("Internal server error").WithCause(panicErr),
				)
				c.Abort()
			}
		}()

		c.Next()
	}
}
