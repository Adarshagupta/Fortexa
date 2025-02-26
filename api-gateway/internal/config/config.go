package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
	Redis    RedisConfig
}

// ServerConfig holds the configuration for the HTTP server
type ServerConfig struct {
	Port            string
	Mode            string
	AllowedOrigins  []string
	RequestTimeout  int
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

// RedisConfig holds the configuration for Redis
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// New returns a new Config struct
func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8000"),
			Mode:            getEnv("GIN_MODE", "debug"),
			AllowedOrigins:  getEnvAsSlice("ALLOWED_ORIGINS", []string{"*"}),
			RequestTimeout:  getEnvAsInt("REQUEST_TIMEOUT", 30),
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
			ConsumerGroup:   getEnv("KAFKA_CONSUMER_GROUP", "api-gateway"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
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