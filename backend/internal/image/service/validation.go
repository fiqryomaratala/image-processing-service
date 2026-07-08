package service

import (
	"strings"

	apperrors "github.com/fiqryomaratala/image-processing-service/backend/internal/errors"
	"github.com/fiqryomaratala/image-processing-service/backend/internal/image/dto"
)

const maxFilenameLength = 255

func validateUploadRequest(request dto.UploadRequest) error {
	var fieldErrors []apperrors.FieldError

	if strings.TrimSpace(request.OriginalFilename) == "" {
		fieldErrors = append(fieldErrors, apperrors.NewFieldError("original_filename", "original filename is required"))
	}

	if len(strings.TrimSpace(request.OriginalFilename)) > maxFilenameLength {
		fieldErrors = append(fieldErrors, apperrors.NewFieldError("original_filename", "original filename must not exceed 255 characters"))
	}

	if strings.TrimSpace(request.StoredFilename) == "" {
		fieldErrors = append(fieldErrors, apperrors.NewFieldError("stored_filename", "stored filename is required"))
	}

	if len(strings.TrimSpace(request.StoredFilename)) > maxFilenameLength {
		fieldErrors = append(fieldErrors, apperrors.NewFieldError("stored_filename", "stored filename must not exceed 255 characters"))
	}

	if strings.TrimSpace(request.ObjectKey) == "" {
		fieldErrors = append(fieldErrors, apperrors.NewFieldError("object_key", "object key is required"))
	}

	if strings.TrimSpace(request.BucketName) == "" {
		fieldErrors = append(fieldErrors, apperrors.NewFieldError("bucket_name", "bucket name is required"))
	}

	if strings.TrimSpace(request.ContentType) == "" {
		fieldErrors = append(fieldErrors, apperrors.NewFieldError("content_type", "content type is required"))
	}

	if request.FileSize < 0 {
		fieldErrors = append(fieldErrors, apperrors.NewFieldError("file_size", "file size must be greater than or equal to 0"))
	}

	if request.Width < 0 {
		fieldErrors = append(fieldErrors, apperrors.NewFieldError("width", "width must be greater than or equal to 0"))
	}

	if request.Height < 0 {
		fieldErrors = append(fieldErrors, apperrors.NewFieldError("height", "height must be greater than or equal to 0"))
	}

	if len(fieldErrors) > 0 {
		return apperrors.Validation("invalid image upload request", fieldErrors)
	}

	return nil
}
