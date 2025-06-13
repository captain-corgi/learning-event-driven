package main

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// ErrorType represents different types of errors in our system
type ErrorType string

const (
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound   ErrorType = "NOT_FOUND_ERROR"
	ErrorTypeConflict   ErrorType = "CONFLICT_ERROR"
	ErrorTypeInternal   ErrorType = "INTERNAL_ERROR"
)

// AppError represents a custom application error
type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Field   string    `json:"field,omitempty"`
	Cause   error     `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s (field: %s)", e.Type, e.Message, e.Field)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying cause error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// HTTPStatusCode returns the appropriate HTTP status code for the error type
func (e *AppError) HTTPStatusCode() int {
	switch e.Type {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Field:   field,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource, id string) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s with id '%s' not found", resource, id),
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
	}
}

// NewInternalError creates a new internal error with cause
func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Cause:   cause,
	}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, message string) error {
	return errors.Wrap(err, message)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
