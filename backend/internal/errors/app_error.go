package apperrors

type AppError struct {
	StatusCode int
	Code       string
	Message    string
	Details    []FieldError
	Cause      error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func (e *AppError) WithDetails(details []FieldError) *AppError {
	e.Details = details
	return e
}

func (e *AppError) WithCause(err error) *AppError {
	e.Cause = err
	return e
}

func New(statusCode int, code, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}
