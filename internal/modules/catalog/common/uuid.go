package common

import "github.com/google/uuid"

func ParseUUID(id string) (uuid.UUID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, ErrInvalidInput
	}
	return uid, nil
}
