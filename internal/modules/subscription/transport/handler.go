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

// ListPlans godoc
// @Summary List subscription plans
// @Description Returns available subscription plans.
// @Tags subscription
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /subscription/plans [get]
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

// CurrentSubscription godoc
// @Summary Get current subscription
// @Description Returns the current subscription for the authenticated user.
// @Tags subscription
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} SubscriptionResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscription/me [get]
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

// ListSubscriptions godoc
// @Summary List subscription history
// @Description Returns subscription history for the authenticated user.
// @Tags subscription
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscription/history [get]
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

// Checkout godoc
// @Summary Create subscription checkout
// @Description Creates a checkout session or subscription payment for the authenticated user.
// @Tags subscription
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CheckoutRequest true "Checkout request"
// @Success 200 {object} CheckoutResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscription/checkout [post]
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

// CancelCurrentSubscription godoc
// @Summary Cancel current subscription
// @Description Cancels the authenticated user's active subscription.
// @Tags subscription
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} CancelSubscriptionResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscription/cancel [post]
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

// ListPayments godoc
// @Summary List payments
// @Description Returns payments for the authenticated user.
// @Tags subscription
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscription/payments [get]
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

// ListPayments godoc
// @Summary List payments for the authenticated user
// @Description Returns a list of payments associated with the user
// @Tags subscription
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=[]subscription.PaymentResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security Bearer
// @Router /subscription/payments [get]
