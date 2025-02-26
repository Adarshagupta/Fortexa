package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter is a simple in-memory rate limiter
type RateLimiter struct {
	limits     map[string]int
	windowSize time.Duration
	mutex      sync.RWMutex
	requests   map[string][]time.Time
}

// NewRateLimiter creates a new RateLimiter with the specified window size
func NewRateLimiter(windowSize time.Duration) *RateLimiter {
	return &RateLimiter{
		limits: map[string]int{
			"default": 100, // Default rate limit: 100 requests per window
		},
		windowSize: windowSize,
		requests:   make(map[string][]time.Time),
	}
}

// cleanupOldRequests removes requests that are outside the current time window
func (rl *RateLimiter) cleanupOldRequests(key string, now time.Time) {
	cutoff := now.Add(-rl.windowSize)
	
	if times, exists := rl.requests[key]; exists {
		newTimes := make([]time.Time, 0, len(times))
		for _, t := range times {
			if t.After(cutoff) {
				newTimes = append(newTimes, t)
			}
		}
		rl.requests[key] = newTimes
	}
}

// increment adds a new request timestamp for the given key
func (rl *RateLimiter) increment(key string, now time.Time) int {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Clean up old requests first
	rl.cleanupOldRequests(key, now)

	// Add the new request
	if _, exists := rl.requests[key]; !exists {
		rl.requests[key] = make([]time.Time, 0, 10)
	}
	rl.requests[key] = append(rl.requests[key], now)

	return len(rl.requests[key])
}

// isOverLimit checks if the number of requests for the key exceeds the limit
func (rl *RateLimiter) isOverLimit(key string, count int) bool {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	limit, exists := rl.limits[key]
	if !exists {
		limit = rl.limits["default"]
	}

	return count > limit
}

// RateLimit returns a Gin middleware function that limits request rates
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Determine the rate limit key (use IP address as a fallback)
		key := c.ClientIP()
		
		// Use API key as the rate limit key if available
		if apiKey := c.GetHeader("X-API-Key"); apiKey != "" {
			key = apiKey
		}

		now := time.Now()
		count := rl.increment(key, now)

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", "100")
		c.Header("X-RateLimit-Remaining", "100")
		c.Header("X-RateLimit-Reset", "60")

		if rl.isOverLimit(key, count) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
} 