package models

// Currency описывает валюту из /currencies/ru.
type Currency struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// Rate описывает один вариант обмена из /v2/{apiKey}/rates/{from}-{to}.
type Rate struct {
	ExchangerID string   `json:"exchanger_id"`
	Rate        string   `json:"rate"`
	RankRate    string   `json:"rankrate"`
	InMin       string   `json:"inmin"`
	InMax       string   `json:"inmax"`
	Reserve     string   `json:"reserve"`
	Marks       []string `json:"marks"`
	FromAmount  string   `json:"from_amount,omitempty"`
	ToAmount    string   `json:"to_amount,omitempty"`
}

// RatesResponse — ответ от /v2/{apiKey}/rates/{from}-{to}.
type RatesResponse struct {
    Rates map[string][]Rate `json:"rates"`
}

// BestRateRequest - входные параметры API сервиса.
type BestRateRequest struct {
	FromID int      `json:"from_id" binding:"required"`
	ToID   int      `json:"to_id" binding:"required"`
	Amount float64  `json:"amount" binding:"required"`
	Marks  []string `json:"marks"`
}

// BestRateResponse - результат с лучшим курсом.
type BestRateResponse struct {
	FromID    int      `json:"from_id"`
	ToID      int      `json:"to_id"`
	Amount    float64  `json:"amount"`
	Marks     []string `json:"marks"`
	BestRate  Rate     `json:"best_rate"`
	Source    string   `json:"source"`
}
