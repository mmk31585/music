package follow

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// FollowArtist godoc
// @Summary Follow an artist
// @Description Follows the artist specified by id for the authenticated user.
// @Tags follow
// @Accept json
// @Produce json
// @Param id path string true "Artist ID"
// @Success 200 {object} FollowResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security Bearer
// @Router /artists/{id}/follow [post]
func (h *Handler) FollowArtist(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	artistID := c.Param("id")

	if err := h.service.FollowArtist(c.Request.Context(), userID, artistID); err != nil {
		switch {
		case errors.Is(err, ErrInvalidArtistID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		case errors.Is(err, ErrArtistNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "artist not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to follow artist"})
		}
		return
	}

	c.JSON(http.StatusOK, FollowResponse{
		Message: "artist followed",
	})
}

// UnfollowArtist godoc
// @Summary Unfollow an artist
// @Description Unfollows the artist specified by id for the authenticated user.
// @Tags follow
// @Accept json
// @Produce json
// @Param id path string true "Artist ID"
// @Success 200 {object} FollowResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security Bearer
// @Router /artists/{id}/follow [delete]
func (h *Handler) UnfollowArtist(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	artistID := c.Param("id")

	if err := h.service.UnfollowArtist(c.Request.Context(), userID, artistID); err != nil {
		switch {
		case errors.Is(err, ErrInvalidArtistID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unfollow artist"})
		}
		return
	}

	c.JSON(http.StatusOK, FollowResponse{
		Message: "artist unfollowed",
	})
}

// FollowUser godoc
// @Summary Follow a user
// @Description Follows the user specified by id for the authenticated user.
// @Tags follow
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} FollowResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security Bearer
// @Router /users/{id}/follow [post]
func (h *Handler) FollowUser(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	targetUserID := c.Param("id")

	if err := h.service.FollowUser(c.Request.Context(), userID, targetUserID); err != nil {
		switch {
		case errors.Is(err, ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		case errors.Is(err, ErrCannotSelfFollow):
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot follow yourself"})
		case errors.Is(err, ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to follow user"})
		}
		return
	}

	c.JSON(http.StatusOK, FollowResponse{
		Message: "user followed",
	})
}

// UnfollowUser godoc
// @Summary Unfollow a user
// @Description Unfollows the user specified by id for the authenticated user.
// @Tags follow
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} FollowResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security Bearer
// @Router /users/{id}/follow [delete]
func (h *Handler) UnfollowUser(c *gin.Context) {
	userID, ok := getUserIDFromGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	targetUserID := c.Param("id")

	if err := h.service.UnfollowUser(c.Request.Context(), userID, targetUserID); err != nil {
		switch {
		case errors.Is(err, ErrInvalidUserID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		case errors.Is(err, ErrCannotSelfFollow):
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot unfollow yourself"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unfollow user"})
		}
		return
	}

	c.JSON(http.StatusOK, FollowResponse{
		Message: "user unfollowed",
	})
}

func getUserIDFromGin(c *gin.Context) (uuid.UUID, bool) {
	value, exists := c.Get("auth_user_id")
	if !exists {
		value, exists = c.Get("userID")
		if !exists {
			value, exists = c.Get("user_id")
			if !exists {
				return uuid.Nil, false
			}
		}
	}

	switch v := value.(type) {
	case uuid.UUID:
		return v, true
	case string:
		parsed, err := uuid.Parse(v)
		if err == nil {
			return parsed, true
		}
	}
	return uuid.Nil, false
}
