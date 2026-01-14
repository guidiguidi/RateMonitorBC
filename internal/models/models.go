package models

// Rate описывает один вариант обмена из /v2/{apiKey}/rates/{from}-{to}.
type Rate struct {
    Changer  int               `json:"changer"`
    Rate     string            `json:"rate"`
    RankRate string            `json:"rankrate"`
    Reserve  string            `json:"reserve"`
    InMin    string            `json:"inmin"`
    InMax    string            `json:"inmax"`
    Marks    []string          `json:"marks"`
    Extra    map[string]any    `json:"extra"`
}

// RatesResponse — ответ от /v2/{apiKey}/rates/{from}-{to}.
type RatesResponse struct {
    Rates map[string][]Rate `json:"rates"` // ключ вида "10-36"
}

// BestRateRequest — вход в твой API.
type BestRateRequest struct {
    FromID int      `json:"from_id"`
    ToID   int      `json:"to_id"`
    Amount float64  `json:"amount"`
    Marks  []string `json:"marks"` // какие метки обязательны (можно пустой список)
}

// BestRateResponse — выход из твоего API.
type BestRateResponse struct {
    Best   Rate   `json:"best"`
    Source string `json:"source"`
}
