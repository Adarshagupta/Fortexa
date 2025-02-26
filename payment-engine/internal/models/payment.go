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

// PaymentAuthorizationRequest represents a request to authorize a payment with a payment processor
type PaymentAuthorizationRequest struct {
	PaymentID       uuid.UUID      `json:"payment_id"`
	Amount          float64        `json:"amount"`
	Currency        string         `json:"currency"`
	PaymentMethodType PaymentMethod `json:"payment_method_type"`
	CardDetails     *CardDetails   `json:"card_details,omitempty"`
	UPIDetails      *UPIDetails    `json:"upi_details,omitempty"`
	BankDetails     *BankDetails   `json:"bank_details,omitempty"`
}

// PaymentAuthorizationResponse represents a response from a payment processor
type PaymentAuthorizationResponse struct {
	PaymentID       uuid.UUID     `json:"payment_id"`
	ProcessorID     string        `json:"processor_id"`
	Approved        bool          `json:"approved"`
	AuthorizationID string        `json:"authorization_id,omitempty"`
	Error           string        `json:"error,omitempty"`
	Timestamp       time.Time     `json:"timestamp"`
}

// CardDetails represents credit/debit card details
type CardDetails struct {
	CardNumber     string `json:"card_number"`
	ExpiryMonth    string `json:"expiry_month"`
	ExpiryYear     string `json:"expiry_year"`
	CVV            string `json:"cvv"`
	CardholderName string `json:"cardholder_name"`
}

// UPIDetails represents UPI payment details
type UPIDetails struct {
	UPIID string `json:"upi_id"`
}

// BankDetails represents bank transfer details
type BankDetails struct {
	AccountNumber string `json:"account_number"`
	IFSC          string `json:"ifsc"`
	AccountName   string `json:"account_name"`
} 