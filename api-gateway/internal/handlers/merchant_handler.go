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

// MerchantHandler handles merchant-related API endpoints
type MerchantHandler struct{}

// NewMerchantHandler creates a new MerchantHandler
func NewMerchantHandler() *MerchantHandler {
	return &MerchantHandler{}
}

// OnboardMerchant handles the merchant onboarding request
// @Summary Onboard a new merchant
// @Description Register a new merchant in the system
// @Tags merchants
// @Accept json
// @Produce json
// @Param merchant body models.MerchantRequest true "Merchant Request"
// @Success 200 {object} models.MerchantResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/merchants/onboard [post]
func (h *MerchantHandler) OnboardMerchant(c *gin.Context) {
	var req models.MerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate merchant ID
	merchantID := uuid.New()

	// Generate API keys
	apiKey, secretKey, err := generateMerchantKeys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate merchant keys"})
		return
	}

	// Create merchant record
	merchant := models.Merchant{
		ID:           merchantID,
		Name:         req.Name,
		BusinessName: req.BusinessName,
		Email:        req.Email,
		Phone:        req.Phone,
		Website:      req.Website,
		APIKey:       apiKey,
		SecretKey:    secretKey,
		Status:       models.MerchantStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// In a real implementation, this would save the merchant to a database
	// For the MVP, we'll just return the created merchant

	// Return the merchant response (without the secret key)
	c.JSON(http.StatusOK, models.MerchantResponse{
		ID:           merchant.ID,
		Name:         merchant.Name,
		BusinessName: merchant.BusinessName,
		Email:        merchant.Email,
		Phone:        merchant.Phone,
		Website:      merchant.Website,
		APIKey:       merchant.APIKey,
		Status:       merchant.Status,
		CreatedAt:    merchant.CreatedAt,
	})
}

// GetMerchant retrieves merchant details
// @Summary Get merchant details
// @Description Get the details of a specific merchant
// @Tags merchants
// @Accept json
// @Produce json
// @Param id path string true "Merchant ID"
// @Success 200 {object} models.MerchantResponse
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/merchants/{id} [get]
func (h *MerchantHandler) GetMerchant(c *gin.Context) {
	merchantID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid merchant ID"})
		return
	}

	// In a real implementation, this would fetch the merchant from a database
	// For the MVP, we'll return a mock response
	merchant := models.MerchantResponse{
		ID:           merchantID,
		Name:         "Test Merchant",
		BusinessName: "Test Business",
		Email:        "test@example.com",
		Phone:        "+919999999999",
		Website:      "https://example.com",
		APIKey:       "pk_test_123456",
		Status:       models.MerchantStatusActive,
		CreatedAt:    time.Now().Add(-24 * time.Hour),
	}

	c.JSON(http.StatusOK, merchant)
}

// generateMerchantKeys generates a new API key and secret key for the merchant
func generateMerchantKeys() (string, string, error) {
	// Generate API key (public)
	apiKeyBytes := make([]byte, 16)
	if _, err := rand.Read(apiKeyBytes); err != nil {
		return "", "", err
	}
	apiKey := "pk_test_" + hex.EncodeToString(apiKeyBytes)

	// Generate secret key (private)
	secretKeyBytes := make([]byte, 32)
	if _, err := rand.Read(secretKeyBytes); err != nil {
		return "", "", err
	}
	secretKey := "sk_test_" + hex.EncodeToString(secretKeyBytes)

	return apiKey, secretKey, nil
}

// RegisterMerchantRoutes registers the merchant routes with the given router group
func RegisterMerchantRoutes(router *gin.RouterGroup) {
	h := NewMerchantHandler()

	merchants := router.Group("/merchants")
	{
		merchants.POST("/onboard", h.OnboardMerchant)
		merchants.GET("/:id", h.GetMerchant)
	}
} 