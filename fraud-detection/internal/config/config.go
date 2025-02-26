package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the service
type Config struct {
	App   AppConfig
	Kafka KafkaConfig
}

// AppConfig holds the configuration for the application
type AppConfig struct {
	Name            string
	Environment     string
	ShutdownTimeout int
	FraudThreshold  float64
}

// KafkaConfig holds the configuration for Kafka
type KafkaConfig struct {
	Brokers       []string
	PaymentsTopic string
	FraudTopic    string
	ConsumerGroup string
}

// New returns a new Config struct
func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		App: AppConfig{
			Name:            getEnv("APP_NAME", "fraud-detection"),
			Environment:     getEnv("APP_ENV", "development"),
			ShutdownTimeout: getEnvAsInt("SHUTDOWN_TIMEOUT", 5),
			FraudThreshold:  getEnvAsFloat("FRAUD_THRESHOLD", 0.7),
		},
		Kafka: KafkaConfig{
			Brokers:       getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			PaymentsTopic: getEnv("KAFKA_PAYMENTS_TOPIC", "payments"),
			FraudTopic:    getEnv("KAFKA_FRAUD_TOPIC", "fraud"),
			ConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "fraud-detection"),
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

// Simple helper function to read an environment variable into a float or return a default value
func getEnvAsFloat(key string, defaultVal float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultVal
} 