package moderation

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Report(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	var req struct {
		TargetID    string  `json:"target_id" binding:"required"`
		TargetType  string  `json:"target_type" binding:"required"`
		Reason      string  `json:"reason" binding:"required"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.ReportContent(c.Request.Context(), userID, req.TargetID, req.TargetType, req.Reason, req.Description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit report"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "report submitted"})
}

func (h *Handler) ListPending(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, err := h.service.GetPendingReports(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reports"})
		return
	}
	if items == nil {
		items = []ContentReport{}
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "limit": limit, "offset": offset})
}

func (h *Handler) ListByStatus(c *gin.Context) {
	status := c.Param("status")
	if status == "" {
		status = "pending"
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, err := h.service.GetReportsByStatus(c.Request.Context(), status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reports"})
		return
	}
	if items == nil {
		items = []ContentReport{}
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "limit": limit, "offset": offset})
}

func (h *Handler) GetReport(c *gin.Context) {
	id := c.Param("id")
	report, err := h.service.GetReportByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": report})
}

func (h *Handler) Resolve(c *gin.Context) {
	moderatorID := c.GetString("auth_user_id")
	reportID := c.Param("id")
	status := c.DefaultQuery("status", "resolved")
	note := c.DefaultQuery("note", "")
	if err := h.service.ResolveReport(c.Request.Context(), reportID, moderatorID, status, note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve report"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) BulkAction(c *gin.Context) {
	moderatorID := c.GetString("auth_user_id")
	var req struct {
		IDs    []string `json:"ids" binding:"required"`
		Action string   `json:"action" binding:"required"`
		Note   string   `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.BulkResolve(c.Request.Context(), req.IDs, moderatorID, req.Action, req.Note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed bulk action"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "count": len(req.IDs)})
}

func (h *Handler) FlagContent(c *gin.Context) {
	var req struct {
		TargetID       string `json:"target_id" binding:"required"`
		TargetType     string `json:"target_type" binding:"required"`
		FlagType       string `json:"flag_type" binding:"required"`
		ExpiresInHours int    `json:"expires_in_hours"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.FlagContent(c.Request.Context(), req.TargetID, req.TargetType, req.FlagType, req.ExpiresInHours); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to flag content"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ListFlags(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	includeExpired := c.DefaultQuery("expired", "false") == "true"
	items, err := h.service.GetAllFlags(c.Request.Context(), limit, offset, includeExpired)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get flags"})
		return
	}
	if items == nil {
		items = []ContentFlag{}
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "limit": limit, "offset": offset})
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *Handler) GetActions(c *gin.Context) {
	reportID := c.Query("report_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	items, err := h.service.GetActions(c.Request.Context(), reportID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get actions"})
		return
	}
	if items == nil {
		items = []ModerationAction{}
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "limit": limit, "offset": offset})
}

func parseLimitOffset(c *gin.Context) (int, int) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	return limit, offset
}

func sanitize(s string) string {
	return strings.TrimSpace(s)
}
