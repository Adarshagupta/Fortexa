package processor

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/adarshagupta/fortexa/settlement-engine/internal/models"
	"github.com/adarshagupta/fortexa/settlement-engine/internal/repository"
)

// SettlementProcessor processes payments and creates settlements
type SettlementProcessor struct {
	repository              repository.Repository
	defaultFeePercent       float64
	minimumSettlementAmount float64
}

// NewSettlementProcessor creates a new settlement processor
func NewSettlementProcessor(
	repository repository.Repository,
	defaultFeePercent float64,
	minimumSettlementAmount float64,
) *SettlementProcessor {
	return &SettlementProcessor{
		repository:              repository,
		defaultFeePercent:       defaultFeePercent,
		minimumSettlementAmount: minimumSettlementAmount,
	}
}

// ProcessPayment processes a payment event and marks it for settlement if eligible
func (p *SettlementProcessor) ProcessPayment(event models.PaymentEvent) error {
	log.Printf("Processing payment event: %s, type: %s", event.ID, event.Type)

	// Check if the payment is captured (eligible for settlement)
	if event.Payment.Status != models.PaymentStatusCaptured {
		log.Printf("Payment %s is not captured, status: %s, skipping", event.Payment.ID, event.Payment.Status)
		return nil
	}

	// Mark the payment for settlement
	if err := p.repository.MarkPaymentForSettlement(event.Payment.ID); err != nil {
		return fmt.Errorf("failed to mark payment for settlement: %w", err)
	}

	log.Printf("Payment %s marked for settlement", event.Payment.ID)
	return nil
}

// CreateSettlementBatch creates settlements for eligible payments
func (p *SettlementProcessor) CreateSettlementBatch(startDate, endDate time.Time) ([]models.Settlement, error) {
	log.Printf("Creating settlement batch for period: %s to %s", 
		startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// Get eligible payments for settlement
	payments, err := p.repository.GetEligiblePayments(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get eligible payments: %w", err)
	}

	log.Printf("Found %d eligible payment groups for settlement", len(payments))
	if len(payments) == 0 {
		return []models.Settlement{}, nil
	}

	var settlements []models.Settlement

	// Process each payment group (merchant + currency)
	for _, paymentSummary := range payments {
		// Get merchant settlement configuration
		config, err := p.repository.GetMerchantSettlementConfig(paymentSummary.MerchantID)
		if err != nil {
			log.Printf("Error getting settlement config for merchant %s: %v, using defaults", 
				paymentSummary.MerchantID, err)
			// Continue with default config
		}

		// Check if amount meets minimum threshold
		if paymentSummary.TotalAmount < config.MinimumSettlementAmount {
			log.Printf("Payment amount %.2f is below minimum settlement amount %.2f for merchant %s, skipping", 
				paymentSummary.TotalAmount, config.MinimumSettlementAmount, paymentSummary.MerchantID)
			continue
		}

		// Get payment IDs for this batch - this is for informational and tracking purposes
		_, err = p.repository.GetPaymentIDs(
			paymentSummary.MerchantID,
			paymentSummary.Currency,
			startDate,
			endDate,
		)
		if err != nil {
			log.Printf("Error getting payment IDs for merchant %s: %v, skipping", 
				paymentSummary.MerchantID, err)
			continue
		}

		// Calculate fees and net amount
		feePercent := config.FeePercent
		if feePercent <= 0 {
			feePercent = p.defaultFeePercent
		}

		feeAmount := roundToTwoDecimals(paymentSummary.TotalAmount * feePercent / 100)
		taxAmount := roundToTwoDecimals(feeAmount * 0.18) // 18% GST on fees
		netAmount := roundToTwoDecimals(paymentSummary.TotalAmount - feeAmount - taxAmount)

		// Create settlement
		settlement := models.Settlement{
			ID:               uuid.New(),
			MerchantID:       paymentSummary.MerchantID,
			Amount:           paymentSummary.TotalAmount,
			Currency:         paymentSummary.Currency,
			Status:           models.SettlementStatusPending,
			PaymentCount:     paymentSummary.PaymentCount,
			FeeAmount:        feeAmount,
			TaxAmount:        taxAmount,
			NetAmount:        netAmount,
			SettlementDate:   time.Now(),
			BankAccountID:    config.BankAccountID,
			SettlementMethod: config.SettlementMethod,
			Reference:        fmt.Sprintf("SET_%s", uuid.New().String()[:8]),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		// Save settlement to database
		if err := p.repository.CreateSettlement(settlement); err != nil {
			log.Printf("Error creating settlement for merchant %s: %v", 
				paymentSummary.MerchantID, err)
			continue
		}

		settlements = append(settlements, settlement)
		log.Printf("Created settlement %s for merchant %s, amount: %.2f %s, net: %.2f", 
			settlement.ID, settlement.MerchantID, settlement.Amount, 
			settlement.Currency, settlement.NetAmount)
	}

	return settlements, nil
}

// ProcessSettlement processes a settlement and completes the funds transfer
func (p *SettlementProcessor) ProcessSettlement(settlement models.Settlement) error {
	log.Printf("Processing settlement %s for merchant %s", settlement.ID, settlement.MerchantID)

	// Update status to processing
	if err := p.repository.UpdateSettlementStatus(settlement.ID, models.SettlementStatusProcessing); err != nil {
		return fmt.Errorf("failed to update settlement status to processing: %w", err)
	}

	// Simulate funds transfer to merchant bank account
	log.Printf("Transferring %.2f %s to bank account %s for merchant %s", 
		settlement.NetAmount, settlement.Currency, settlement.BankAccountID, settlement.MerchantID)
	
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	// Update status to completed
	if err := p.repository.UpdateSettlementStatus(settlement.ID, models.SettlementStatusCompleted); err != nil {
		return fmt.Errorf("failed to update settlement status to completed: %w", err)
	}

	log.Printf("Settlement %s completed", settlement.ID)
	return nil
}

// Helper function to round to two decimal places
func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

// getMerchantFeePercent returns the fee percentage for a merchant
// In a real implementation, this would fetch the merchant's fee configuration from the database
func (p *SettlementProcessor) getMerchantFeePercent(merchantID uuid.UUID) float64 {
	// For the MVP, we'll use the default fee percentage
	return p.defaultFeePercent
}

// getMerchantSettlementMethod returns the preferred settlement method for a merchant
// In a real implementation, this would fetch the merchant's settlement configuration from the database
func (p *SettlementProcessor) getMerchantSettlementMethod(merchantID uuid.UUID) models.SettlementMethod {
	// For the MVP, we'll use bank transfer as the default
	return models.SettlementMethodBankTransfer
}

// getMerchantBankAccountID returns the bank account ID for a merchant
// In a real implementation, this would fetch the merchant's bank details from the database
func (p *SettlementProcessor) getMerchantBankAccountID(merchantID uuid.UUID) string {
	// For the MVP, we'll use a placeholder value
	return "bank_acc_" + merchantID.String()[:8]
} 