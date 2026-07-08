package errors

import (
	"net/http"

	apperrors "github.com/fiqryomaratala/image-processing-service/backend/internal/errors"
)

const (
	CodeInvalidFileType = "INVALID_FILE_TYPE"
	CodeFileTooLarge    = "FILE_TOO_LARGE"
	CodeImageNotFound   = "IMAGE_NOT_FOUND"
)

func InvalidFileType() *apperrors.AppError {
	return apperrors.New(http.StatusBadRequest, CodeInvalidFileType, "invalid file type")
}

func FileTooLarge() *apperrors.AppError {
	return apperrors.New(http.StatusBadRequest, CodeFileTooLarge, "file too large")
}

func ImageNotFound() *apperrors.AppError {
	return apperrors.New(http.StatusNotFound, CodeImageNotFound, "image not found")
}
