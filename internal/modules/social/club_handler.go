package social

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateClub(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	var req CreateClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	club, err := h.service.CreateClub(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": club})
}

func (h *Handler) GetClub(c *gin.Context) {
	id := c.Param("id")
	club, err := h.service.GetClub(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "club not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": club})
}

func (h *Handler) ListClubs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	clubs, err := h.service.ListClubs(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list clubs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": clubs})
}

func (h *Handler) JoinClub(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.service.JoinClub(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join club"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) LeaveClub(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	_ = h.service.LeaveClub(c.Request.Context(), id, userID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetClubMembers(c *gin.Context) {
	id := c.Param("id")
	members, err := h.service.GetClubMembers(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get members"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": members})
}

func (h *Handler) CreateClubPost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	var req ClubPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post, err := h.service.CreateClubPost(c.Request.Context(), id, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": post})
}

func (h *Handler) GetClubPosts(c *gin.Context) {
	id := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	posts, err := h.service.GetClubPosts(c.Request.Context(), id, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get posts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": posts})
}

// --- Club Enhancements (Phase 5) ---

func (h *Handler) ListClubsWithGenre(c *gin.Context) {
	genre := c.Query("genre")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	clubs, err := h.service.ListClubsWithGenre(c.Request.Context(), genre, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list clubs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": clubs})
}

func (h *Handler) GetClubDetail(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	detail, err := h.service.GetClubDetail(c.Request.Context(), id, userID)
	if err != nil {
		code := http.StatusInternalServerError
		if err == ErrClubNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": detail})
}

func (h *Handler) LaunchPartyFromClub(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	var req LaunchPartyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	party, err := h.service.LaunchListeningParty(c.Request.Context(), id, userID, req)
	if err != nil {
		code := http.StatusInternalServerError
		if err == ErrLaunchNotMember || err == ErrClubNotFound {
			code = http.StatusForbidden
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": party})
}

// --- Discussions ---
