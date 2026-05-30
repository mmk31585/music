package subscription

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"music/internal/platform/web"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListPlans(c *gin.Context) {
	plans, err := h.service.ListPlans(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to list plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plans": plans,
	})
}

func (h *Handler) CurrentSubscription(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	sub, err := h.service.CurrentSubscription(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to get subscription"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (h *Handler) ListSubscriptions(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	subs, err := h.service.ListSubscriptions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to list subscriptions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subscriptions": subs,
	})
}

func (h *Handler) Checkout(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	res, err := h.service.Checkout(c.Request.Context(), userID, req)
	if errors.Is(err, ErrPremiumNotAvailable) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "premium is not available yet"})
		return
	}

	if errors.Is(err, ErrInactivePlan) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "plan is inactive"})
		return
	}

	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "plan not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create checkout"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) CancelCurrentSubscription(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	res, err := h.service.CancelCurrentSubscription(c.Request.Context(), userID)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "active subscription not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to cancel subscription"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) ListPayments(c *gin.Context) {
	userID, ok := web.GetUserIDString(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	payments, err := h.service.ListPayments(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to list payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
	})
}
