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

// PaymentEvent represents a payment event received from or published to Kafka
type PaymentEvent struct {
	ID        uuid.UUID     `json:"id"`
	Type      string        `json:"type"`
	Payment   Payment       `json:"payment"`
	Timestamp time.Time     `json:"timestamp"`
}

// FraudCheck represents a fraud check result
type FraudCheck struct {
	PaymentID   uuid.UUID      `json:"payment_id"`
	MerchantID  uuid.UUID      `json:"merchant_id"`
	CustomerID  uuid.UUID      `json:"customer_id,omitempty"`
	RiskScore   float64        `json:"risk_score"`
	IsFraudulent bool          `json:"is_fraudulent"`
	Reason      string         `json:"reason,omitempty"`
	Checks      []FraudCheckItem `json:"checks"`
	CreatedAt   time.Time      `json:"created_at"`
}

// FraudCheckItem represents an individual fraud check item
type FraudCheckItem struct {
	Type  string  `json:"type"`
	Score float64 `json:"score"`
	Info  string  `json:"info,omitempty"`
}

// FraudEvent represents a fraud check event to be published to Kafka
type FraudEvent struct {
	ID        uuid.UUID  `json:"id"`
	Type      string     `json:"type"`
	FraudCheck FraudCheck `json:"fraud_check"`
	Timestamp time.Time  `json:"timestamp"`
} 