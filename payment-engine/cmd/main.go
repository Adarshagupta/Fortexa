package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/yourusername/fortexa/payment-engine/internal/config"
	"github.com/yourusername/fortexa/payment-engine/internal/handlers"
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

	// Create Kafka writer for publishing events
	kafkaWriter := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Brokers...),
		Topic:        cfg.Kafka.PaymentsTopic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}
	defer kafkaWriter.Close()

	// Create payment handler
	paymentHandler := handlers.NewPaymentHandler(kafkaReader, kafkaWriter)

	// Start the payment handler
	log.Println("Starting payment processing engine")
	err := paymentHandler.Start(ctx)
	if err != nil {
		log.Fatalf("Error starting payment handler: %v", err)
	}

	log.Println("Payment processing engine shutdown complete")
} 