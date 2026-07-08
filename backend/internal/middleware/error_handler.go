package middleware

import (
	stderrors "errors"
	"net/http"

	apperrors "github.com/fiqryomaratala/image-processing-service/backend/internal/errors"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() || len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		handleError(c, log, err)
	}
}

func handleError(c *gin.Context, log *zap.Logger, err error) {
	var appErr *apperrors.AppError
	if stderrors.As(err, &appErr) {
		if appErr.StatusCode >= http.StatusInternalServerError {
			log.Error("request failed",
				zap.Error(err),
				zap.String("request_id", getRequestID(c)),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
			)
		}

		errItems := make([]response.ErrorItem, 0, len(appErr.Details))
		for _, detail := range appErr.Details {
			errItems = append(errItems, response.ErrorItem{
				Field:   detail.Field,
				Message: detail.Message,
				Code:    appErr.Code,
			})
		}

		if len(errItems) == 0 {
			errItems = append(errItems, response.ErrorItem{
				Message: appErr.Message,
				Code:    appErr.Code,
			})
		}

		switch appErr.StatusCode {
		case http.StatusBadRequest:
			response.BadRequest(c, appErr.Message, errItems)
		case http.StatusUnauthorized:
			response.Unauthorized(c, appErr.Message, errItems)
		case http.StatusForbidden:
			response.Forbidden(c, appErr.Message, errItems)
		case http.StatusNotFound:
			response.NotFound(c, appErr.Message, errItems)
		case http.StatusConflict:
			response.Conflict(c, appErr.Message, errItems)
		case http.StatusUnprocessableEntity:
			response.UnprocessableEntity(c, appErr.Message, errItems)
		default:
			response.InternalServerError(c, appErr.Message, errItems)
		}

		return
	}

	log.Error("unhandled request error",
		zap.Error(err),
		zap.String("request_id", getRequestID(c)),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	response.InternalServerError(c, "Internal server error", []response.ErrorItem{
		{
			Message: "Internal server error",
			Code:    apperrors.CodeInternalServer,
		},
	})
}
