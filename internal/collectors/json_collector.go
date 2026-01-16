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

func GetCurrencies(filePath string) ([]models.Currency, error) {
	var currencies []models.Currency
	if err := LoadJSONFromFile(filePath, &currencies); err != nil {
		return nil, err
	}
	return currencies, nil
}

func FindByCode(currencies []models.Currency, code string) *models.Currency {
	for i := range currencies {
		if currencies[i].Code == code {
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
