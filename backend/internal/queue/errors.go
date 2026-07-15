package queue

import (
	"net/http"

	apperrors "github.com/fiqryomaratala/image-processing-service/backend/internal/errors"
)

const (
	CodeQueueUnavailable = "QUEUE_UNAVAILABLE"
	CodeQueueNotFound    = "QUEUE_NOT_FOUND"
	CodePublishFailed    = "QUEUE_PUBLISH_FAILED"
)

func QueueUnavailable(message string) *apperrors.AppError {
	return apperrors.New(http.StatusServiceUnavailable, CodeQueueUnavailable, message)
}

func QueueNotFound(queueName string) *apperrors.AppError {
	return apperrors.New(http.StatusNotFound, CodeQueueNotFound, "queue not found").WithDetails([]apperrors.FieldError{
		apperrors.NewFieldError("queue", queueName),
	})
}

func PublishFailed(message string) *apperrors.AppError {
	return apperrors.New(http.StatusBadGateway, CodePublishFailed, message)
}
