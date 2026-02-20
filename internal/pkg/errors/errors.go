package errors

import (
	"errors"
	"fmt"
)

// AppError represents an application error with code and optional cause
type AppError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying cause error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New creates a new AppError
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap creates a new AppError with a cause
func Wrap(code ErrorCode, message string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Wrapf creates a new AppError with formatted message and cause
func Wrapf(code ErrorCode, format string, cause error, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Cause:   cause,
	}
}

// IsNotFound checks if an error is a "not found" error
func IsNotFound(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Code == NotFound
}

// IsInvalidInput checks if an error is an "invalid input" error
func IsInvalidInput(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.Code == InvalidInput
}

// Predefined errors for common cases
var (
	ErrProductNotFound   = New(NotFound, "Product not found")
	ErrInvalidPrice      = New(InvalidInput, "Invalid price value")
	ErrInvalidDateRange  = New(InvalidInput, "Invalid date range")
	ErrInvalidPagination = New(InvalidInput, "Invalid pagination parameters")
)
