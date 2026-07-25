package web

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserIDString(c *gin.Context) (string, bool) {
	keys := []string{"auth_user_id", "userID", "userId", "user_id", "sub"}

	for _, key := range keys {
		value, exists := c.Get(key)
		if !exists {
			continue
		}

		switch v := value.(type) {
		case string:
			v = strings.TrimSpace(v)
			if v != "" {
				return v, true
			}
		case uuid.UUID:
			if v != uuid.Nil {
				return v.String(), true
			}
		}
	}

	return "", false
}

func GetOptionalUserUUID(c *gin.Context) *uuid.UUID {
	keys := []string{"auth_user_id", "userID", "userId", "user_id", "sub"}

	for _, key := range keys {
		value, exists := c.Get(key)
		if !exists {
			continue
		}

		switch v := value.(type) {
		case uuid.UUID:
			if v != uuid.Nil {
				id := v
				return &id
			}
		case string:
			id, err := uuid.Parse(strings.TrimSpace(v))
			if err == nil {
				return &id
			}
		}
	}

	return nil
}

func GetRequiredUserUUID(c *gin.Context) (uuid.UUID, bool) {
	id := GetOptionalUserUUID(c)
	if id == nil {
		return uuid.Nil, false
	}
	return *id, true
}
