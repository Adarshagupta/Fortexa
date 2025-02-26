package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/yourusername/fortexa/fraud-detection/internal/analyzer"
	"github.com/yourusername/fortexa/fraud-detection/internal/config"
	"github.com/yourusername/fortexa/fraud-detection/internal/models"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Create context that will be canceled on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	// Create Kafka reader for consuming payment events
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Kafka.Brokers,
		Topic:       cfg.Kafka.PaymentsTopic,
		GroupID:     cfg.Kafka.ConsumerGroup,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.FirstOffset,
		MaxWait:     1 * time.Second,
	})
	defer kafkaReader.Close()

	// Create Kafka writer for publishing fraud events
	kafkaWriter := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Brokers...),
		Topic:        cfg.Kafka.PaymentsTopic, // We'll publish back to the same topic
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}
	defer kafkaWriter.Close()

	// Create fraud analyzer
	fraudAnalyzer := analyzer.NewFraudAnalyzer(cfg.App.FraudThreshold)

	// Start processing payments
	log.Println("Starting fraud detection service")
	err := processPayments(ctx, kafkaReader, kafkaWriter, fraudAnalyzer)
	if err != nil {
		log.Fatalf("Error processing payments: %v", err)
	}

	log.Println("Fraud detection service shutdown complete")
}

// processPayments continuously reads payment events from Kafka and analyzes them for fraud
func processPayments(ctx context.Context, reader *kafka.Reader, writer *kafka.Writer, analyzer *analyzer.FraudAnalyzer) error {
	log.Println("Processing payments for fraud detection")

	for {
		select {
		case <-ctx.Done():
			log.Println("Payment processing shutting down")
			return nil
		default:
			message, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			// Process the message
			go processMessage(ctx, message, writer, analyzer)
		}
	}
}

// processMessage processes a Kafka message containing a payment event
func processMessage(ctx context.Context, message kafka.Message, writer *kafka.Writer, analyzer *analyzer.FraudAnalyzer) {
	log.Printf("Processing message with key: %s", string(message.Key))

	var event models.PaymentEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		log.Printf("Error unmarshaling event: %v", err)
		return
	}

	// Only process payment initiated, authorized, or captured events
	if event.Type != "payment.initiated" && 
	   event.Type != "payment.authorized" && 
	   event.Type != "payment.captured" {
		return
	}

	log.Printf("Analyzing payment for fraud: %s, Event: %s", event.Payment.ID, event.Type)

	// Analyze the payment for fraud
	fraudCheck := analyzer.AnalyzePayment(event.Payment)

	// If fraudulent, publish a fraud event
	if fraudCheck.IsFraudulent {
		log.Printf("FRAUD DETECTED: Payment ID: %s, Risk Score: %.2f, Reason: %s", 
			fraudCheck.PaymentID, fraudCheck.RiskScore, fraudCheck.Reason)
		
		// Create a fraud event
		fraudEvent := models.FraudEvent{
			ID:         uuid.New(),
			Type:       "fraud.detected",
			FraudCheck: fraudCheck,
			Timestamp:  time.Now(),
		}

		// Serialize the event
		eventJSON, err := json.Marshal(fraudEvent)
		if err != nil {
			log.Printf("Error marshaling fraud event: %v", err)
			return
		}

		// Publish the fraud event
		err = writer.WriteMessages(ctx, kafka.Message{
			Key:   message.Key,
			Value: eventJSON,
		})

		if err != nil {
			log.Printf("Error publishing fraud event: %v", err)
			return
		}

		log.Printf("Published fraud event for payment: %s", fraudCheck.PaymentID)
	} else {
		log.Printf("Payment passed fraud checks: %s, Risk Score: %.2f", 
			fraudCheck.PaymentID, fraudCheck.RiskScore)
	}
} 