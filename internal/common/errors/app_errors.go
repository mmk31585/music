package errors

import "net/http"

type AppError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
	StatusCode int         `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(statusCode int, code, message string, details interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		StatusCode: statusCode,
	}
}

func Internal(message string, details interface{}) *AppError {
	if message == "" {
		message = "internal server error"
	}

	return New(http.StatusInternalServerError, CodeInternal, message, details)
}

func BadRequest(message string, details interface{}) *AppError {
	if message == "" {
		message = "bad request"
	}

	return New(http.StatusBadRequest, CodeBadRequest, message, details)
}

func Validation(message string, details interface{}) *AppError {
	if message == "" {
		message = "validation failed"
	}

	return New(http.StatusUnprocessableEntity, CodeValidation, message, details)
}

func Unauthorized(message string, details interface{}) *AppError {
	if message == "" {
		message = "unauthorized"
	}

	return New(http.StatusUnauthorized, CodeUnauthorized, message, details)
}

func Forbidden(message string, details interface{}) *AppError {
	if message == "" {
		message = "forbidden"
	}

	return New(http.StatusForbidden, CodeForbidden, message, details)
}

func NotFound(message string, details interface{}) *AppError {
	if message == "" {
		message = "resource not found"
	}

	return New(http.StatusNotFound, CodeNotFound, message, details)
}

func Conflict(message string, details interface{}) *AppError {
	if message == "" {
		message = "conflict"
	}

	return New(http.StatusConflict, CodeConflict, message, details)
}

func MethodNotAllowed(message string, details interface{}) *AppError {
	if message == "" {
		message = "method not allowed"
	}

	return New(http.StatusMethodNotAllowed, CodeMethodNotAllowed, message, details)
}
