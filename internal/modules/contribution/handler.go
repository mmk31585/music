package contribution

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"music/internal/platform/web"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		return
	}

	var req CreateContributionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request: " + err.Error()})
		return
	}

	result, err := h.service.Create(c.Request.Context(), req, userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	var data any
	_ = json.Unmarshal([]byte(result.Data), &data)

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": ToContributionResponse(result, data)})
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	result, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	var data any
	_ = json.Unmarshal([]byte(result.Data), &data)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": ToContributionResponse(result, data)})
}

func (h *Handler) ListMyContributions(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		return
	}

	page, pageSize := parsePagination(c)
	items, total, err := h.service.ListByUser(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		h.handleError(c, err)
		return
	}

	responses := make([]ContributionResponse, 0, len(items))
	for _, item := range items {
		var data any
		_ = json.Unmarshal([]byte(item.Data), &data)
		responses = append(responses, ToContributionResponse(item, data))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    responses,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *Handler) ListPending(c *gin.Context) {
	page, pageSize := parsePagination(c)
	items, total, err := h.service.ListPending(c.Request.Context(), page, pageSize)
	if err != nil {
		h.handleError(c, err)
		return
	}

	responses := make([]ContributionResponse, 0, len(items))
	for _, item := range items {
		var data any
		_ = json.Unmarshal([]byte(item.Data), &data)
		responses = append(responses, ToContributionResponse(item, data))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    responses,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *Handler) ListByTarget(c *gin.Context) {
	targetType := c.Param("targetType")
	targetID := c.Param("targetID")

	items, err := h.service.ListByTarget(c.Request.Context(), targetType, targetID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	responses := make([]ContributionResponse, 0, len(items))
	for _, item := range items {
		var data any
		_ = json.Unmarshal([]byte(item.Data), &data)
		responses = append(responses, ToContributionResponse(item, data))
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": responses})
}

func (h *Handler) Review(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	var req ReviewContributionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request: " + err.Error()})
		return
	}

	result, err := h.service.Review(c.Request.Context(), id, userID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	var data any
	_ = json.Unmarshal([]byte(result.Data), &data)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": ToContributionResponse(result, data)})
}

func (h *Handler) GetHistory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	items, err := h.service.GetHistory(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	responses := make([]ContributionHistoryResponse, 0, len(items))
	for _, item := range items {
		var data any
		_ = json.Unmarshal([]byte(item.Data), &data)
		var prev any
		if item.PreviousData != nil {
			_ = json.Unmarshal([]byte(*item.PreviousData), &prev)
		}
		responses = append(responses, ToHistoryResponse(item, data, prev))
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": responses})
}

func (h *Handler) GetContentVersions(c *gin.Context) {
	targetType := c.Param("targetType")
	targetID := c.Param("targetID")

	items, err := h.service.GetContentVersions(c.Request.Context(), targetType, targetID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	versionResponses := make([]gin.H, 0, len(items))
	for _, v := range items {
		var data any
		_ = json.Unmarshal([]byte(v.Data), &data)
		versionResponses = append(versionResponses, gin.H{
			"id":         v.ID,
			"version":    v.Version,
			"data":       data,
			"applied_by": v.AppliedBy,
			"created_at": v.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": versionResponses})
}

func (h *Handler) ListByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "status query param required"})
		return
	}

	contributionType := c.Query("type")
	search := c.Query("q")
	page, pageSize := parsePagination(c)
	items, total, err := h.service.ListByStatusFiltered(c.Request.Context(), status, contributionType, search, page, pageSize)
	if err != nil {
		h.handleError(c, err)
		return
	}

	responses := make([]ContributionResponse, 0, len(items))
	for _, item := range items {
		var data any
		_ = json.Unmarshal([]byte(item.Data), &data)
		responses = append(responses, ToContributionResponse(item, data))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    responses,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

func (h *Handler) Apply(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "contribution id required"})
		return
	}

	result, err := h.service.Apply(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result, "message": "contribution applied to target"})
}

func (h *Handler) GetLeaderboard(c *gin.Context) {
	limit := 20
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	items, err := h.service.GetLeaderboard(c.Request.Context(), limit)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrContributionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "contribution not found"})
	case errors.Is(err, ErrInvalidTarget):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid target"})
	case errors.Is(err, ErrAlreadyDecided):
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "contribution already decided"})
	case errors.Is(err, ErrForbiddenReview):
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "not authorized to review"})
	case errors.Is(err, ErrNotApproved):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "contribution is not approved"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
	}
}

func parsePagination(c *gin.Context) (page, pageSize int) {
	page = 1
	pageSize = 20
	if p, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(c.DefaultQuery("page_size", "20")); err == nil && ps > 0 && ps <= 100 {
		pageSize = ps
	}
	return
}
