package social

import (
	"errors"
	"net/http"
	"strconv"

	appErr "music/internal/common/errors"
	"music/internal/common/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service      *Service
	queueEngine  *QueueEngine
	stageManager *StageManager
	logger       *zap.Logger
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service, logger: zap.L()}
}

func NewHandlerWithEngine(service *Service, queueEngine *QueueEngine) *Handler {
	return &Handler{service: service, queueEngine: queueEngine, logger: zap.L()}
}

func NewHandlerFull(service *Service, queueEngine *QueueEngine, stageManager *StageManager) *Handler {
	return &Handler{service: service, queueEngine: queueEngine, stageManager: stageManager, logger: zap.L()}
}

func NewHandlerWithLogger(service *Service, queueEngine *QueueEngine, stageManager *StageManager, logger *zap.Logger) *Handler {
	return &Handler{service: service, queueEngine: queueEngine, stageManager: stageManager, logger: logger}
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
		h.logger.Error("update party status failed", zap.String("party_id", id), zap.String("status", req.Status), zap.Error(err))
		response.Error(c, appErr.Internal("failed to update party", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) JoinParty(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.service.JoinParty(c.Request.Context(), id, userID); err != nil {
		h.logger.Error("join party failed", zap.String("party_id", id), zap.String("user_id", userID), zap.Error(err))
		response.Error(c, appErr.Internal("failed to join party", err))
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

// --- Democratic Voting Queue ---

func (h *Handler) SuggestTrack(c *gin.Context) {
	if h.queueEngine == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "queue engine not available"})
		return
	}
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	var req struct {
		TrackID string `json:"track_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.queueEngine.SuggestTrack(c.Request.Context(), id, req.TrackID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) CastVote(c *gin.Context) {
	if h.queueEngine == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "queue engine not available"})
		return
	}
	id := c.Param("id")
	candidateID := c.Param("candidateId")
	userID := c.GetString("auth_user_id")
	if err := h.queueEngine.CastVote(c.Request.Context(), candidateID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": id})
}

func (h *Handler) RemoveVote(c *gin.Context) {
	if h.queueEngine == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "queue engine not available"})
		return
	}
	candidateID := c.Param("candidateId")
	userID := c.GetString("auth_user_id")
	if err := h.queueEngine.RemoveVote(c.Request.Context(), candidateID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetQueueState(c *gin.Context) {
	if h.queueEngine == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "queue engine not available"})
		return
	}
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	state, err := h.queueEngine.GetQueueState(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": state})
}

func (h *Handler) TrackEnded(c *gin.Context) {
	if h.queueEngine == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "queue engine not available"})
		return
	}
	id := c.Param("id")
	var req struct {
		TrackID string `json:"track_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.queueEngine.TrackEnded(c.Request.Context(), id, req.TrackID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
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

func (h *Handler) RaiseHand(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.stageManager.RaiseHand(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) LowerHand(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	id := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.stageManager.LowerHand(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ApproveHand(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	roomID := c.Param("id")
	targetUserID := c.Param("userId")
	hostID := c.GetString("auth_user_id")
	if err := h.stageManager.ApproveHand(c.Request.Context(), roomID, hostID, targetUserID); err != nil {
		if err == ErrNotHost {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) DenyHand(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	roomID := c.Param("id")
	targetUserID := c.Param("userId")
	hostID := c.GetString("auth_user_id")
	if err := h.stageManager.DenyHand(c.Request.Context(), roomID, hostID, targetUserID); err != nil {
		if err == ErrNotHost {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) RemoveFromStage(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	roomID := c.Param("id")
	targetUserID := c.Param("userId")
	hostID := c.GetString("auth_user_id")
	if err := h.stageManager.RemoveFromStage(c.Request.Context(), roomID, hostID, targetUserID); err != nil {
		if err == ErrNotHost {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) LeaveStage(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	roomID := c.Param("id")
	userID := c.GetString("auth_user_id")
	if err := h.stageManager.LeaveStage(c.Request.Context(), roomID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ToggleMute(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	roomID := c.Param("id")
	targetUserID := c.Param("userId")
	actorID := c.GetString("auth_user_id")
	var req struct {
		Muted bool `json:"muted"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.stageManager.ToggleMute(c.Request.Context(), roomID, actorID, targetUserID, req.Muted); err != nil {
		if err == ErrNotHost {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetStageState(c *gin.Context) {
	if h.stageManager == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "stage manager not available"})
		return
	}
	roomID := c.Param("id")
	userID := c.GetString("auth_user_id")
	state, err := h.stageManager.GetStageState(c.Request.Context(), roomID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": state})
}

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
