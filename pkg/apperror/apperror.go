package apperror

import (
	"errors"
	"net/http"
)

// AppError represents a structured domain/application error with an HTTP status code.
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func BadRequest(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  http.StatusBadRequest,
		Err:     err,
	}
}

func NotFound(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  http.StatusNotFound,
		Err:     err,
	}
}

func Conflict(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  http.StatusConflict,
		Err:     err,
	}
}

func Internal(message string, err error) *AppError {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: message,
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}

// FromError extracts an *AppError from err or wraps it as an internal server error.
func FromError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return Internal("An unexpected error occurred", err)
}
