package bestchange

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "github.com/guidiguidi/bestchange/internal/models"
)

type Client struct {
    apiKey string
    http   *http.Client
}

func NewClient(apiKey string) *Client {
    return &Client{
        apiKey: apiKey,
        http:   http.DefaultClient,
    }
}

// hasAllMarks проверяет, что у курса есть все требуемые метки.
func hasAllMarks(rate models.Rate, required []string) bool {
    if len(required) == 0 {
        return true
    }

    set := make(map[string]struct{}, len(rate.Marks))
    for _, m := range rate.Marks {
        set[m] = struct{}{}
    }

    for _, m := range required {
        if _, ok := set[m]; !ok {
            return false
        }
    }
    return true
}

// GetBestRateWithFilters получает курсы для пары и выбирает лучший по фильтрам.
func (c *Client) GetBestRateWithFilters(
    ctx context.Context,
    fromID, toID int,
    amount float64,
    requiredMarks []string,
) (*models.Rate, error) {
    url := fmt.Sprintf("https://bestchange.app/v2/%s/rates/%d-%d", c.apiKey, fromID, toID)

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := c.http.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    var rr models.RatesResponse
    if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
        return nil, err
    }

    key := fmt.Sprintf("%d-%d", fromID, toID)
    rates, ok := rr.Rates[key]
    if !ok || len(rates) == 0 {
        return nil, fmt.Errorf("no rates for pair %s", key)
    }

    // Фильтрация по сумме и меткам.
    filtered := make([]models.Rate, 0, len(rates))
    for _, r := range rates {
        inMin, err := strconv.ParseFloat(r.InMin, 64)
        if err != nil {
            continue
        }
        inMax, err := strconv.ParseFloat(r.InMax, 64)
        if err != nil {
            continue
        }

        if amount < inMin || amount > inMax {
            continue
        }
        if !hasAllMarks(r, requiredMarks) {
            continue
        }

        filtered = append(filtered, r)
    }

    if len(filtered) == 0 {
        return nil, fmt.Errorf("no rates match filters for %s", key)
    }

    // BestChange уже сортирует по rankrate, первый — лучший.
    best := filtered[0]
    return &best, nil
}
