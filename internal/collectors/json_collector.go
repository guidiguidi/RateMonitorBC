package collectors

import (
    "encoding/json"
    "io"
    "os"
    "github.com/guidiguidi/RateMonitorBC/internal/models"
)

const CurrenciesFile = "data/currencies.json"  

type Currency struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`   
    URLName  string  `json:"urlname"`
    ViewName string  `json:"viewname"`
    Code     string  `json:"code"`     
    Crypto   bool    `json:"crypto"`
    Cash     bool    `json:"cash"`
    PS       int     `json:"ps"`
    Group    int     `json:"group"`
    DefAmt   int     `json:"defamt"`
    BigAmt   int     `json:"bigamt"`
    Pos      int     `json:"pos"`
    Rank     int     `json:"rank"`
    Keywords string  `json:"keywords"`
}

func LoadJSONFromFile(filePath string, v interface{}) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()
    data, err := io.ReadAll(file)
    if err != nil {
        return err
    }
    return json.Unmarshal(data, v)
}

func GetCurrencies() ([]models.Currency, error) {
    var currencies []models.Currency
    if err := LoadJSONFromFile(CurrenciesFile, &currencies); err != nil {
        return nil, err
    }
    return currencies, nil
}

func FindByCode(currencies []models.Currency, code string) *models.Currency {
    codeMap := make(map[string]*models.Currency)
    for i := range currencies {
        codeMap[currencies[i].Code] = &currencies[i]
    }
    if c, ok := codeMap[code]; ok {
        return c
    }
    return nil
}

