package errors

import (
	stderrors "errors"
	"net/http"
)

func ToAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if stderrors.As(err, &appErr) {
		return appErr
	}

	return New(
		http.StatusInternalServerError,
		CodeInternal,
		"internal server error",
		nil,
	)
}
