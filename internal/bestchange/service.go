package bestchange

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/guidiguidi/RateMonitorBC/internal/models"
)

var (
	ErrNoSuitableRates = fmt.Errorf("no suitable rates found")
)

type Service struct {
	client *Client
}

func NewService(client *Client) *Service {
	return &Service{client: client}
}

func (s *Service) GetBestRate(ctx context.Context, req *models.BestRateRequest) (*models.BestRateResponse, error) {
	ratesResponse, err := s.client.GetRates(ctx, req.FromID, req.ToID)
	if err != nil {
		return nil, fmt.Errorf("get rates: %w", err)
	}

	key := fmt.Sprintf("%d-%d", req.FromID, req.ToID)
	rates, ok := ratesResponse.Rates[key]
	if !ok || len(rates) == 0 {
		return nil, ErrNoSuitableRates
	}

	filteredRates := s.filterRates(rates, req.Amount, req.Marks)
	if len(filteredRates) == 0 {
		return nil, ErrNoSuitableRates
	}

	bestRate := s.findBestRate(filteredRates)

	fromAmount := fmt.Sprintf("%.8f", req.Amount)
	rateValue, _ := strconv.ParseFloat(bestRate.Rate, 64)
	toAmount := fmt.Sprintf("%.8f", req.Amount*rateValue)
	bestRate.FromAmount = fromAmount
	bestRate.ToAmount = toAmount

	return &models.BestRateResponse{
		FromID:   req.FromID,
		ToID:     req.ToID,
		Amount:   req.Amount,
		Marks:    req.Marks,
		BestRate: bestRate,
		Source:   "bestchange",
	}, nil
}

func (s *Service) filterRates(rates []models.Rate, amount float64, marks []string) []models.Rate {
	var filtered []models.Rate
	for _, rate := range rates {
		inMin, err := strconv.ParseFloat(rate.InMin, 64)
		if err != nil {
			continue
		}

		if amount < inMin {
			continue
		}

		if rate.InMax != "0" {
			inMax, err := strconv.ParseFloat(rate.InMax, 64)
			if err != nil {
				continue
			}
			if amount > inMax {
				continue
			}
		}

		if !s.hasAllMarks(rate.Marks, marks) {
			continue
		}

		filtered = append(filtered, rate)
	}
	return filtered
}

func (s *Service) hasAllMarks(rateMarks, requiredMarks []string) bool {
	if len(requiredMarks) == 0 {
		return true
	}

	markSet := make(map[string]struct{})
	for _, m := range rateMarks {
		markSet[m] = struct{}{}
	}

	for _, m := range requiredMarks {
		if _, ok := markSet[m]; !ok {
			return false
		}
	}
	return true
}

func (s *Service) findBestRate(rates []models.Rate) models.Rate {
	sort.Slice(rates, func(i, j int) bool {
		rankRateI, _ := strconv.ParseFloat(rates[i].RankRate, 64)
		rankRateJ, _ := strconv.ParseFloat(rates[j].RankRate, 64)

		if rankRateI != rankRateJ {
			return rankRateI < rankRateJ
		}

		toAmountI, _ := strconv.ParseFloat(rates[i].ToAmount, 64)
		toAmountJ, _ := strconv.ParseFloat(rates[j].ToAmount, 64)
		return toAmountI > toAmountJ
	})

	return rates[0]
}
