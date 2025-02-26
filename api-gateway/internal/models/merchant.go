package models

import (
	"time"

	"github.com/google/uuid"
)

// MerchantStatus represents the status of a merchant
type MerchantStatus string

// Merchant statuses
const (
	MerchantStatusActive   MerchantStatus = "ACTIVE"
	MerchantStatusInactive MerchantStatus = "INACTIVE"
	MerchantStatusPending  MerchantStatus = "PENDING"
	MerchantStatusBlocked  MerchantStatus = "BLOCKED"
)

// Merchant represents a merchant in the system
type Merchant struct {
	ID           uuid.UUID      `json:"id"`
	Name         string         `json:"name"`
	BusinessName string         `json:"business_name"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone,omitempty"`
	Website      string         `json:"website,omitempty"`
	APIKey       string         `json:"api_key"`
	SecretKey    string         `json:"-"` // Secret key is never exposed in JSON
	Status       MerchantStatus `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// MerchantRequest represents a request to create or update a merchant
type MerchantRequest struct {
	Name         string `json:"name" binding:"required"`
	BusinessName string `json:"business_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Phone        string `json:"phone" binding:"omitempty"`
	Website      string `json:"website" binding:"omitempty,url"`
}

// MerchantResponse represents a response with merchant details
type MerchantResponse struct {
	ID           uuid.UUID      `json:"id"`
	Name         string         `json:"name"`
	BusinessName string         `json:"business_name"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone,omitempty"`
	Website      string         `json:"website,omitempty"`
	APIKey       string         `json:"api_key"`
	Status       MerchantStatus `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
}

// MerchantEvent represents a merchant event to be published to Kafka
type MerchantEvent struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Merchant  Merchant  `json:"merchant"`
	Timestamp time.Time `json:"timestamp"`
} 