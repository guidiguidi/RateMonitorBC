package bestchange

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/guidiguidi/RateMonitorBC/internal/models"
	"golang.org/x/time/rate"
)

type Client struct {
	apiKey      string
	baseURL     string
	httpClient  *http.Client
	rateLimiter *rate.Limiter
}

func NewClient(apiKey, baseURL string, rateLimit int) *Client {
	return &Client{
		apiKey:      apiKey,
		baseURL:     baseURL,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		rateLimiter: rate.NewLimiter(rate.Limit(rateLimit), 1),
	}
}

func (c *Client) GetRates(ctx context.Context, fromID, toID int) (*models.RatesResponse, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter wait: %w", err)
	}

	url := fmt.Sprintf("%s/%s/rates/%d-%d", c.baseURL, c.apiKey, fromID, toID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip reader: %w", err)
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	var ratesResponse models.RatesResponse
	if err := json.NewDecoder(reader).Decode(&ratesResponse); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	return &ratesResponse, nil
}

