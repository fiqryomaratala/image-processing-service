package handler

import (
	apperrors "github.com/fiqryomaratala/image-processing-service/backend/internal/errors"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/response"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Get godoc
// @Summary Health check
// @Description Returns the current API health status.
// @Tags System
// @Produce json
// @Param fail query bool false "Trigger validation error example"
// @Success 200 {object} response.HealthSuccessResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/v1/health [get]
func (h *HealthHandler) Get(c *gin.Context) {
	if c.Query("fail") == "true" {
		_ = c.Error(
			apperrors.Validation(
				"Validation failed",
				[]apperrors.FieldError{
					apperrors.NewFieldError("fail", "fail must not be true"),
				},
			),
		)
		return
	}

	response.Success(c, "Image Processing Service API is running", gin.H{
		"status": "healthy",
	}, nil)
}
