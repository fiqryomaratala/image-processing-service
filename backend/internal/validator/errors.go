package validator

import (
	"net/http"

	apperrors "github.com/fiqryomaratala/image-processing-service/backend/internal/errors"
)

const (
	CodeInvalidExtension = "INVALID_EXTENSION"
	CodeInvalidMimeType  = "INVALID_MIME_TYPE"
	CodeInvalidImage     = "INVALID_IMAGE"
	CodeFileTooLarge     = "FILE_TOO_LARGE"
	CodeImageTooLarge    = "IMAGE_TOO_LARGE"
	CodeFilenameInvalid  = "FILENAME_INVALID"
)

func InvalidExtension() *apperrors.AppError {
	return apperrors.New(http.StatusBadRequest, CodeInvalidExtension, "invalid file extension")
}

func InvalidMimeType() *apperrors.AppError {
	return apperrors.New(http.StatusBadRequest, CodeInvalidMimeType, "invalid mime type")
}

func InvalidImage() *apperrors.AppError {
	return apperrors.New(http.StatusBadRequest, CodeInvalidImage, "invalid image file")
}

func FileTooLarge() *apperrors.AppError {
	return apperrors.New(http.StatusBadRequest, CodeFileTooLarge, "file size exceeds the allowed limit")
}

func ImageTooLarge() *apperrors.AppError {
	return apperrors.New(http.StatusBadRequest, CodeImageTooLarge, "image dimensions are out of the allowed range")
}

func FilenameInvalid() *apperrors.AppError {
	return apperrors.New(http.StatusBadRequest, CodeFilenameInvalid, "invalid filename")
}
