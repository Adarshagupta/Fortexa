package models

import (
	"time"

	"github.com/google/uuid"
)

// WebhookStatus represents the status of a webhook
type WebhookStatus string

// Webhook statuses
const (
	WebhookStatusActive   WebhookStatus = "ACTIVE"
	WebhookStatusInactive WebhookStatus = "INACTIVE"
)

// WebhookEventStatus represents the status of a webhook event delivery
type WebhookEventStatus string

// Webhook event statuses
const (
	WebhookEventStatusPending  WebhookEventStatus = "PENDING"
	WebhookEventStatusDelivered WebhookEventStatus = "DELIVERED"
	WebhookEventStatusFailed    WebhookEventStatus = "FAILED"
	WebhookEventStatusRetrying  WebhookEventStatus = "RETRYING"
)

// EventType represents the type of event that can trigger a webhook
type EventType string

// Event types
const (
	EventTypePaymentInitiated  EventType = "payment.initiated"
	EventTypePaymentAuthorized EventType = "payment.authorized"
	EventTypePaymentCaptured   EventType = "payment.captured"
	EventTypePaymentFailed     EventType = "payment.failed"
	EventTypePaymentRefunded   EventType = "payment.refunded"
	EventTypeSettlementCreated EventType = "settlement.created"
	EventTypeSettlementPaid    EventType = "settlement.paid"
)

// Webhook represents a webhook configuration for a merchant
type Webhook struct {
	ID         uuid.UUID    `json:"id"`
	MerchantID uuid.UUID    `json:"merchant_id"`
	URL        string       `json:"url"`
	EventTypes []string     `json:"event_types"`
	Status     WebhookStatus `json:"status"`
	Secret     string       `json:"-"` // Secret is never exposed in JSON
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

// WebhookRequest represents a request to create or update a webhook
type WebhookRequest struct {
	URL        string   `json:"url" binding:"required,url"`
	EventTypes []string `json:"event_types" binding:"required"`
}

// WebhookResponse represents a response with webhook details
type WebhookResponse struct {
	ID         uuid.UUID    `json:"id"`
	MerchantID uuid.UUID    `json:"merchant_id"`
	URL        string       `json:"url"`
	EventTypes []string     `json:"event_types"`
	Status     WebhookStatus `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
}

// WebhookEvent represents an event to be sent to a webhook
type WebhookEvent struct {
	ID         uuid.UUID          `json:"id"`
	WebhookID  uuid.UUID          `json:"webhook_id"`
	EventType  EventType          `json:"event_type"`
	Payload    map[string]interface{} `json:"payload"`
	Status     WebhookEventStatus `json:"status"`
	Attempts   int                `json:"attempts"`
	LastAttempt *time.Time        `json:"last_attempt,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
} 