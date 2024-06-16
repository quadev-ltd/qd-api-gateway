package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/quadev-ltd/qd-qpi-gateway/internal/errors"
)

// RateLimiter settings type
type RateLimiter struct {
	rate  rate.Limit
	burst int
	ips   map[string]*rate.Limiter
	mtx   sync.Mutex
}

// NewRateLimiter Return new RateLimiter
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		rate:  r,
		burst: b,
		ips:   make(map[string]*rate.Limiter),
	}
}

// GetLimiter returns rate limiter instance for a given IP address locking instance
func (rateLimitter *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rateLimitter.mtx.Lock()
	defer rateLimitter.mtx.Unlock()

	limiter, exists := rateLimitter.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(rateLimitter.rate, rateLimitter.burst)
		rateLimitter.ips[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware returns the rate limiter middleware
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.GetLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": errors.TooManyRequests,
			})
			return
		}

		c.Next()
	}
}
