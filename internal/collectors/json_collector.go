package collectors

import (
	"encoding/json"
	"io"
	"os"

	"github.com/guidiguidi/RateMonitorBC/internal/models"
)

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

type CurrenciesResponse struct {
	Currencies []models.Currency `json:"currencies"`
}

func GetCurrencies(filePath string) ([]models.Currency, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var resp CurrenciesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Currencies, nil
}

func FindByName(currencies []models.Currency, name string) *models.Currency {
	for i := range currencies {
		if currencies[i].Name == name {
			return &currencies[i]
		}
	}
	return nil
}

func FindByID(currencies []models.Currency, id int) *models.Currency {
	for i := range currencies {
		if currencies[i].ID == id {
			return &currencies[i]
		}
	}
	return nil
}
