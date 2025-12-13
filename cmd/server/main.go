package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/singh-anurag-7991/shield/internal/db"
	"github.com/singh-anurag-7991/shield/internal/middleware"
	"github.com/singh-anurag-7991/shield/internal/models"
	"github.com/singh-anurag-7991/shield/internal/rate"
	"github.com/singh-anurag-7991/shield/internal/storage"
)

func main() {
	r := gin.Default()

	// Storage (Memory or Redis)
	var storageImpl rate.Storage
	storageImpl = storage.NewMemoryStorage()

	// Postgres for configs
	var configs []models.LimiterConfig
	configStore, err := db.NewPostgresConfigStore("postgres://anuragsingh@localhost:5432/shield?sslmode=disable")
	if err == nil {
		if err := configStore.InitTable(context.Background()); err != nil {
			log.Println("Table init warning:", err)
		}
		configs, err = configStore.LoadConfigs(context.Background())
		if err != nil {
			log.Println("Load configs warning:", err)
		}
	}

	// Fallback to hard-coded configs if DB fails
	if len(configs) == 0 {
		configs = []models.LimiterConfig{
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
		}
		log.Println("Using hard-coded configs (DB not available)")
	}

	// Unprotected routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.File("internal/dashboard/index.html")
	})

	// Protected routes with rate limit
	protected := r.Group("/api")
	protected.Use(middleware.RateLimit(storageImpl, configs))
	{
		protected.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message":   "success",
				"timestamp": time.Now(),
			})
		})
	}

	log.Println("🚀 Shield Rate Limiter running on :8080")
	log.Println("Test: curl -N http://localhost:8080/api/test")
	log.Println("Health: curl -N http://localhost:8080/health (no limit)")
	log.Println("Dashboard: http://localhost:8080/dashboard")
	log.Println("💡 Should 429 after ~10 /api/test, /health always 200!")
	r.Run(":8080")
}
