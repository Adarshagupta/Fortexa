package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the service
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
}

// AppConfig holds the configuration for the application
type AppConfig struct {
	Name            string
	Environment     string
	ShutdownTimeout int
}

// DatabaseConfig holds the configuration for the database
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// KafkaConfig holds the configuration for Kafka
type KafkaConfig struct {
	Brokers         []string
	PaymentsTopic   string
	SettlementTopic string
	FraudTopic      string
	ConsumerGroup   string
}

// New returns a new Config struct
func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		App: AppConfig{
			Name:            getEnv("APP_NAME", "payment-engine"),
			Environment:     getEnv("APP_ENV", "development"),
			ShutdownTimeout: getEnvAsInt("SHUTDOWN_TIMEOUT", 5),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "fortexa"),
			Password: getEnv("DB_PASSWORD", "fortexa123"),
			DBName:   getEnv("DB_NAME", "fortexa"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Kafka: KafkaConfig{
			Brokers:         getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			PaymentsTopic:   getEnv("KAFKA_PAYMENTS_TOPIC", "payments"),
			SettlementTopic: getEnv("KAFKA_SETTLEMENT_TOPIC", "settlements"),
			FraudTopic:      getEnv("KAFKA_FRAUD_TOPIC", "fraud"),
			ConsumerGroup:   getEnv("KAFKA_CONSUMER_GROUP", "payment-engine"),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// Simple helper function to read an environment variable into a string slice or return a default value
func getEnvAsSlice(key string, defaultVal []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return []string{value}
	}
	return defaultVal
}

// Simple helper function to read an environment variable into an integer or return a default value
func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
} 