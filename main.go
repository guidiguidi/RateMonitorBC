package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guidiguidi/RateMonitorBC/internal/bestchange"
	"github.com/guidiguidi/RateMonitorBC/config"
	"github.com/guidiguidi/RateMonitorBC/internal/httpapi"
)

func main() {
	// Initialize BestChange service
	client := bestchange.NewClient(
        config.Cfg.BestChange.APIKey,
        config.Cfg.BestChange.BaseURL,
        config.Cfg.BestChange.RateLimit,
    )
    
    service, err := bestchange.NewService(client, "data/currencies.json")
    if err != nil {
        log.Fatalf("failed to create service: %v", err)
    }

	// Initialize Gin router
	router := gin.Default()
	handler := httpapi.NewHandler(service)
	api := router.Group("/api/v1")
	{
		api.POST("/best-rate", handler.GetBestRate)
	}
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Run the HTTP server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("listen: %s\n", err)
	}
}
