package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yourusername/fortexa/api-gateway/internal/config"
	"github.com/yourusername/fortexa/api-gateway/internal/handlers"
	"github.com/yourusername/fortexa/api-gateway/internal/middleware"
)

// @title Fortexa Payment API
// @version 1.0
// @description A Distributed, Fault-Tolerant, Event-Driven Transaction Processing System
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.fortexa.io/support
// @contact.email support@fortexa.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8000
// @BasePath /api/v1
// @schemes http https
func main() {
	// Load configuration
	cfg := config.New()

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create a new Gin router
	router := gin.Default()

	// Setup CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-API-Key")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Setup rate limiting middleware
	rateLimiter := middleware.NewRateLimiter(time.Minute)
	router.Use(rateLimiter.RateLimit())

	// Create a Kafka writer
	kafkaWriter := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Brokers...),
		Topic:        cfg.Kafka.PaymentsTopic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}
	defer kafkaWriter.Close()

	// Create authentication middleware
	authMiddleware := middleware.NewAuthMiddleware()

	// Setup API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		handlers.RegisterMerchantRoutes(v1)
		
		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			handlers.RegisterPaymentRoutes(protected, kafkaWriter)
			handlers.RegisterWebhookRoutes(protected)
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	// Create an HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("API Gateway server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
} 