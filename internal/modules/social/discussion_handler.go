package social

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateDiscussion godoc
// @Summary Create a discussion
// @Description Creates a new discussion for a target (track, album, playlist).
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateDiscussionRequest true "Create discussion request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/discussions [post]
func (h *Handler) CreateDiscussion(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	var req CreateDiscussionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d, err := h.service.CreateDiscussion(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create discussion"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": d})
}

// GetDiscussions godoc
// @Summary Get discussions for a target
// @Description Returns paginated discussions for a specific target (track, album, playlist).
// @Tags social
// @Produce json
// @Param target_type query string true "Target type (track, album, playlist)"
// @Param target_id query string true "Target ID"
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/discussions [get]
func (h *Handler) GetDiscussions(c *gin.Context) {
	targetType := c.Query("target_type")
	targetID := c.Query("target_id")
	if targetType == "" || targetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_type and target_id required"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, err := h.service.GetDiscussions(c.Request.Context(), targetType, targetID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get discussions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

// GetDiscussionReplies godoc
// @Summary Get discussion replies
// @Description Returns all replies for a discussion.
// @Tags social
// @Produce json
// @Param id path string true "Discussion ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/discussions/{id}/replies [get]
func (h *Handler) GetDiscussionReplies(c *gin.Context) {
	parentID := c.Param("id")
	replies, err := h.service.GetDiscussionReplies(c.Request.Context(), parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get replies"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": replies})
}

// --- Club Discussions (Phase 6) ---

// CreateClubDiscussion godoc
// @Summary Create club discussion
// @Description Creates a new discussion thread in a club.
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Club ID"
// @Param request body CreateClubDiscussionRequest true "Create discussion request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/clubs/{id}/discussions [post]
func (h *Handler) CreateClubDiscussion(c *gin.Context) {
	clubID := c.Param("id")
	userID := c.GetString("auth_user_id")
	var req CreateClubDiscussionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d, err := h.service.CreateClubDiscussion(c.Request.Context(), clubID, userID, req.Title, req.Body)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, ErrNotClubMember) {
			code = http.StatusForbidden
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": d})
}

// ListClubDiscussions godoc
// @Summary List club discussions
// @Description Returns paginated discussions for a club.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Club ID"
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/clubs/{id}/discussions [get]
func (h *Handler) ListClubDiscussions(c *gin.Context) {
	clubID := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, err := h.service.ListClubDiscussions(c.Request.Context(), clubID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

// GetClubDiscussion godoc
// @Summary Get club discussion
// @Description Returns a discussion by its ID.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Discussion ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/discussions/{id} [get]
func (h *Handler) GetClubDiscussion(c *gin.Context) {
	id := c.Param("id")
	d, err := h.service.GetClubDiscussion(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrDiscussionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "discussion not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": d})
}

// GetClubDiscussionReplies godoc
// @Summary Get club discussion replies
// @Description Returns all replies for a club discussion.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Discussion ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/discussions/{id}/replies [get]
func (h *Handler) GetClubDiscussionReplies(c *gin.Context) {
	id := c.Param("id")
	items, err := h.service.GetClubDiscussionReplies(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

// CreateClubDiscussionReply godoc
// @Summary Create club discussion reply
// @Description Creates a reply to a club discussion.
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Discussion ID"
// @Param request body CreateDiscussionReplyRequest true "Reply body"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/discussions/{id}/replies [post]
func (h *Handler) CreateClubDiscussionReply(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	var req CreateDiscussionReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reply, err := h.service.CreateClubDiscussionReply(c.Request.Context(), id, userID, req.Body)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, ErrNotClubMember) {
			code = http.StatusForbidden
		}
		if errors.Is(err, ErrDiscussionNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": reply})
}

// DeleteClubDiscussion godoc
// @Summary Delete club discussion
// @Description Deletes a club discussion (author or owner only).
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Discussion ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/discussions/{id} [delete]
func (h *Handler) DeleteClubDiscussion(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.service.DeleteClubDiscussion(c.Request.Context(), id, userID); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, ErrDiscussionNotFound) {
			code = http.StatusNotFound
		}
		if errors.Is(err, ErrNotAuthorOrOwner) {
			code = http.StatusForbidden
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteClubDiscussionReply godoc
// @Summary Delete club discussion reply
// @Description Deletes a reply from a club discussion (author or owner only).
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Discussion ID"
// @Param replyId path string true "Reply ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/discussions/{id}/replies/{replyId} [delete]
func (h *Handler) DeleteClubDiscussionReply(c *gin.Context) {
	replyID := c.Param("replyId")
	userID := c.GetString("auth_user_id")
	if err := h.service.DeleteClubDiscussionReply(c.Request.Context(), replyID, userID); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, ErrReplyNotFound) {
			code = http.StatusNotFound
		}
		if errors.Is(err, ErrNotAuthorOrOwner) {
			code = http.StatusForbidden
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// --- Track Ratings ---

// --- Stage & Raise-Hand ---
