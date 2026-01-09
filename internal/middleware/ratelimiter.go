package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"task-manager/internal/utils"
)

type IPTokenBucket struct {
	tokens         float64
	lastRefillTime time.Time
}

type RateLimiter struct {
	mu          sync.Mutex
	ips         map[string]*IPTokenBucket
	rate        float64
	maxCapacity float64
}

var ratelimitRoutes = map[string]bool{
	"/api/auth/login":               true,
	"/api/auth/forgot-password":     true,
	"/api/auth/reset-password":      true,
	"/api/auth/verify-email":        true,
	"/api/auth/resend-verification": true,
}

func NewRateLimiter(limit int, ratePerSecond float64) *RateLimiter {
	rl := &RateLimiter{
		ips:         make(map[string]*IPTokenBucket),
		rate:        ratePerSecond,
		maxCapacity: float64(limit),
	}

	return rl
}

func (rl *RateLimiter) LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !ratelimitRoutes[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			utils.Internal("Internal server error", nil)
			return
		}

		rl.mu.Lock()
		defer rl.mu.Unlock()

		bucket, exists := rl.ips[ip]
		if !exists {
			bucket = &IPTokenBucket{
				tokens:         rl.maxCapacity,
				lastRefillTime: time.Now(),
			}

			rl.ips[ip] = bucket
		}

		now := time.Now()
		elapsed := now.Sub(bucket.lastRefillTime).Seconds()
		bucket.tokens = min(rl.maxCapacity, bucket.tokens+(elapsed*rl.rate))
		bucket.lastRefillTime = now

		if bucket.tokens >= 1 {
			bucket.tokens--
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Too Many requests", http.StatusTooManyRequests)
	})
}
