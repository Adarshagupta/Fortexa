package models

import (
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

// Payment statuses
const (
	PaymentStatusInitiated  PaymentStatus = "INITIATED"
	PaymentStatusAuthorized PaymentStatus = "AUTHORIZED"
	PaymentStatusCaptured   PaymentStatus = "CAPTURED"
	PaymentStatusSettled    PaymentStatus = "SETTLED"
	PaymentStatusRefunded   PaymentStatus = "REFUNDED"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusChargeback PaymentStatus = "CHARGEBACK"
)

// PaymentMethod represents the payment method used
type PaymentMethod string

// Payment methods
const (
	PaymentMethodCreditCard   PaymentMethod = "CREDIT_CARD"
	PaymentMethodDebitCard    PaymentMethod = "DEBIT_CARD"
	PaymentMethodUPI          PaymentMethod = "UPI"
	PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
	PaymentMethodWallet       PaymentMethod = "WALLET"
	PaymentMethodCrypto       PaymentMethod = "CRYPTO"
	PaymentMethodBNPL         PaymentMethod = "BNPL"
)

// Payment represents a payment transaction
type Payment struct {
	ID               uuid.UUID      `json:"id"`
	MerchantID       uuid.UUID      `json:"merchant_id"`
	CustomerID       uuid.UUID      `json:"customer_id,omitempty"`
	Amount           float64        `json:"amount"`
	Currency         string         `json:"currency"`
	Status           PaymentStatus  `json:"status"`
	PaymentMethodID  *uuid.UUID     `json:"payment_method_id,omitempty"`
	PaymentMethodType PaymentMethod `json:"payment_method_type"`
	Description      string         `json:"description,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	IdempotencyKey   string         `json:"idempotency_key,omitempty"`
	ReferenceID      string         `json:"reference_id,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// PaymentRequest represents a request to create a new payment
type PaymentRequest struct {
	MerchantID       uuid.UUID      `json:"merchant_id" binding:"required"`
	CustomerID       *uuid.UUID     `json:"customer_id"`
	Amount           float64        `json:"amount" binding:"required,gt=0"`
	Currency         string         `json:"currency" binding:"required,len=3"`
	PaymentMethodID  *uuid.UUID     `json:"payment_method_id"`
	PaymentMethodType PaymentMethod `json:"payment_method_type" binding:"required"`
	Description      string         `json:"description"`
	Metadata         map[string]interface{} `json:"metadata"`
	IdempotencyKey   string         `json:"idempotency_key"`
	ReferenceID      string         `json:"reference_id"`
}

// PaymentResponse represents a response with payment details
type PaymentResponse struct {
	ID               uuid.UUID      `json:"id"`
	MerchantID       uuid.UUID      `json:"merchant_id"`
	CustomerID       *uuid.UUID     `json:"customer_id,omitempty"`
	Amount           float64        `json:"amount"`
	Currency         string         `json:"currency"`
	Status           PaymentStatus  `json:"status"`
	PaymentMethodType PaymentMethod `json:"payment_method_type"`
	Description      string         `json:"description,omitempty"`
	ReferenceID      string         `json:"reference_id,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
}

// PaymentEvent represents a payment event to be published to Kafka
type PaymentEvent struct {
	ID        uuid.UUID     `json:"id"`
	Type      string        `json:"type"`
	Payment   Payment       `json:"payment"`
	Timestamp time.Time     `json:"timestamp"`
}

// RefundRequest represents a request to refund a payment
type RefundRequest struct {
	PaymentID      uuid.UUID `json:"payment_id" binding:"required"`
	Amount         float64   `json:"amount" binding:"required,gt=0"`
	Reason         string    `json:"reason"`
	IdempotencyKey string    `json:"idempotency_key"`
}

// RefundResponse represents a response with refund details
type RefundResponse struct {
	ID        uuid.UUID    `json:"id"`
	PaymentID uuid.UUID    `json:"payment_id"`
	Amount    float64      `json:"amount"`
	Status    PaymentStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
} 