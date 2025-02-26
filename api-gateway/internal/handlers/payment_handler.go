package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/yourusername/fortexa/api-gateway/internal/models"
)

// PaymentHandler handles payment-related API endpoints
type PaymentHandler struct {
	kafkaWriter *kafka.Writer
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(kafkaWriter *kafka.Writer) *PaymentHandler {
	return &PaymentHandler{
		kafkaWriter: kafkaWriter,
	}
}

// InitiatePayment handles the payment initiation request
// @Summary Initiate a new payment
// @Description Create a new payment transaction
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body models.PaymentRequest true "Payment Request"
// @Success 200 {object} models.PaymentResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/payments/initiate [post]
func (h *PaymentHandler) InitiatePayment(c *gin.Context) {
	var req models.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a new payment ID
	paymentID := uuid.New()

	// Create payment record
	payment := models.Payment{
		ID:               paymentID,
		MerchantID:       req.MerchantID,
		Amount:           req.Amount,
		Currency:         req.Currency,
		Status:           models.PaymentStatusInitiated,
		PaymentMethodType: req.PaymentMethodType,
		Description:      req.Description,
		Metadata:         req.Metadata,
		IdempotencyKey:   req.IdempotencyKey,
		ReferenceID:      req.ReferenceID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if req.CustomerID != nil {
		payment.CustomerID = *req.CustomerID
	}

	if req.PaymentMethodID != nil {
		payment.PaymentMethodID = req.PaymentMethodID
	}

	// Create a payment event to publish to Kafka
	event := models.PaymentEvent{
		ID:        uuid.New(),
		Type:      "payment.initiated",
		Payment:   payment,
		Timestamp: time.Now(),
	}

	// Serialize the event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize payment event"})
		return
	}

	// Publish the event to Kafka
	err = h.kafkaWriter.WriteMessages(c.Request.Context(), kafka.Message{
		Key:   []byte(payment.ID.String()),
		Value: eventJSON,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish payment event"})
		return
	}

	// Return the payment response
	c.JSON(http.StatusOK, models.PaymentResponse{
		ID:               payment.ID,
		MerchantID:       payment.MerchantID,
		Amount:           payment.Amount,
		Currency:         payment.Currency,
		Status:           payment.Status,
		PaymentMethodType: payment.PaymentMethodType,
		Description:      payment.Description,
		ReferenceID:      payment.ReferenceID,
		CreatedAt:        payment.CreatedAt,
	})
}

// GetPaymentStatus retrieves the status of a payment
// @Summary Get payment status
// @Description Get the current status of a payment
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} models.PaymentResponse
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/payments/{id} [get]
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	// In a real implementation, this would fetch the payment from a database
	// For the MVP, we'll return a mock response
	paymentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	// Mock payment response
	payment := models.PaymentResponse{
		ID:               paymentID,
		MerchantID:       uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Amount:           100.00,
		Currency:         "INR",
		Status:           models.PaymentStatusAuthorized,
		PaymentMethodType: models.PaymentMethodCreditCard,
		CreatedAt:        time.Now().Add(-time.Hour),
	}

	c.JSON(http.StatusOK, payment)
}

// RequestRefund handles payment refund requests
// @Summary Request a refund
// @Description Process a refund for a payment
// @Tags payments
// @Accept json
// @Produce json
// @Param refund body models.RefundRequest true "Refund Request"
// @Success 200 {object} models.RefundResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/refunds [post]
func (h *PaymentHandler) RequestRefund(c *gin.Context) {
	var req models.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a new refund ID
	refundID := uuid.New()

	// Create refund response
	refund := models.RefundResponse{
		ID:        refundID,
		PaymentID: req.PaymentID,
		Amount:    req.Amount,
		Status:    models.PaymentStatusRefunded,
		CreatedAt: time.Now(),
	}

	// In a real implementation, this would:
	// 1. Validate the refund request
	// 2. Check if the payment exists and can be refunded
	// 3. Process the refund with the payment provider
	// 4. Update the payment status
	// 5. Publish an event to Kafka

	c.JSON(http.StatusOK, refund)
}

// RegisterPaymentRoutes registers the payment routes with the given router group
func RegisterPaymentRoutes(router *gin.RouterGroup, kafkaWriter *kafka.Writer) {
	h := NewPaymentHandler(kafkaWriter)

	payments := router.Group("/payments")
	{
		payments.POST("/initiate", h.InitiatePayment)
		payments.GET("/:id", h.GetPaymentStatus)
	}

	refunds := router.Group("/refunds")
	{
		refunds.POST("", h.RequestRefund)
	}
} 