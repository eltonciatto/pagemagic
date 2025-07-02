package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"pagemagic/billing-svc/internal/models"
	"pagemagic/billing-svc/internal/services"

	"github.com/gin-gonic/gin"
)

type BillingHandler struct {
	billingService *services.BillingService
	webhookSecret  string
}

func NewBillingHandler(billingService *services.BillingService, webhookSecret string) *BillingHandler {
	return &BillingHandler{
		billingService: billingService,
		webhookSecret:  webhookSecret,
	}
}

// HealthCheck godoc
// @Summary Health check
// @Description Health check endpoint
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *BillingHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "billing-svc",
	})
}

// CreateCustomer godoc
// @Summary Create Stripe customer
// @Description Create a new Stripe customer for a user
// @Tags billing
// @Accept json
// @Produce json
// @Param request body CreateCustomerRequest true "Customer creation request"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/customers [post]
func (h *BillingHandler) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.billingService.CreateCustomer(context.Background(), req.UserID, req.Email)
	if err != nil {
		log.Printf("Failed to create customer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// CreateSubscription godoc
// @Summary Create subscription
// @Description Create a new subscription for a user
// @Tags billing
// @Accept json
// @Produce json
// @Param request body CreateSubscriptionRequest true "Subscription creation request"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/subscriptions [post]
func (h *BillingHandler) CreateSubscription(c *gin.Context) {
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.billingService.CreateSubscription(context.Background(), req.UserID, req.PlanID)
	if err != nil {
		log.Printf("Failed to create subscription: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// GetUserSubscriptions godoc
// @Summary Get user subscriptions
// @Description Get all subscriptions for a user
// @Tags billing
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} models.Subscription
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/users/{user_id}/subscriptions [get]
func (h *BillingHandler) GetUserSubscriptions(c *gin.Context) {
	userID := c.Param("user_id")

	subscriptions, err := h.billingService.GetUserSubscriptions(context.Background(), userID)
	if err != nil {
		log.Printf("Failed to get user subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscriptions"})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

// HandleWebhook godoc
// @Summary Handle Stripe webhook
// @Description Handle incoming Stripe webhook events
// @Tags billing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /webhooks/stripe [post]
func (h *BillingHandler) HandleWebhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read webhook payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read payload"})
		return
	}

	// TODO: Verify webhook signature using h.webhookSecret

	var event models.WebhookPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Printf("Failed to parse webhook payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Save event to database for audit trail
	if err := h.billingService.SaveBillingEvent(context.Background(), event.ID, event.Type, event.Data); err != nil {
		log.Printf("Failed to save billing event: %v", err)
		// Don't return error - we still want to process the event
	}

	// Process the webhook event
	if err := h.billingService.HandleWebhook(context.Background(), event.Type, event.Data); err != nil {
		log.Printf("Failed to handle webhook %s: %v", event.Type, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// Request structures
type CreateCustomerRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
}

type CreateSubscriptionRequest struct {
	UserID string `json:"user_id" binding:"required"`
	PlanID string `json:"plan_id" binding:"required"`
}
