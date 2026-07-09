package handler

import (
	"fmt"

	apperrors "github.com/fiqryomaratala/image-processing-service/backend/internal/errors"
	filepkg "github.com/fiqryomaratala/image-processing-service/backend/internal/file"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/dto"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Upload godoc
// @Summary Upload image for validation
// @Description Accepts a multipart image file and validates it without storing the file.
// @Tags Images
// @Accept mpfd
// @Produce json
// @Param file formData file true "Image file"
// @Success 200 {object} dto.UploadSuccessResponse
// @Failure 400 {object} response.ErrorBody
// @Failure 500 {object} response.ErrorBody
// @Router /api/v1/images/upload [post]
func (h *Handler) Upload(c *gin.Context) {
	h.logger.Info("upload request received",
		zap.String("request_id", c.GetString("request_id")),
	)

	header, err := c.FormFile("file")
	if err != nil {
		_ = c.Error(apperrors.Validation("invalid upload request", []apperrors.FieldError{
			apperrors.NewFieldError("file", "file is required"),
		}))
		return
	}

	openedFile, err := header.Open()
	if err != nil {
		_ = c.Error(apperrors.Internal("failed to open uploaded file").WithCause(fmt.Errorf("open multipart file: %w", err)))
		return
	}
	defer func() {
		_ = openedFile.Close()
	}()

	result, err := h.service.Upload(c.Request.Context(), dto.UploadRequest{
		File: filepkg.Upload{
			Filename: header.Filename,
			Size:     header.Size,
			Reader:   openedFile,
		},
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Success(c, "File validation passed", result, nil)
}
