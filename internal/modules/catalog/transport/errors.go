package common

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource conflict")
	ErrInvalidInput = errors.New("invalid input")
	ErrForeignKey   = errors.New("foreign key violation")
	ErrInternal     = errors.New("internal error")
)
