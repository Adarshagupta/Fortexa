package analyzer

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/fortexa/fraud-detection/internal/models"
)

// FraudAnalyzer analyzes payments for potential fraud
type FraudAnalyzer struct {
	// In a real implementation, this would contain sophisticated fraud detection algorithms,
	// machine learning models, and connections to fraud databases
	fraudThreshold float64
}

// NewFraudAnalyzer creates a new FraudAnalyzer
func NewFraudAnalyzer(fraudThreshold float64) *FraudAnalyzer {
	return &FraudAnalyzer{
		fraudThreshold: fraudThreshold,
	}
}

// AnalyzePayment checks a payment for potential fraud
func (a *FraudAnalyzer) AnalyzePayment(payment models.Payment) models.FraudCheck {
	log.Printf("Analyzing payment for fraud: %s", payment.ID)

	// In a real implementation, this would run multiple sophisticated checks
	// For the MVP, we'll implement some basic checks
	checks := []models.FraudCheckItem{
		a.checkAmount(payment),
		a.checkVelocity(payment),
		a.checkGeolocation(payment),
		a.checkPaymentMethod(payment),
	}

	// Calculate overall risk score (average of all checks)
	var totalScore float64
	for _, check := range checks {
		totalScore += check.Score
	}
	riskScore := totalScore / float64(len(checks))

	// Determine if payment is fraudulent based on threshold
	isFraudulent := riskScore > a.fraudThreshold
	reason := ""
	if isFraudulent {
		reason = "Multiple risk factors detected"
	}

	return models.FraudCheck{
		PaymentID:    payment.ID,
		MerchantID:   payment.MerchantID,
		CustomerID:   payment.CustomerID,
		RiskScore:    riskScore,
		IsFraudulent: isFraudulent,
		Reason:       reason,
		Checks:       checks,
		CreatedAt:    time.Now(),
	}
}

// checkAmount analyzes the payment amount for unusual patterns
func (a *FraudAnalyzer) checkAmount(payment models.Payment) models.FraudCheckItem {
	// In a real implementation, this would check if the amount is typical for the customer,
	// merchant, and payment method
	
	var score float64
	var info string
	
	// Simplified check: Transactions > 10000 are considered higher risk
	if payment.Amount > 10000 {
		score = 0.8
		info = "Unusually large transaction amount"
	} else if payment.Amount > 5000 {
		score = 0.5
		info = "Larger than average transaction amount"
	} else {
		score = 0.1
		info = "Normal transaction amount"
	}
	
	return models.FraudCheckItem{
		Type:  "amount_check",
		Score: score,
		Info:  info,
	}
}

// checkVelocity checks for unusual transaction frequency
func (a *FraudAnalyzer) checkVelocity(payment models.Payment) models.FraudCheckItem {
	// In a real implementation, this would check the number of transactions
	// from this customer in the last hour/day
	
	// For the MVP, we'll simulate with random data
	score := rand.Float64() * 0.5
	
	var info string
	if score > 0.3 {
		info = "Multiple transactions detected in a short period"
	} else {
		info = "Normal transaction frequency"
	}
	
	return models.FraudCheckItem{
		Type:  "velocity_check",
		Score: score,
		Info:  info,
	}
}

// checkGeolocation checks for unusual location patterns
func (a *FraudAnalyzer) checkGeolocation(payment models.Payment) models.FraudCheckItem {
	// In a real implementation, this would check if the transaction location
	// is consistent with the customer's typical locations
	
	// For the MVP, we'll check if the payment has location data in metadata
	var score float64
	var info string
	
	if payment.Metadata != nil {
		if location, ok := payment.Metadata["location"].(string); ok {
			if strings.Contains(strings.ToLower(location), "nigeria") || 
			   strings.Contains(strings.ToLower(location), "ukraine") {
				// Some regions have higher fraud rates (this is just for demo purposes)
				score = 0.9
				info = "Transaction from high-risk region"
			} else {
				score = 0.2
				info = "Transaction from normal region"
			}
		} else {
			score = 0.5
			info = "No location data provided"
		}
	} else {
		score = 0.5
		info = "No location data provided"
	}
	
	return models.FraudCheckItem{
		Type:  "geolocation_check",
		Score: score,
		Info:  info,
	}
}

// checkPaymentMethod analyzes the payment method for risk
func (a *FraudAnalyzer) checkPaymentMethod(payment models.Payment) models.FraudCheckItem {
	// Different payment methods have different risk profiles
	
	var score float64
	var info string
	
	switch payment.PaymentMethodType {
	case models.PaymentMethodCreditCard, models.PaymentMethodDebitCard:
		// Card payments have moderate risk
		score = 0.4
		info = "Card payment - moderate risk"
	case models.PaymentMethodUPI:
		// UPI is generally linked to bank accounts and has lower risk
		score = 0.2
		info = "UPI payment - lower risk"
	case models.PaymentMethodBankTransfer:
		// Bank transfers have lower risk
		score = 0.1
		info = "Bank transfer - lower risk"
	case models.PaymentMethodCrypto:
		// Crypto has higher risk due to anonymity
		score = 0.7
		info = "Crypto payment - higher risk"
	default:
		score = 0.5
		info = "Unknown payment method"
	}
	
	return models.FraudCheckItem{
		Type:  "payment_method_check",
		Score: score,
		Info:  info,
	}
} 