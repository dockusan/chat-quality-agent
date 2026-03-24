package middleware

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

type visitor struct {
	count    int
	lastSeen time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
	// Cleanup old entries every minute
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[security] panic in rate limiter cleanup goroutine: %v", r)
			}
		}()
		for {
			time.Sleep(time.Minute)
			rl.cleanup()
		}
	}()
	return rl
}

func (rl *rateLimiter) allow(key string) (bool, int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, exists := rl.visitors[key]
	if !exists || now.Sub(v.lastSeen) > rl.window {
		// New visitor or window expired — reset counter
		rl.visitors[key] = &visitor{count: 1, lastSeen: now}
		return true, rl.limit - 1
	}

	v.count++
	remaining := rl.limit - v.count
	if remaining < 0 {
		remaining = 0
	}
	return v.count <= rl.limit, remaining
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for key, v := range rl.visitors {
		if time.Since(v.lastSeen) > rl.window {
			delete(rl.visitors, key)
		}
	}
}

// RateLimit creates a rate limiting middleware.
func RateLimit(limitPerMinute int) gin.HandlerFunc {
	limiter := newRateLimiter(limitPerMinute, time.Minute)
	return func(c *gin.Context) {
		key := c.ClientIP()
		// If authenticated, use user_id for more granular limiting
		if userID := GetUserID(c); userID != "" {
			key = "user:" + userID
		}

		allowed, remaining := limiter.allow(key)
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limitPerMinute))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		if !allowed {
			log.Printf("[security] rate limit exceeded: key=%s ip=%s path=%s", key, c.ClientIP(), c.Request.URL.Path)
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate_limit_exceeded",
			})
			return
		}
		c.Next()
	}
}
