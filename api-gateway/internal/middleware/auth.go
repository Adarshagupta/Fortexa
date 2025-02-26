package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authentication errors
var (
	ErrMissingAPIKey = errors.New("missing API key")
	ErrInvalidAPIKey = errors.New("invalid API key")
)

// AuthMiddleware middleware for API authentication
type AuthMiddleware struct {
	// In a real implementation, this would be replaced with a database lookup
	// or a caching layer to validate API keys against merchant records
	apiKeys map[string]string // map of API key to merchant ID
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware() *AuthMiddleware {
	// For MVP, we'll use a simple in-memory map of API keys
	// In production, this would be fetched from the database
	return &AuthMiddleware{
		apiKeys: map[string]string{
			"pk_test_123456": "test-merchant-id",
			// Add more test keys as needed
		},
	}
}

// Authenticate validates the API key in the request
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// Also check Authorization header with Bearer scheme
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrMissingAPIKey.Error(),
			})
			c.Abort()
			return
		}

		merchantID, exists := m.apiKeys[apiKey]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrInvalidAPIKey.Error(),
			})
			c.Abort()
			return
		}

		// Store merchant ID in context for later use
		c.Set("merchantID", merchantID)
		c.Next()
	}
}

// RequireAuth is a convenience wrapper for Authenticate
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return m.Authenticate()
} 