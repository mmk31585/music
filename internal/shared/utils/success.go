package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse is a Swagger-friendly non-generic success response schema.
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"request completed successfully"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// SuccessResponseData is an optional generic runtime response type.
// You can keep using this in code, but Swagger annotations should use SuccessResponse.
type SuccessResponseData[T any] struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    T           `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// Success sends a standardized success response with generic typed data.
func Success[T any](c *gin.Context, statusCode int, message string, data T) {
	c.JSON(statusCode, SuccessResponseData[T]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta sends a success response including metadata.
func SuccessWithMeta[T any](c *gin.Context, statusCode int, message string, data T, meta interface{}) {
	c.JSON(statusCode, SuccessResponseData[T]{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// SuccessWithMap sends a success response using gin.H for dynamic payloads.
func SuccessWithMap(c *gin.Context, statusCode int, message string, data gin.H) {
	c.JSON(statusCode, SuccessResponseData[gin.H]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// OK is a convenience wrapper for 200 OK responses.
func OK[T any](c *gin.Context, message string, data T) {
	Success(c, http.StatusOK, message, data)
}

// Created is a convenience wrapper for 201 Created responses.
func Created[T any](c *gin.Context, message string, data T) {
	Success(c, http.StatusCreated, message, data)
}

// SuccessNoContent sends a 204 No Content response.
func SuccessNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Redirect sends an HTTP redirect using Gin's context.
func Redirect(c *gin.Context, statusCode int, url string) {
	c.Redirect(statusCode, url)
}
