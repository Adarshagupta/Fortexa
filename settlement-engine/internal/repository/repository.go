package repository

import (
	"time"

	"github.com/adarshagupta/fortexa/settlement-engine/internal/models"
	"github.com/google/uuid"
)

// Repository defines the interface for database operations
type Repository interface {
	// MarkPaymentForSettlement marks a payment as ready for settlement
	MarkPaymentForSettlement(paymentID uuid.UUID) error
	
	// GetEligiblePayments gets eligible payments for settlement
	GetEligiblePayments(startDate, endDate time.Time) ([]models.PaymentSummary, error)
	
	// CreateSettlement creates a new settlement record
	CreateSettlement(settlement models.Settlement) error
	
	// UpdateSettlementStatus updates the status of a settlement
	UpdateSettlementStatus(settlementID uuid.UUID, status models.SettlementStatus) error
	
	// GetMerchantSettlementConfig gets the settlement configuration for a merchant
	GetMerchantSettlementConfig(merchantID uuid.UUID) (models.MerchantSettlementConfig, error)
	
	// GetPaymentIDs gets all payment IDs for a given merchant, currency, and time range
	GetPaymentIDs(merchantID uuid.UUID, currency string, startDate, endDate time.Time) ([]uuid.UUID, error)
}

// Ensure DBRepository implements Repository interface
var _ Repository = (*DBRepository)(nil) 