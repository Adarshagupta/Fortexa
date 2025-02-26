package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/adarshagupta/fortexa/settlement-engine/internal/models"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct{}

// NewMockRepository creates a new mock repository for demonstration
func NewMockRepository() Repository {
	log.Println("Using mock repository for database operations")
	return &MockRepository{}
}

// MarkPaymentForSettlement mocks marking a payment for settlement
func (r *MockRepository) MarkPaymentForSettlement(paymentID uuid.UUID) error {
	log.Printf("[MOCK] Marked payment %s for settlement", paymentID)
	return nil
}

// GetEligiblePayments mocks retrieving eligible payments
func (r *MockRepository) GetEligiblePayments(startDate, endDate time.Time) ([]models.PaymentSummary, error) {
	log.Printf("[MOCK] Getting eligible payments between %s and %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	
	// Create some mock payment summaries
	merchantID := uuid.New()
	return []models.PaymentSummary{
		{
			MerchantID:      merchantID,
			TotalAmount:     1000.50,
			Currency:        "USD",
			PaymentCount:    5,
			EarliestPayment: startDate,
			LatestPayment:   endDate,
		},
		{
			MerchantID:      merchantID,
			TotalAmount:     750.25,
			Currency:        "EUR",
			PaymentCount:    3,
			EarliestPayment: startDate,
			LatestPayment:   endDate,
		},
	}, nil
}

// CreateSettlement mocks creating a settlement record
func (r *MockRepository) CreateSettlement(settlement models.Settlement) error {
	log.Printf("[MOCK] Created settlement %s for merchant %s (%s %.2f)", 
		settlement.ID, settlement.MerchantID, settlement.Currency, settlement.Amount)
	return nil
}

// UpdateSettlementStatus mocks updating settlement status
func (r *MockRepository) UpdateSettlementStatus(settlementID uuid.UUID, status models.SettlementStatus) error {
	log.Printf("[MOCK] Updated settlement %s status to %s", settlementID, status)
	return nil
}

// GetMerchantSettlementConfig mocks getting merchant settlement config
func (r *MockRepository) GetMerchantSettlementConfig(merchantID uuid.UUID) (models.MerchantSettlementConfig, error) {
	log.Printf("[MOCK] Retrieved settlement config for merchant %s", merchantID)
	return models.MerchantSettlementConfig{
		MerchantID:              merchantID,
		SettlementCycle:         "DAILY",
		PreferredSettlementDay:  1,
		SettlementMethod:        models.SettlementMethodBankTransfer,
		BankAccountID:           fmt.Sprintf("bank_acc_%s", merchantID.String()[:8]),
		FeePercent:              2.5,
		MinimumSettlementAmount: 100,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}, nil
}

// GetPaymentIDs mocks getting payment IDs for settlement
func (r *MockRepository) GetPaymentIDs(merchantID uuid.UUID, currency string, startDate, endDate time.Time) ([]uuid.UUID, error) {
	log.Printf("[MOCK] Retrieved payment IDs for merchant %s, currency %s between %s and %s", 
		merchantID, currency, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	
	// Generate some mock payment IDs
	ids := make([]uuid.UUID, 5)
	for i := 0; i < 5; i++ {
		ids[i] = uuid.New()
	}
	
	return ids, nil
}