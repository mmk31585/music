package response

import (
	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"request completed successfully"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

type SuccessResponseData[T any] struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    T           `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success[T any](c *gin.Context, statusCode int, message string, data T, meta ...interface{}) {
	resp := SuccessResponseData[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
	if len(meta) > 0 {
		resp.Meta = meta[0]
	}
	c.JSON(statusCode, resp)
}
