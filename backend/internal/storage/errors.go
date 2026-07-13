package storage

import (
	"net/http"

	apperrors "github.com/fiqryomaratala/image-processing-service/backend/internal/errors"
)

const (
	CodeStorageUnavailable = "STORAGE_UNAVAILABLE"
	CodeBucketUnavailable  = "BUCKET_UNAVAILABLE"
	CodeUploadFailed       = "UPLOAD_FAILED"
	CodeObjectNotFound     = "OBJECT_NOT_FOUND"
)

func StorageUnavailable(message string) *apperrors.AppError {
	return apperrors.New(http.StatusServiceUnavailable, CodeStorageUnavailable, message)
}

func BucketUnavailable(bucket string) *apperrors.AppError {
	return apperrors.New(http.StatusBadGateway, CodeBucketUnavailable, "storage bucket is not available").WithDetails([]apperrors.FieldError{
		apperrors.NewFieldError("bucket", bucket),
	})
}

func UploadFailed(message string) *apperrors.AppError {
	return apperrors.New(http.StatusBadGateway, CodeUploadFailed, message)
}

func ObjectNotFound(objectKey string) *apperrors.AppError {
	return apperrors.New(http.StatusNotFound, CodeObjectNotFound, "object not found").WithDetails([]apperrors.FieldError{
		apperrors.NewFieldError("object_key", objectKey),
	})
}
