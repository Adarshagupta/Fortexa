package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/yourusername/fortexa/payment-engine/internal/models"
	"github.com/yourusername/fortexa/payment-engine/internal/processors"
)

// PaymentHandler handles payment-related events from Kafka
type PaymentHandler struct {
	kafkaReader *kafka.Reader
	kafkaWriter *kafka.Writer
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(reader *kafka.Reader, writer *kafka.Writer) *PaymentHandler {
	return &PaymentHandler{
		kafkaReader: reader,
		kafkaWriter: writer,
	}
}

// Start begins listening for payment events
func (h *PaymentHandler) Start(ctx context.Context) error {
	log.Println("Payment handler started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Payment handler shutting down")
			return nil
		default:
			message, err := h.kafkaReader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			// Process the message
			go h.processMessage(ctx, message)
		}
	}
}

// processMessage processes a Kafka message containing a payment event
func (h *PaymentHandler) processMessage(ctx context.Context, message kafka.Message) {
	log.Printf("Processing message with key: %s", string(message.Key))

	var event models.PaymentEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		log.Printf("Error unmarshaling event: %v", err)
		return
	}

	log.Printf("Received payment event: %s, Payment ID: %s", event.Type, event.Payment.ID)

	switch event.Type {
	case "payment.initiated":
		h.handlePaymentInitiated(ctx, event)
	case "payment.authorization.requested":
		h.handlePaymentAuthorizationRequested(ctx, event)
	case "payment.capture.requested":
		h.handlePaymentCaptureRequested(ctx, event)
	case "payment.refund.requested":
		h.handlePaymentRefundRequested(ctx, event)
	default:
		log.Printf("Unknown event type: %s", event.Type)
	}
}

// handlePaymentInitiated processes a payment.initiated event
func (h *PaymentHandler) handlePaymentInitiated(ctx context.Context, event models.PaymentEvent) {
	payment := event.Payment

	// Create a new event for authorization
	authEvent := models.PaymentEvent{
		ID:        uuid.New(),
		Type:      "payment.authorization.requested",
		Payment:   payment,
		Timestamp: time.Now(),
	}

	// Publish the authorization event
	h.publishEvent(ctx, payment.ID.String(), authEvent)
}

// handlePaymentAuthorizationRequested processes a payment.authorization.requested event
func (h *PaymentHandler) handlePaymentAuthorizationRequested(ctx context.Context, event models.PaymentEvent) {
	payment := event.Payment

	// Get the appropriate payment processor for the payment method
	processor, err := processors.PaymentProcessorFactory(payment.PaymentMethodType)
	if err != nil {
		log.Printf("Error creating processor: %v", err)
		h.publishFailedEvent(ctx, payment, "payment.authorization.failed", err.Error())
		return
	}

	// Create an authorization request
	// In a real implementation, card/UPI/bank details would be fetched from a secure vault
	// For the MVP, we'll simulate with minimal data
	authReq := models.PaymentAuthorizationRequest{
		PaymentID:         payment.ID,
		Amount:            payment.Amount,
		Currency:          payment.Currency,
		PaymentMethodType: payment.PaymentMethodType,
	}

	// Add payment method details based on type
	switch payment.PaymentMethodType {
	case models.PaymentMethodCreditCard, models.PaymentMethodDebitCard:
		authReq.CardDetails = &models.CardDetails{
			CardNumber:     "4111111111111111", // Test card number
			ExpiryMonth:    "12",
			ExpiryYear:     "25",
			CVV:            "123",
			CardholderName: "Test User",
		}
	case models.PaymentMethodUPI:
		authReq.UPIDetails = &models.UPIDetails{
			UPIID: "testuser@upi",
		}
	case models.PaymentMethodBankTransfer:
		authReq.BankDetails = &models.BankDetails{
			AccountNumber: "1234567890",
			IFSC:          "TEST0001",
			AccountName:   "Test User",
		}
	}

	// Process the authorization
	authRes, err := processor.Authorize(authReq)
	if err != nil || !authRes.Approved {
		errorMsg := "Payment authorization failed"
		if err != nil {
			errorMsg = err.Error()
		} else if authRes.Error != "" {
			errorMsg = authRes.Error
		}
		log.Printf("Authorization failed: %s", errorMsg)
		h.publishFailedEvent(ctx, payment, "payment.authorization.failed", errorMsg)
		return
	}

	// Update payment status to AUTHORIZED
	payment.Status = models.PaymentStatusAuthorized
	payment.UpdatedAt = time.Now()

	// Add authorization details to metadata
	if payment.Metadata == nil {
		payment.Metadata = make(map[string]interface{})
	}
	payment.Metadata["authorization_id"] = authRes.AuthorizationID
	payment.Metadata["processor_id"] = authRes.ProcessorID

	// Create a new event for authorization successful
	authSuccessEvent := models.PaymentEvent{
		ID:        uuid.New(),
		Type:      "payment.authorized",
		Payment:   payment,
		Timestamp: time.Now(),
	}

	// Publish the authorization successful event
	h.publishEvent(ctx, payment.ID.String(), authSuccessEvent)

	// For automatic capture, create a capture request event
	captureEvent := models.PaymentEvent{
		ID:        uuid.New(),
		Type:      "payment.capture.requested",
		Payment:   payment,
		Timestamp: time.Now(),
	}

	// Publish the capture request event
	h.publishEvent(ctx, payment.ID.String(), captureEvent)
}

