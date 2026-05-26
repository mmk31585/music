package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func OK(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusOK, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func OKWithMeta(w http.ResponseWriter, message string, data interface{}, meta interface{}) {
	JSON(w, http.StatusOK, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Created(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusCreated, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}
