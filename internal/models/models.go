package models

import (
	"encoding/json"
	"strings"
)

// Currency описывает валюту из /currencies/ru.
type Currency struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Code     string  `json:"code"`
	Crypto   bool    `json:"crypto"`
	Cash     bool    `json:"cash"`
	Ps       int     `json:"ps"`
	Group    int     `json:"group"`
	Defamt   float64 `json:"defamt"`
	Bigamt   float64 `json:"bigamt"`
	Pos      int     `json:"pos"`
	Rank     int     `json:"rank"`
	Keywords string  `json:"keywords"`
}

// Marks is a custom type to handle inconsistent API responses for the "marks" field.
type Marks []string

// UnmarshalJSON implements the json.Unmarshaler interface for the Marks type.
// It can handle both an array of strings and an empty object.
func (m *Marks) UnmarshalJSON(data []byte) error {
	var s []string
	if err := json.Unmarshal(data, &s); err == nil {
		*m = s
		return nil
	}

	if strings.TrimSpace(string(data)) == "{}" {
		*m = []string{}
		return nil
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = []string{}
	return nil
}

// Rate описывает один вариант обмена из /v2/{apiKey}/rates/{from}-{to}.
type Rate struct {
	ExchangerID string  `json:"exchanger_id"`
	Rate        string  `json:"rate"`
	RankRate    string  `json:"rankrate"`
	InMin       string  `json:"inmin"`
	InMax       string  `json:"inmax"`
	Reserve     string  `json:"reserve"`
	Marks       Marks   `json:"marks"`
	FromAmount  string  `json:"from_amount,omitempty"`
	ToAmount    string  `json:"to_amount,omitempty"`
}

// RatesResponse — ответ от /v2/{apiKey}/rates/{from}-{to}.
type RatesResponse struct {
	Rates map[string][]Rate `json:"rates"`
}

// BestRateRequest - входные параметры API сервиса.
type BestRateRequest struct {
    FromID int      `json:"from_id" binding:"required"`
    ToID   int      `json:"to_id" binding:"required"`
    Amount float64  `json:"amount" binding:"required,gt=0"`
    Marks  []string `json:"marks"`
}

// BestRateResponse - результат с лучшим курсом.
type BestRateResponse struct {
	FromID   int       `json:"from_id"`
	ToID     int       `json:"to_id"`
	Amount   float64  `json:"amount"`
	Marks    []string  `json:"marks"`
	BestRate *Rate     `json:"best_rate"`
	Source   string    `json:"source"`
}