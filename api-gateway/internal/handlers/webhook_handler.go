package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/fortexa/api-gateway/internal/models"
)

// WebhookHandler handles webhook-related API endpoints
type WebhookHandler struct{}

// NewWebhookHandler creates a new WebhookHandler
func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

// RegisterWebhook handles the webhook registration request
// @Summary Register a webhook
// @Description Register a new webhook for event notifications
// @Tags webhooks
// @Accept json
// @Produce json
// @Param webhook body models.WebhookRequest true "Webhook Request"
// @Success 200 {object} models.WebhookResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/webhooks/register [post]
func (h *WebhookHandler) RegisterWebhook(c *gin.Context) {
	var req models.WebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the merchant ID from the request context (set by the auth middleware)
	// In a real implementation, this would come from the authenticated context
	// For the MVP, we'll use a placeholder value
	merchantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// Generate webhook ID
	webhookID := uuid.New()

	// Generate webhook secret
	secret, err := generateWebhookSecret()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate webhook secret"})
		return
	}

	// Create webhook record
	webhook := models.Webhook{
		ID:         webhookID,
		MerchantID: merchantID,
		URL:        req.URL,
		EventTypes: req.EventTypes,
		Status:     models.WebhookStatusActive,
		Secret:     secret,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// In a real implementation, this would save the webhook to a database
	// For the MVP, we'll just return the created webhook

	// Return the webhook response (without the secret)
	c.JSON(http.StatusOK, models.WebhookResponse{
		ID:         webhook.ID,
		MerchantID: webhook.MerchantID,
		URL:        webhook.URL,
		EventTypes: webhook.EventTypes,
		Status:     webhook.Status,
		CreatedAt:  webhook.CreatedAt,
	})
}

// GetWebhooks retrieves all webhooks for a merchant
// @Summary Get merchant webhooks
// @Description Get all webhooks registered for a merchant
// @Tags webhooks
// @Accept json
// @Produce json
// @Success 200 {array} models.WebhookResponse
// @Failure 500 {object} gin.H
// @Router /api/v1/webhooks [get]
func (h *WebhookHandler) GetWebhooks(c *gin.Context) {
	// Get the merchant ID from the request context (set by the auth middleware)
	// In a real implementation, this would come from the authenticated context
	// For the MVP, we'll use a placeholder value
	merchantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// In a real implementation, this would fetch the webhooks from a database
	// For the MVP, we'll return a mock response with a sample webhook
	webhooks := []models.WebhookResponse{
		{
			ID:         uuid.New(),
			MerchantID: merchantID,
			URL:        "https://example.com/webhook",
			EventTypes: []string{"payment.initiated", "payment.captured"},
			Status:     models.WebhookStatusActive,
			CreatedAt:  time.Now().Add(-24 * time.Hour),
		},
	}

	c.JSON(http.StatusOK, webhooks)
}

// generateWebhookSecret generates a new webhook secret for signing webhook payloads
func generateWebhookSecret() (string, error) {
	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return "", err
	}
	return "whsec_" + hex.EncodeToString(secretBytes), nil
}

// RegisterWebhookRoutes registers the webhook routes with the given router group
func RegisterWebhookRoutes(router *gin.RouterGroup) {
	h := NewWebhookHandler()

	webhooks := router.Group("/webhooks")
	{
		webhooks.POST("/register", h.RegisterWebhook)
		webhooks.GET("", h.GetWebhooks)
	}
} 