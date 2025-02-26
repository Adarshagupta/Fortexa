package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adarshagupta/fortexa/settlement-engine/internal/models"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// DBRepository handles database operations
type DBRepository struct {
	db *sql.DB
}

// NewDBRepository creates a new database repository
func NewDBRepository(connectionString string) (*DBRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &DBRepository{db: db}, nil
}

// Close closes the database connection
func (r *DBRepository) Close() error {
	return r.db.Close()
}

// MarkPaymentForSettlement marks a payment as ready for settlement
func (r *DBRepository) MarkPaymentForSettlement(paymentID uuid.UUID) error {
	query := `
        UPDATE payments 
        SET settlement_ready = true, updated_at = $1 
        WHERE id = $2 AND status = $3
    `
	_, err := r.db.Exec(query, time.Now(), paymentID, models.PaymentStatusCaptured)
	if err != nil {
		return fmt.Errorf("failed to mark payment for settlement: %w", err)
	}
	return nil
}

// GetEligiblePayments gets eligible payments for settlement
func (r *DBRepository) GetEligiblePayments(startDate, endDate time.Time) ([]models.PaymentSummary, error) {
	query := `
        SELECT 
            merchant_id,
            currency,
            SUM(amount) as total_amount,
            COUNT(*) as payment_count,
            MIN(created_at) as earliest_payment,
            MAX(created_at) as latest_payment
        FROM payments
        WHERE 
            settlement_ready = true 
            AND status = $1
            AND created_at BETWEEN $2 AND $3
        GROUP BY merchant_id, currency
    `

	rows, err := r.db.Query(query, models.PaymentStatusCaptured, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query eligible payments: %w", err)
	}
	defer rows.Close()

	var summaries []models.PaymentSummary
	for rows.Next() {
		var summary models.PaymentSummary
		var merchantIDStr string

		err := rows.Scan(
			&merchantIDStr,
			&summary.Currency,
			&summary.TotalAmount,
			&summary.PaymentCount,
			&summary.EarliestPayment,
			&summary.LatestPayment,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment summary row: %w", err)
		}

		merchantID, err := uuid.Parse(merchantIDStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse merchant ID: %w", err)
		}
		summary.MerchantID = merchantID

		summaries = append(summaries, summary)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payment summary rows: %w", err)
	}

	return summaries, nil
}

// CreateSettlement creates a new settlement record
func (r *DBRepository) CreateSettlement(settlement models.Settlement) error {
	query := `
        INSERT INTO settlements (
            id, merchant_id, amount, currency, status, payment_count,
            fee_amount, tax_amount, net_amount, settlement_date,
            bank_account_id, settlement_method, reference, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
        )
    `

	_, err := r.db.Exec(
		query,
		settlement.ID,
		settlement.MerchantID,
		settlement.Amount,
		settlement.Currency,
		settlement.Status,
		settlement.PaymentCount,
		settlement.FeeAmount,
		settlement.TaxAmount,
		settlement.NetAmount,
		settlement.SettlementDate,
		settlement.BankAccountID,
		settlement.SettlementMethod,
		settlement.Reference,
		settlement.CreatedAt,
		settlement.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create settlement: %w", err)
	}

	return nil
}

// UpdateSettlementStatus updates the status of a settlement
func (r *DBRepository) UpdateSettlementStatus(settlementID uuid.UUID, status models.SettlementStatus) error {
	query := `
        UPDATE settlements 
        SET status = $1, updated_at = $2 
        WHERE id = $3
    `
	_, err := r.db.Exec(query, status, time.Now(), settlementID)
	if err != nil {
		return fmt.Errorf("failed to update settlement status: %w", err)
	}
	return nil
}

// GetMerchantSettlementConfig gets the settlement configuration for a merchant
func (r *DBRepository) GetMerchantSettlementConfig(merchantID uuid.UUID) (models.MerchantSettlementConfig, error) {
	// For MVP, we'll just return a default configuration
	// In a real implementation, this would fetch from the database
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

// GetPaymentIDs gets all payment IDs for a given merchant, currency, and time range
func (r *DBRepository) GetPaymentIDs(merchantID uuid.UUID, currency string, startDate, endDate time.Time) ([]uuid.UUID, error) {
	query := `
        SELECT id 
        FROM payments
        WHERE 
            merchant_id = $1 
            AND currency = $2
            AND settlement_ready = true
            AND status = $3
            AND created_at BETWEEN $4 AND $5
    `

	rows, err := r.db.Query(query, merchantID, currency, models.PaymentStatusCaptured, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query payment IDs: %w", err)
	}
	defer rows.Close()

	var paymentIDs []uuid.UUID
	for rows.Next() {
		var idStr string
		if err := rows.Scan(&idStr); err != nil {
			return nil, fmt.Errorf("failed to scan payment ID: %w", err)
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse payment ID: %w", err)
		}

		paymentIDs = append(paymentIDs, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payment ID rows: %w", err)
	}

	return paymentIDs, nil
} 