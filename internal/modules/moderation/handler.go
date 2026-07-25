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

// Report godoc
// @Summary Report content
// @Description Submits a content report for moderation review.
// @Tags moderation
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body object{target_id=string,target_type=string,reason=string,description=string} true "Report payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /moderation/report [post]
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

// ListPending godoc
// @Summary List pending reports
// @Description Returns a paginated list of pending content reports (moderator/admin only).
// @Tags moderation
// @Produce json
// @Security Bearer
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /moderation/pending [get]
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

// ListByStatus godoc
// @Summary List reports by status
// @Description Returns a paginated list of reports filtered by status (moderator/admin only).
// @Tags moderation
// @Produce json
// @Security Bearer
// @Param status path string true "Report status (pending, resolved, dismissed)"
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /moderation/status/{status} [get]
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

// GetReport godoc
// @Summary Get report by ID
// @Description Returns a single content report by its ID.
// @Tags moderation
// @Produce json
// @Security Bearer
// @Param id path string true "Report ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /moderation/reports/{id} [get]
func (h *Handler) GetReport(c *gin.Context) {
	id := c.Param("id")
	report, err := h.service.GetReportByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": report})
}

// Resolve godoc
// @Summary Resolve a report
// @Description Resolves a content report with a status and optional note (moderator/admin only).
// @Tags moderation
// @Produce json
// @Security Bearer
// @Param id path string true "Report ID"
// @Param status query string false "Resolution status" default(resolved)
// @Param note query string false "Resolution note"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /moderation/resolve/{id} [post]
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

// BulkAction godoc
// @Summary Bulk resolve reports
// @Description Resolves multiple reports in bulk (moderator/admin only).
// @Tags moderation
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body object{ids=[]string,action=string,note=string} true "Bulk action payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /moderation/bulk [post]
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

// FlagContent godoc
// @Summary Flag content
// @Description Flags a target content item for review or automatic action.
// @Tags moderation
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body object{target_id=string,target_type=string,flag_type=string,expires_in_hours=int} true "Flag payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /moderation/flag [post]
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

// ListFlags godoc
// @Summary List content flags
// @Description Returns a paginated list of content flags, optionally including expired ones.
// @Tags moderation
// @Produce json
// @Security Bearer
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Param expired query bool false "Include expired flags" default(false)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /moderation/flags [get]
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

// GetStats godoc
// @Summary Get moderation stats
// @Description Returns moderation statistics including total reports, pending, resolved today (moderator/admin only).
// @Tags moderation
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /moderation/stats [get]
func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// GetActions godoc
// @Summary Get moderation actions
// @Description Returns a paginated list of moderation actions (moderator/admin only).
// @Tags moderation
// @Produce json
// @Security Bearer
// @Param report_id query string false "Filter by report ID"
// @Param limit query int false "Items per page" default(50)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /moderation/actions [get]
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
