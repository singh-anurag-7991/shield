package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/singh-anurag-7991/shield/internal/limiter"
	"github.com/singh-anurag-7991/shield/internal/models"
	"github.com/singh-anurag-7991/shield/internal/storage"
)

func RateLimit(storage storage.Storage, configs []models.LimiterConfig) gin.HandlerFunc {
	factory := limiter.NewFactory(storage)
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		for _, cfg := range configs {
			limiterInstance := factory.Create(cfg)
			if !limiterInstance.Allow(clientIP) {
				stats := limiterInstance.GetStats(clientIP)
				c.Header("X-RateLimit-Limit", strconv.FormatInt(stats.Limit, 10))
				c.Header("X-RateLimit-Remaining", strconv.FormatInt(stats.Remaining, 10))
				c.Header("X-RateLimit-Reset", strconv.FormatInt(stats.Reset, 10))
				c.JSON(429, gin.H{
					"error":   "rate limit exceeded",
					"limiter": cfg.Name,
					"stats":   stats,
				})
				c.Abort()
				return
			}
			stats := limiterInstance.GetStats(clientIP)
			c.Header("X-RateLimit-Limit", strconv.FormatInt(stats.Limit, 10))
			c.Header("X-RateLimit-Remaining", strconv.FormatInt(stats.Remaining, 10))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(stats.Reset, 10))
		}
		c.Next()
	}
}
