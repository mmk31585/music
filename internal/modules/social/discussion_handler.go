package social

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func (h *Handler) GetClubDiscussionReplies(c *gin.Context) {
	id := c.Param("id")
	items, err := h.service.GetClubDiscussionReplies(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

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
