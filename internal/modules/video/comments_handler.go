package video

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "music/internal/common/errors"
	"music/internal/common/response"
	"music/internal/modules/auth"
)

// ListComments handles GET /api/v1/videos/:id/comments
func (h *Handler) ListComments(c *gin.Context) {
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid video id", err))
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	comments, total, err := h.commentRepo.ListByVideo(c.Request.Context(), videoID, limit, offset)
	if err != nil {
		response.Error(c, apperrors.Internal("failed to load comments", err))
		return
	}

	// Convert to response objects with author names
	var items []CommentResponse
	for _, cm := range comments {
		author, _ := h.commentRepo.GetUserDisplayName(c.Request.Context(), cm.UserID)
		items = append(items, CommentResponse{
			ID:        cm.ID.String(),
			Author:    author,
			Content:   cm.Content,
			CreatedAt: cm.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	if items == nil {
		items = []CommentResponse{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": items,
			"total": total,
		},
	})
}

// CreateComment handles POST /api/v1/videos/:id/comments
func (h *Handler) CreateComment(c *gin.Context) {
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid video id", err))
		return
	}

	userIDStr := auth.UserIDFromContext(c)
	if userIDStr == "" {
		response.Error(c, apperrors.Unauthorized("authentication required", nil))
		return
	}
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(c, apperrors.Unauthorized("invalid user id", nil))
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request", err))
		return
	}

	cm, err := h.commentRepo.Create(c.Request.Context(), videoID, uid, req.Content)
	if err != nil {
		response.Error(c, apperrors.Internal("failed to create comment", err))
		return
	}

	author, _ := h.commentRepo.GetUserDisplayName(c.Request.Context(), uid)
	resp := CommentResponse{
		ID:        cm.ID.String(),
		Author:    author,
		Content:   cm.Content,
		CreatedAt: cm.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	response.Success(c, http.StatusCreated, "comment created", resp)
}

// AdminListVideoComments handles GET /admin/videos/:id/comments
func (h *Handler) AdminListVideoComments(c *gin.Context) {
	videoIDStr := c.Param("id")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid video id", err))
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit < 1 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := h.commentRepo.AdminListByVideo(c.Request.Context(), videoID, limit, offset)
	if err != nil {
		response.Error(c, apperrors.Internal("failed to load comments", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": items,
			"total": total,
		},
	})
}

// AdminUpdateComment handles PUT /admin/comments/:id
func (h *Handler) AdminUpdateComment(c *gin.Context) {
	idStr := c.Param("id")
	commentID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid comment id", err))
		return
	}

	var req AdminUpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request", err))
		return
	}

	comment, err := h.commentRepo.UpdateComment(c.Request.Context(), commentID, req.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Error(c, apperrors.NotFound("comment not found", nil))
			return
		}
		response.Error(c, apperrors.Internal("failed to update comment", err))
		return
	}

	author, _ := h.commentRepo.GetUserDisplayName(c.Request.Context(), comment.UserID)
	response.Success(c, http.StatusOK, "comment updated", CommentResponse{
		ID:        comment.ID.String(),
		Author:    author,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// AdminDeleteComment handles DELETE /admin/comments/:id
func (h *Handler) AdminDeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	commentID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, apperrors.BadRequest("invalid comment id", err))
		return
	}

	if err := h.commentRepo.DeleteComment(c.Request.Context(), commentID); err != nil {
		response.Error(c, apperrors.Internal("failed to delete comment", err))
		return
	}

	response.Success[any](c, http.StatusOK, "comment deleted", nil)
}
