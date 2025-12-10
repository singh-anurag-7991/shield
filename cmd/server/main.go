package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/singh-anurag-7991/shield/internal/middleware"
	"github.com/singh-anurag-7991/shield/internal/models"
	"github.com/singh-anurag-7991/shield/internal/rate" // ← YE IMPORT ADD KARO (for Storage interface)
	"github.com/singh-anurag-7991/shield/internal/storage"
)

func main() {
	r := gin.Default()

	var storageImpl rate.Storage

	storageImpl = storage.NewMemoryStorage()

	redisStorage, err := storage.NewRedisStorage("redis://default:AXOBAAIncDJlMTJiNjc2NTlmNDU0MjI1OThjYjFjYjFlNDZjYThlZHAyMjk1Njk@related-weevil-29569.upstash.io:6379")
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}
	storageImpl = redisStorage

	configs := []models.LimiterConfig{
		// {
		// 	Name:      "global",
		// 	Algorithm: "token",
		// 	Capacity:  10,
		// 	Rate:      10,
		// },
		{
			Name:      "burst",
			Algorithm: "leaky",
			Capacity:  5,
			Rate:      2,
		},
	}

	r.Use(middleware.RateLimit(storageImpl, configs))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "success",
			"timestamp": time.Now(),
		})
	})

	log.Println("🚀 Shield Rate Limiter running on :8080")
	log.Println("Test: curl -N http://localhost:8080/api/test")
	log.Println("💡 Should 429 after ~10 fast requests!")
	r.Run(":8080")
}