// handlePaymentCaptureRequested processes a payment.capture.requested event
func (h *PaymentHandler) handlePaymentCaptureRequested(ctx context.Context, event models.PaymentEvent) {
	payment := event.Payment

	// Get the appropriate payment processor for the payment method
	processor, err := processors.PaymentProcessorFactory(payment.PaymentMethodType)
	if err != nil {
		log.Printf("Error creating processor: %v", err)
		h.publishFailedEvent(ctx, payment, "payment.capture.failed", err.Error())
		return
	}

	// Process the capture
	err = processor.Capture(payment.ID, payment.Amount)
	if err != nil {
		log.Printf("Capture failed: %v", err)
		h.publishFailedEvent(ctx, payment, "payment.capture.failed", err.Error())
		return
	}

	// Update payment status to CAPTURED
	payment.Status = models.PaymentStatusCaptured
	payment.UpdatedAt = time.Now()

	// Create a new event for capture successful
	captureSuccessEvent := models.PaymentEvent{
		ID:        uuid.New(),
		Type:      "payment.captured",
		Payment:   payment,
		Timestamp: time.Now(),
	}

	// Publish the capture successful event
	h.publishEvent(ctx, payment.ID.String(), captureSuccessEvent)

	// Create a settlement request event
	settlementEvent := models.PaymentEvent{
		ID:        uuid.New(),
		Type:      "payment.settlement.requested",
		Payment:   payment,
		Timestamp: time.Now(),
	}

	// Publish the settlement request event
	h.publishEvent(ctx, payment.ID.String(), settlementEvent)
}

// handlePaymentRefundRequested processes a payment.refund.requested event
func (h *PaymentHandler) handlePaymentRefundRequested(ctx context.Context, event models.PaymentEvent) {
	payment := event.Payment

	// Get the appropriate payment processor for the payment method
	processor, err := processors.PaymentProcessorFactory(payment.PaymentMethodType)
	if err != nil {
		log.Printf("Error creating processor: %v", err)
		h.publishFailedEvent(ctx, payment, "payment.refund.failed", err.Error())
		return
	}

	// Get refund amount from metadata (in a real implementation, this would be part of the refund request)
	refundAmount := payment.Amount // Default to full refund
	if amount, ok := payment.Metadata["refund_amount"].(float64); ok {
		refundAmount = amount
	}

	// Process the refund
	err = processor.Refund(payment.ID, refundAmount)
	if err != nil {
		log.Printf("Refund failed: %v", err)
		h.publishFailedEvent(ctx, payment, "payment.refund.failed", err.Error())
		return
	}

	// Update payment status to REFUNDED
	payment.Status = models.PaymentStatusRefunded
	payment.UpdatedAt = time.Now()

	// Update refund details in metadata
	if payment.Metadata == nil {
		payment.Metadata = make(map[string]interface{})
	}
	payment.Metadata["refund_amount"] = refundAmount
	payment.Metadata["refund_id"] = uuid.New().String()
	payment.Metadata["refund_time"] = time.Now().Format(time.RFC3339)

	// Create a new event for refund successful
	refundSuccessEvent := models.PaymentEvent{
		ID:        uuid.New(),
		Type:      "payment.refunded",
		Payment:   payment,
		Timestamp: time.Now(),
	}

	// Publish the refund successful event
	h.publishEvent(ctx, payment.ID.String(), refundSuccessEvent)
}

// publishFailedEvent publishes a failure event with the error message
func (h *PaymentHandler) publishFailedEvent(ctx context.Context, payment models.Payment, eventType, errorMessage string) {
	// Update payment status to FAILED
	payment.Status = models.PaymentStatusFailed
	payment.UpdatedAt = time.Now()

	// Add error details to metadata
	if payment.Metadata == nil {
		payment.Metadata = make(map[string]interface{})
	}
	payment.Metadata["error"] = errorMessage
	payment.Metadata["failure_time"] = time.Now().Format(time.RFC3339)

	// Create a failure event
	failureEvent := models.PaymentEvent{
		ID:        uuid.New(),
		Type:      eventType,
		Payment:   payment,
		Timestamp: time.Now(),
	}

	// Publish the failure event
	h.publishEvent(ctx, payment.ID.String(), failureEvent)
}

// publishEvent publishes an event to Kafka
func (h *PaymentHandler) publishEvent(ctx context.Context, key string, event models.PaymentEvent) {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}

	err = h.kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: eventJSON,
	})

	if err != nil {
		log.Printf("Error publishing event: %v", err)
		return
	}

	log.Printf("Published event: %s, Payment ID: %s", event.Type, event.Payment.ID)
} 