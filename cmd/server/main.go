package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/singh-anurag-7991/shield/internal/middleware"
	"github.com/singh-anurag-7991/shield/internal/models"
	"github.com/singh-anurag-7991/shield/internal/storage"
)

func main() {
	r := gin.Default()

	memStorage := storage.NewMemoryStorage()

	configs := []models.LimiterConfig{
		{
			Name:      "global",
			Algorithm: "token",
			Capacity:  10,
			Rate:      10,
		},
		{
			Name:      "burst",
			Algorithm: "leaky",
			Capacity:  5,
			Rate:      2,
		},
		// {
		//     Name:      "windowed",
		//     Algorithm: "sliding",
		//     Capacity:  100,
		//     Window:    "60s",
		// },
	}

	r.Use(middleware.RateLimit(memStorage, configs))

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
