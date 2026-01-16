package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/guidiguidi/RateMonitorBC/config"
	"github.com/guidiguidi/RateMonitorBC/internal/bestchange"
	"github.com/guidiguidi/RateMonitorBC/internal/httpapi"
	"github.com/guidiguidi/RateMonitorBC/internal/models"
	"github.com/spf13/cobra"
)

var (
	fromCode string
	toCode   string
	amount float64
	marks  []string
)

var rootCmd = &cobra.Command{
	Use:   "bestchange",
	Short: "A tool to find the best exchange rates and run a web server.",
}

var bestCmd = &cobra.Command{
	Use:   "best",
	Short: "Find the best exchange rate",
	Run:   findBestRate,
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the web server",
	Run:   runServer,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	bestCmd.Flags().StringVar(&fromCode, "from", "", "From currency code")
	bestCmd.Flags().StringVar(&toCode, "to", "", "To currency code")
	bestCmd.Flags().Float64Var(&amount, "amount", 0, "Amount to exchange")
	bestCmd.Flags().StringSliceVar(&marks, "marks", []string{}, "Required marks (comma-separated)")

	bestCmd.MarkFlagRequired("from")
	bestCmd.MarkFlagRequired("to")
	bestCmd.MarkFlagRequired("amount")

	rootCmd.AddCommand(bestCmd)
	rootCmd.AddCommand(serveCmd)
}

func initConfig() {
	if err := config.Load(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
}

func findBestRate(cmd *cobra.Command, args []string) {
	cfg := config.Cfg

	client := bestchange.NewClient(cfg.BestChange.APIKey, cfg.BestChange.BaseURL, cfg.BestChange.RateLimit)
	service, err := bestchange.NewService(client, "data/currencies.json")
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	req := &models.BestRateRequest{
		FromCode: fromCode,
		ToCode:   toCode,
		Amount: amount,
		Marks:  marks,
	}

	bestRate, err := service.GetBestRate(context.Background(), req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	rate := bestRate.BestRate
	fmt.Println("Best Rate Found:")
	fmt.Printf("  Exchanger ID: %s\n", rate.ExchangerID)
	fmt.Printf("  Rate: %s\n", rate.Rate)
	fmt.Printf("  To Amount: %s\n", rate.ToAmount)
	fmt.Printf("  Marks: %s\n", strings.Join(rate.Marks, ", "))
}

func runServer(cmd *cobra.Command, args []string) {
	cfg := config.Cfg

	bcClient := bestchange.NewClient(cfg.BestChange.APIKey, cfg.BestChange.BaseURL, cfg.BestChange.RateLimit)
	bcService, err := bestchange.NewService(bcClient, "data/currencies.json")
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}
	h := httpapi.NewHandler(bcService)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.Static("/static", "./web/static")
	r.StaticFile("/", "./web/index.html")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/best-rate", h.GetBestRate)
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		log.Printf("🚀 Server starting on :%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
