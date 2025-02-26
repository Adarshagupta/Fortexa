package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// Config represents the application configuration
type Config struct {
	App        AppConfig
	Database   DBConfig
	Kafka      KafkaConfig
	Settlement SettlementConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name    string
	Version string
	Env     string
}

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Broker          string
	PaymentTopic    string
	SettlementTopic string
	ConsumerGroup   string
}

// SettlementConfig holds settlement configuration
type SettlementConfig struct {
	DefaultFeePercent        float64
	MinimumSettlementAmount  float64
	DefaultSettlementCycle   string
	PreferredSettlementDay   int
	SettlementBatchTimeStart string
	SettlementBatchTimeEnd   string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() Config {
	return Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "settlement-engine"),
			Version: getEnv("APP_VERSION", "1.0.0"),
			Env:     getEnv("APP_ENV", "development"),
		},
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "fortexa"),
			Password: getEnv("DB_PASSWORD", "fortexa123"),
			Name:     getEnv("DB_NAME", "fortexa"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Kafka: KafkaConfig{
			Broker:          getEnv("KAFKA_BROKER", "localhost:9092"),
			PaymentTopic:    getEnv("KAFKA_PAYMENT_TOPIC", "payment-events"),
			SettlementTopic: getEnv("KAFKA_SETTLEMENT_TOPIC", "settlement-events"),
			ConsumerGroup:   getEnv("KAFKA_CONSUMER_GROUP", "settlement-engine"),
		},
		Settlement: SettlementConfig{
			DefaultFeePercent:        getEnvAsFloat("SETTLEMENT_DEFAULT_FEE_PERCENT", 2.5),
			MinimumSettlementAmount:  getEnvAsFloat("SETTLEMENT_MINIMUM_AMOUNT", 100.0),
			DefaultSettlementCycle:   getEnv("SETTLEMENT_DEFAULT_CYCLE", "DAILY"),
			PreferredSettlementDay:   getEnvAsInt("SETTLEMENT_PREFERRED_DAY", 1),
			SettlementBatchTimeStart: getEnv("SETTLEMENT_BATCH_TIME_START", "00:00"),
			SettlementBatchTimeEnd:   getEnv("SETTLEMENT_BATCH_TIME_END", "23:59"),
		},
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Failed to parse %s as int, using default value %d: %v", key, defaultValue, err)
		return defaultValue
	}
	return value
}

// getEnvAsFloat retrieves an environment variable as a float or returns a default value
func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		log.Printf("Warning: Failed to parse %s as float, using default value %.2f: %v", key, defaultValue, err)
		return defaultValue
	}
	return value
}

// getEnvAsBool retrieves an environment variable as a boolean or returns a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		// Check for common string representations
		valueStr = strings.ToLower(valueStr)
		if valueStr == "true" || valueStr == "yes" || valueStr == "1" {
			return true
		}
		if valueStr == "false" || valueStr == "no" || valueStr == "0" {
			return false
		}
		
		log.Printf("Warning: Failed to parse %s as bool, using default value %t: %v", key, defaultValue, err)
		return defaultValue
	}
	return value
} 