package common

import (
	"errors"

	"github.com/jackc/pgconn"
)

func MapPGError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23505":
		return ErrConflict
	case "23503":
		return ErrForeignKey
	default:
		return err
	}
}
