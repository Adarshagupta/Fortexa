package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/adarshagupta/fortexa/settlement-engine/internal/models"
	"github.com/adarshagupta/fortexa/settlement-engine/internal/processor"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// SettlementHandler handles payment and settlement events
type SettlementHandler struct {
	ctx                 context.Context
	kafkaReader         *kafka.Reader
	kafkaWriter         *kafka.Writer
	settlementProcessor *processor.SettlementProcessor
}

// NewSettlementHandler creates a new settlement handler
func NewSettlementHandler(
	ctx context.Context,
	kafkaReader *kafka.Reader,
	kafkaWriter *kafka.Writer,
	settlementProcessor *processor.SettlementProcessor,
) *SettlementHandler {
	return &SettlementHandler{
		ctx:                 ctx,
		kafkaReader:         kafkaReader,
		kafkaWriter:         kafkaWriter,
		settlementProcessor: settlementProcessor,
	}
}

// Start begins processing payment events and creating settlements
func (h *SettlementHandler) Start() error {
	// Start a goroutine for consuming payment events
	go h.consumePaymentEvents()

	// Start a goroutine for creating settlement batches
	go h.createSettlementBatches()

	// Block until context is done
	<-h.ctx.Done()
	
	log.Println("Shutting down settlement handler")
	
	// Close resources
	if err := h.kafkaReader.Close(); err != nil {
		log.Printf("Error closing Kafka reader: %v", err)
	}
	
	if err := h.kafkaWriter.Close(); err != nil {
		log.Printf("Error closing Kafka writer: %v", err)
	}
	
	return nil
}

// consumePaymentEvents consumes payment events from Kafka
func (h *SettlementHandler) consumePaymentEvents() {
	log.Println("Starting payment event consumer")
	
	for {
		select {
		case <-h.ctx.Done():
			return
		default:
			// Set a timeout for reading messages
			readCtx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
			msg, err := h.kafkaReader.ReadMessage(readCtx)
			cancel()
			
			if err != nil {
				// If context was canceled or timeout, continue
				if err == context.Canceled || err == context.DeadlineExceeded {
					continue
				}
				
				log.Printf("Error reading message from Kafka: %v", err)
				time.Sleep(1 * time.Second) // Avoid tight loop on errors
				continue
			}
			
			log.Printf("Received payment event: %s", string(msg.Value))
			
			// Process payment event
			var paymentEvent models.PaymentEvent
			if err := json.Unmarshal(msg.Value, &paymentEvent); err != nil {
				log.Printf("Error unmarshaling payment event: %v", err)
				continue
			}
			
			// Process the payment through the settlement processor
			if err := h.settlementProcessor.ProcessPayment(paymentEvent); err != nil {
				log.Printf("Error processing payment: %v", err)
				continue
			}
			
			log.Printf("Successfully processed payment %s", paymentEvent.Payment.ID)
		}
	}
}

// createSettlementBatches periodically creates settlement batches
func (h *SettlementHandler) createSettlementBatches() {
	log.Println("Starting settlement batch creator")
	
	// Create a ticker for running settlement batch creation
	ticker := time.NewTicker(24 * time.Hour) // Run once per day
	defer ticker.Stop()
	
	// Run once immediately on startup
	h.runSettlementBatch()
	
	for {
		select {
		case <-h.ctx.Done():
			return
		case <-ticker.C:
			h.runSettlementBatch()
		}
	}
}

// runSettlementBatch runs a single settlement batch creation
func (h *SettlementHandler) runSettlementBatch() {
	log.Println("Creating settlement batch")
	
	// Get current time for batch window
	now := time.Now()
	startDate := now.AddDate(0, 0, -1) // 1 day ago
	endDate := now
	
	// Create settlement batch
	settlements, err := h.settlementProcessor.CreateSettlementBatch(startDate, endDate)
	if err != nil {
		log.Printf("Error creating settlement batch: %v", err)
		return
	}
	
	log.Printf("Created %d settlements", len(settlements))
	
	// Publish settlement events
	for _, settlement := range settlements {
		h.publishSettlementEvent(settlement)
	}
}

// publishSettlementEvent publishes a settlement event to Kafka
func (h *SettlementHandler) publishSettlementEvent(settlement models.Settlement) {
	// Create settlement event
	event := models.SettlementEvent{
		ID:         uuid.New(),
		Settlement: settlement,
		Timestamp:  time.Now(),
	}
	
	// Marshal event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling settlement event: %v", err)
		return
	}
	
	// Create Kafka message
	message := kafka.Message{
		Key:   []byte(settlement.ID.String()),
		Value: eventJSON,
		Time:  time.Now(),
	}
	
	// Write message to Kafka
	if err := h.kafkaWriter.WriteMessages(h.ctx, message); err != nil {
		log.Printf("Error publishing settlement event: %v", err)
		return
	}
	
	log.Printf("Published settlement event for settlement %s", settlement.ID)
	
	// Process the settlement (mark as completed, etc.)
	if err := h.settlementProcessor.ProcessSettlement(settlement); err != nil {
		log.Printf("Error processing settlement: %v", err)
		return
	}
} 