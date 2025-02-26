package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/adarshagupta/fortexa/settlement-engine/internal/config"
	"github.com/adarshagupta/fortexa/settlement-engine/internal/handler"
	"github.com/adarshagupta/fortexa/settlement-engine/internal/processor"
	"github.com/adarshagupta/fortexa/settlement-engine/internal/repository"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize configuration
	cfg := config.LoadConfig()

	// Setup context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signal handling for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		log.Println("Received shutdown signal, initiating graceful shutdown")
		cancel()
	}()

	var repo repository.Repository

	// Check if running in mock mode
	mockMode := strings.ToLower(os.Getenv("MOCK_MODE")) == "true"

	if mockMode {
		// Use mock repository
		log.Println("Running in MOCK MODE - No database connection required")
		repo = repository.NewMockRepository()
	} else {
		// Use real database repository
		connectionString := os.Getenv("DB_CONNECTION_STRING")
		if connectionString == "" {
			// Construct connection string from individual parameters
			connectionString = fmt.Sprintf(
				"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
				cfg.Database.User,
				cfg.Database.Password,
				cfg.Database.Host,
				cfg.Database.Port,
				cfg.Database.Name,
				cfg.Database.SSLMode,
			)
		}

		log.Printf("Connecting to database with connection string: %s", connectionString)
		dbRepo, err := repository.NewDBRepository(connectionString)
		if err != nil {
			log.Printf("Failed to initialize database repository: %v", err)
			log.Println("Falling back to MOCK MODE")
			mockMode = true
			repo = repository.NewMockRepository()
		} else {
			log.Println("Successfully connected to database")
			repo = dbRepo
		}
	}

	// Create Kafka reader for payment events
	paymentReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{cfg.Kafka.Broker},
		Topic:       cfg.Kafka.PaymentTopic,
		GroupID:     cfg.Kafka.ConsumerGroup,
		StartOffset: kafka.LastOffset,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
	})

	// Create Kafka writer for settlement events
	settlementWriter := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.Broker),
		Topic:    cfg.Kafka.SettlementTopic,
		Balancer: &kafka.LeastBytes{},
	}

	// Create settlement processor with repository
	settlementProcessor := processor.NewSettlementProcessor(
		repo,
		cfg.Settlement.DefaultFeePercent,
		cfg.Settlement.MinimumSettlementAmount,
	)

	// Create and start settlement handler
	settlementHandler := handler.NewSettlementHandler(
		ctx,
		paymentReader,
		settlementWriter,
		settlementProcessor,
	)

	log.Printf("Starting Settlement Engine service in %s mode", map[bool]string{true: "MOCK", false: "DATABASE"}[mockMode])
	if err := settlementHandler.Start(); err != nil {
		log.Fatalf("Failed to start settlement handler: %v", err)
	}
} 