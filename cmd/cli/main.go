package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidiguidi/RateMonitorBC/config"
	"github.com/guidiguidi/RateMonitorBC/internal/bestchange"
	"github.com/guidiguidi/RateMonitorBC/internal/models"
	"github.com/spf13/cobra"
)

var (
	fromID int
	toID   int
	amount float64
	marks  []string
)

var rootCmd = &cobra.Command{
	Use:   "bestchange",
	Short: "Find the best exchange rate",
	Run:   findBestRate,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().IntVar(&fromID, "from", 0, "From currency ID")
	rootCmd.Flags().IntVar(&toID, "to", 0, "To currency ID")
	rootCmd.Flags().Float64Var(&amount, "amount", 0, "Amount to exchange")
	rootCmd.Flags().StringSliceVar(&marks, "marks", []string{}, "Required marks (comma-separated)")

	rootCmd.MarkFlagRequired("from")
	rootCmd.MarkFlagRequired("to")
	rootCmd.MarkFlagRequired("amount")
}

func initConfig() {
	if err := config.Load(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
}

func findBestRate(cmd *cobra.Command, args []string) {
	cfg := config.Cfg

	client := bestchange.NewClient(cfg.BestChange.APIKey, cfg.BestChange.BaseURL, cfg.BestChange.RateLimit)
	service := bestchange.NewService(client)

	req := &models.BestRateRequest{
		FromID: fromID,
		ToID:   toID,
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
