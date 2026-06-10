package tips

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	var req struct {
		ArtistID    string  `json:"artist_id" binding:"required"`
		TrackID     *string `json:"track_id"`
		AmountCents int64   `json:"amount_cents" binding:"required,min=1000"`
		Currency    string  `json:"currency"`
		Message     string  `json:"message"`
		CallbackURL string  `json:"callback_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Currency == "" {
		req.Currency = "IRR"
	}

	tip, redirectURL, err := h.service.CreateTip(c.Request.Context(), userID, req.ArtistID, req.TrackID, req.AmountCents, req.Currency, req.Message, req.CallbackURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tip"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tip": tip, "redirect_url": redirectURL})
}

func (h *Handler) VerifyCallback(c *gin.Context) {
	tipID := c.Param("id")
	authority := c.DefaultQuery("Authority", c.DefaultQuery("id", ""))
	status := c.DefaultQuery("Status", c.DefaultQuery("status", "OK"))

	tip, err := h.service.VerifyTip(c.Request.Context(), tipID, authority, status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tip not found"})
		return
	}

	if tip.Status == "completed" {
		c.JSON(http.StatusOK, gin.H{"success": true, "tip": tip, "message": "Tip completed!"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "tip": tip, "message": "Tip failed or canceled"})
	}
}

func (h *Handler) ListSent(c *gin.Context) {
	userID := c.GetString("auth_user_id")
	tips, err := h.service.ListBySender(c.Request.Context(), userID, 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tips": tips})
}

func (h *Handler) ListReceived(c *gin.Context) {
	artistID := c.Param("artist_id")
	if artistID == "" {
		artistID = c.GetString("auth_user_id")
	}
	tips, err := h.service.ListByArtist(c.Request.Context(), artistID, 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	total, _ := h.service.TotalForArtist(c.Request.Context(), artistID)
	c.JSON(http.StatusOK, gin.H{"tips": tips, "total_cents": total})
}
