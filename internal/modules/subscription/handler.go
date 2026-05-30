package subscription

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListPlans godoc
// @Summary      List subscription plans
// @Description  Retrieve all available subscription plans
// @Tags         subscription
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  ErrorResponse
// @Router       /subscription/plans [get]
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
// @Summary      Get current subscription
// @Description  Get the active subscription of the authenticated user
// @Tags         subscription
// @Produce      json
// @Success      200  {object}  SubscriptionResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     Bearer
// @Router       /subscription/me [get]
func (h *Handler) CurrentSubscription(c *gin.Context) {
	userID, ok := getUserID(c)
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
// @Summary      List subscription history
// @Description  Retrieve all past and present subscriptions of the user
// @Tags         subscription
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     Bearer
// @Router       /subscription/history [get]
func (h *Handler) ListSubscriptions(c *gin.Context) {
	userID, ok := getUserID(c)
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
// @Summary      Create checkout
// @Description  Initiate a checkout for a subscription plan
// @Tags         subscription
// @Accept       json
// @Produce      json
// @Param        request  body      CheckoutRequest  true  "Plan ID"
// @Success      200      {object}  CheckoutResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     Bearer
// @Router       /subscription/checkout [post]
func (h *Handler) Checkout(c *gin.Context) {
	userID, ok := getUserID(c)
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
// @Summary      Cancel subscription
// @Description  Cancel the current active subscription
// @Tags         subscription
// @Produce      json
// @Success      200  {object}  CancelSubscriptionResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     Bearer
// @Router       /subscription/cancel [post]
func (h *Handler) CancelCurrentSubscription(c *gin.Context) {
	userID, ok := getUserID(c)
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
// @Summary      List payment history
// @Description  Retrieve all payments made by the user
// @Tags         subscription
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     Bearer
// @Router       /subscription/payments [get]
func (h *Handler) ListPayments(c *gin.Context) {
	userID, ok := getUserID(c)
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

func getUserID(c *gin.Context) (string, bool) {
	keys := []string{
		"userID",
		"userId",
		"user_id",
		"sub",
	}

	for _, key := range keys {
		value, exists := c.Get(key)
		if !exists {
			continue
		}

		switch v := value.(type) {
		case string:
			if v != "" {
				return v, true
			}
		}
	}

	return "", false
}
