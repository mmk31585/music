package social

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateClub godoc
// @Summary Create a club
// @Description Creates a new music club for the authenticated user.
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateClubRequest true "Create club request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/clubs [post]
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

// GetClub godoc
// @Summary Get club by ID
// @Description Returns a music club by its ID.
// @Tags social
// @Produce json
// @Param id path string true "Club ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/clubs/{id} [get]
func (h *Handler) GetClub(c *gin.Context) {
	id := c.Param("id")
	club, err := h.service.GetClub(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "club not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": club})
}

// ListClubs godoc
// @Summary List clubs
// @Description Returns a paginated list of music clubs.
// @Tags social
// @Produce json
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/clubs [get]
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

// JoinClub godoc
// @Summary Join a club
// @Description Allows the authenticated user to join a club.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Club ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/clubs/{id}/join [post]
func (h *Handler) JoinClub(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.service.JoinClub(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join club"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// LeaveClub godoc
// @Summary Leave a club
// @Description Allows the authenticated user to leave a club.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Club ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/clubs/{id}/leave [post]
func (h *Handler) LeaveClub(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	_ = h.service.LeaveClub(c.Request.Context(), id, userID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetClubMembers godoc
// @Summary Get club members
// @Description Returns the members of a club.
// @Tags social
// @Produce json
// @Param id path string true "Club ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/clubs/{id}/members [get]
func (h *Handler) GetClubMembers(c *gin.Context) {
	id := c.Param("id")
	members, err := h.service.GetClubMembers(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get members"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": members})
}

// CreateClubPost godoc
// @Summary Create club post
// @Description Creates a new post in a club.
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Club ID"
// @Param request body ClubPostRequest true "Post content"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/clubs/{id}/posts [post]
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

// GetClubPosts godoc
// @Summary Get club posts
// @Description Returns paginated posts from a club.
// @Tags social
// @Produce json
// @Param id path string true "Club ID"
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/clubs/{id}/posts [get]
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

// ListClubsWithGenre godoc
// @Summary List clubs by genre
// @Description Returns a paginated list of music clubs filtered by genre.
// @Tags social
// @Produce json
// @Security Bearer
// @Param genre query string true "Genre filter"
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/clubs/browse [get]
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

// GetClubDetail godoc
// @Summary Get club detail
// @Description Returns detailed club information including membership status for the authenticated user.
// @Tags social
// @Produce json
// @Security Bearer
// @Param id path string true "Club ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/clubs/{id}/detail [get]
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

// LaunchPartyFromClub godoc
// @Summary Launch a listening party from a club
// @Description Launches a listening party within a club for the authenticated user.
// @Tags social
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Club ID"
// @Param request body LaunchPartyRequest true "Launch party request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /social/clubs/{id}/launch-party [post]
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
