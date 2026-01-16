package bestchange

import (
	"testing"

	"github.com/guidiguidi/RateMonitorBC/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFilterRatesByAmount(t *testing.T) {
	rates := []models.Rate{
		{InMin: "100", InMax: "1000"},
		{InMin: "500", InMax: "2000"},
		{InMin: "1500", InMax: "5000"},
	}

	service := &Service{}

	// Test case 1: Amount within range
	filtered1 := service.filterRates(rates, 600, nil)
	assert.Len(t, filtered1, 2)

	// Test case 2: Amount below all ranges
	filtered2 := service.filterRates(rates, 50, nil)
	assert.Len(t, filtered2, 0)

	// Test case 3: Amount above all ranges
	filtered3 := service.filterRates(rates, 6000, nil)
	assert.Len(t, filtered3, 0)
}

func TestFilterRatesByMarks(t *testing.T) {
	rates := []models.Rate{
		{InMin: "100", InMax: "1000", Marks: []string{"manual", "reg"}},
		{InMin: "100", InMax: "1000", Marks: []string{"manual"}},
		{InMin: "100", InMax: "1000", Marks: []string{"reg"}},
	}

	service := &Service{}

	// Test case 1: Require "manual" and "reg"
	filtered1 := service.filterRates(rates, 200, []string{"manual", "reg"})
	assert.Len(t, filtered1, 1)

	// Test case 2: Require "manual"
	filtered2 := service.filterRates(rates, 200, []string{"manual"})
	assert.Len(t, filtered2, 2)

	// Test case 3: Require "floating" (not present)
	filtered3 := service.filterRates(rates, 200, []string{"floating"})
	assert.Len(t, filtered3, 0)
}

func TestFindBestRate(t *testing.T) {
	rates := []models.Rate{
		{RankRate: "0.00002360", ToAmount: "0.02360"},
		{RankRate: "0.00002345", ToAmount: "0.02345"},
		{RankRate: "0.00002345", ToAmount: "0.02355"},
	}

	service := &Service{}

	bestRate := service.findBestRate(rates)
	assert.Equal(t, "0.00002345", bestRate.RankRate)
	assert.Equal(t, "0.02355", bestRate.ToAmount)
}
