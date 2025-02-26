package models

import (
	"time"

	"github.com/google/uuid"
)

// Payment statuses
const (
	PaymentStatusPending   = "PENDING"
	PaymentStatusCaptured  = "CAPTURED"
	PaymentStatusFailed    = "FAILED"
	PaymentStatusRefunded  = "REFUNDED"
	PaymentStatusSettled   = "SETTLED"
)

// SettlementStatus represents the status of a settlement
type SettlementStatus string

// Settlement status constants
const (
	SettlementStatusPending   SettlementStatus = "PENDING"
	SettlementStatusProcessing SettlementStatus = "PROCESSING"
	SettlementStatusCompleted  SettlementStatus = "COMPLETED"
	SettlementStatusFailed     SettlementStatus = "FAILED"
)

// SettlementMethod represents the method used for settlement
type SettlementMethod string

// Settlement method constants
const (
	SettlementMethodBankTransfer SettlementMethod = "BANK_TRANSFER"
	SettlementMethodWallet       SettlementMethod = "WALLET"
)

// Payment represents a payment transaction
type Payment struct {
	ID              uuid.UUID `json:"id"`
	MerchantID      uuid.UUID `json:"merchant_id"`
	OrderID         string    `json:"order_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	PaymentMethod   string    `json:"payment_method"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	SettlementReady bool      `json:"settlement_ready"`
}

// PaymentEvent represents a payment event from Kafka
type PaymentEvent struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Payment   Payment   `json:"payment"`
	Timestamp time.Time `json:"timestamp"`
}

// PaymentSummary represents a summary of payments for a merchant
type PaymentSummary struct {
	MerchantID      uuid.UUID `json:"merchant_id"`
	TotalAmount     float64   `json:"total_amount"`
	Currency        string    `json:"currency"`
	PaymentCount    int       `json:"payment_count"`
	EarliestPayment time.Time `json:"earliest_payment"`
	LatestPayment   time.Time `json:"latest_payment"`
}

// Settlement represents a settlement batch for a merchant
type Settlement struct {
	ID              uuid.UUID        `json:"id"`
	MerchantID      uuid.UUID        `json:"merchant_id"`
	Amount          float64          `json:"amount"`
	Currency        string           `json:"currency"`
	Status          SettlementStatus `json:"status"`
	PaymentCount    int              `json:"payment_count"`
	FeeAmount       float64          `json:"fee_amount"`
	TaxAmount       float64          `json:"tax_amount"`
	NetAmount       float64          `json:"net_amount"`
	SettlementDate  time.Time        `json:"settlement_date"`
	BankAccountID   string           `json:"bank_account_id"`
	SettlementMethod SettlementMethod `json:"settlement_method"`
	Reference       string           `json:"reference"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// SettlementEvent represents a settlement event for Kafka
type SettlementEvent struct {
	ID         uuid.UUID  `json:"id"`
	Settlement Settlement `json:"settlement"`
	Timestamp  time.Time  `json:"timestamp"`
}

// MerchantSettlementConfig represents settlement configuration for a merchant
type MerchantSettlementConfig struct {
	MerchantID              uuid.UUID        `json:"merchant_id"`
	SettlementCycle         string           `json:"settlement_cycle"` // DAILY, WEEKLY, MONTHLY
	PreferredSettlementDay  int              `json:"preferred_settlement_day"` // Day of week/month
	SettlementMethod        SettlementMethod `json:"settlement_method"`
	BankAccountID           string           `json:"bank_account_id"`
	FeePercent              float64          `json:"fee_percent"`
	MinimumSettlementAmount float64          `json:"minimum_settlement_amount"`
	CreatedAt               time.Time        `json:"created_at"`
	UpdatedAt               time.Time        `json:"updated_at"`
} 