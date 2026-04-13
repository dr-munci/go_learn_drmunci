package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	if limiter, exists := limiters[ip]; exists {
		return limiter
	}
	limiter := rate.NewLimiter(5, 10) // 5 istek/sn, burst 10
	limiters[ip] = limiter
	return limiter
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !getLimiter(c.ClientIP()).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Çok fazla istek gönderdiniz, lütfen bekleyin",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
