package apperrors

import "net/http"

type FieldError struct {
	Field   string
	Message string
}

func Validation(message string, details []FieldError) *AppError {
	return New(http.StatusBadRequest, CodeValidationError, message).WithDetails(details)
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, CodeUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, CodeForbidden, message)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, CodeNotFound, message)
}

func Conflict(message string) *AppError {
	return New(http.StatusConflict, CodeConflict, message)
}

func Internal(message string) *AppError {
	return New(http.StatusInternalServerError, CodeInternalServer, message)
}

func NewFieldError(field, message string) FieldError {
	return FieldError{
		Field:   field,
		Message: message,
	}
}
