package service

import (
	apperrors "github.com/fiqryomaratala/image-processing-service/backend/internal/errors"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/dto"
)

func validateUploadRequest(request dto.UploadRequest) error {
	var fieldErrors []apperrors.FieldError

	if request.File.Reader == nil {
		fieldErrors = append(fieldErrors, apperrors.NewFieldError("file", "file is required"))
	}

	if request.File.Size < 0 {
		fieldErrors = append(fieldErrors, apperrors.NewFieldError("file", "file size must be greater than or equal to 0"))
	}

	if len(fieldErrors) > 0 {
		return apperrors.Validation("invalid image upload request", fieldErrors)
	}

	return nil
}
