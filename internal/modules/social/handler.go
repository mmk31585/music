package social

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Follow(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	followedID := c.Param("userId")
	if followedID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id required"})
		return
	}
	if err := h.service.Follow(c.Request.Context(), userID, followedID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) Unfollow(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	followedID := c.Param("userId")
	_ = h.service.Unfollow(c.Request.Context(), userID, followedID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetFollowers(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		userID = c.GetString("auth_user_id")
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, total, err := h.service.GetFollowers(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get followers"})
		return
	}
	c.JSON(http.StatusOK, FollowersResponse{Items: items, TotalCount: total, Limit: limit, Offset: offset})
}

func (h *Handler) GetFollowing(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		userID = c.GetString("auth_user_id")
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, total, err := h.service.GetFollowing(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get following"})
		return
	}
	c.JSON(http.StatusOK, FollowersResponse{Items: items, TotalCount: total, Limit: limit, Offset: offset})
}

func (h *Handler) IsFollowing(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	targetID := c.Param("userId")
	following, err := h.service.IsFollowing(c.Request.Context(), userID, targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "check failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_following": following})
}

func (h *Handler) GetFeed(c *gin.Context) {
	userID := c.DefaultQuery("user_id", c.GetString("auth_user_id"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	types := c.Query("types")
	items, err := h.service.GetFeed(c.Request.Context(), userID, limit, offset, types)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get feed"})
		return
	}
	hasMore := len(items) > limit
	if hasMore {
		items = items[:limit]
	}
	c.JSON(http.StatusOK, ActivityResponse{
		Items:      items,
		Pagination: Pagination{Limit: limit, Offset: offset, Count: len(items), HasMore: hasMore},
	})
}

// --- Listening Parties ---

func (h *Handler) CreateParty(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	var req CreatePartyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	party, err := h.service.CreateParty(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create party"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": party})
}

func (h *Handler) GetParty(c *gin.Context) {
	id := c.Param("id")
	party, err := h.service.GetParty(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "party not found"})
		return
	}
	participants, _ := h.service.GetPartyParticipants(c.Request.Context(), id)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": party, "participants": participants})
}

func (h *Handler) ListActiveParties(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	parties, err := h.service.ListActiveParties(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list parties"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": parties})
}

func (h *Handler) UpdatePartyStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status  string  `json:"status" binding:"required"`
		TrackID *string `json:"track_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdatePartyStatus(c.Request.Context(), id, req.Status, req.TrackID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update party"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) JoinParty(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.service.JoinParty(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join party"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) LeaveParty(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	_ = h.service.LeaveParty(c.Request.Context(), id, userID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// --- Live Rooms ---

func (h *Handler) CreateRoom(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	room, err := h.service.CreateRoom(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": room})
}

func (h *Handler) GetRoom(c *gin.Context) {
	id := c.Param("id")
	room, err := h.service.GetRoom(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": room})
}

func (h *Handler) ListActiveRooms(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	rooms, err := h.service.ListActiveRooms(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list rooms"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rooms})
}

func (h *Handler) JoinRoom(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.service.JoinRoom(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to join room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) LeaveRoom(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	_ = h.service.LeaveRoom(c.Request.Context(), id, userID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetRoomParticipants(c *gin.Context) {
	id := c.Param("id")
	participants, err := h.service.GetRoomParticipants(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get participants"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": participants})
}

func (h *Handler) AddToRoomQueue(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	var req AddToQueueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AddToRoomQueue(c.Request.Context(), id, req.TrackID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add to queue"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetRoomQueue(c *gin.Context) {
	id := c.Param("id")
	queue, err := h.service.GetRoomQueue(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get queue"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": queue})
}

// --- Music Clubs ---

func (h *Handler) CreateClub(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	var req CreateClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	club, err := h.service.CreateClub(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create club"})
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

// --- Discussions ---

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

// --- Track Ratings ---

func (h *Handler) CreateRating(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	var req CreateRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rating, err := h.service.CreateRating(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create rating"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": rating})
}

func (h *Handler) GetTrackRatings(c *gin.Context) {
	trackID := c.Param("trackId")
	ratings, err := h.service.GetTrackRatings(c.Request.Context(), trackID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get ratings"})
		return
	}
	avg, count, _ := h.service.GetTrackRatingAverage(c.Request.Context(), trackID)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ratings,
		"average": avg,
		"count":   count,
	})
}
