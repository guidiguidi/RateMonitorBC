package main

import (
    "fmt"
    "log"
	"sort"
    "github.com/guidiguidi/RateMonitorBC/internal/collectors"
    "github.com/guidiguidi/RateMonitorBC/internal/models"
)

func main() {
    currencies, err := collectors.GetCurrencies()
    if err != nil {
        log.Fatal("Error loading currencies:", err)
    }

    fmt.Printf("📊 Total currencies: %d\n", len(currencies))
    
    // Фильтр крипты
    cryptoCount := 0
    for _, c := range currencies {
        if c.Crypto {
            cryptoCount++
        }
    }
    fmt.Printf("🔥 Crypto currencies: %d\n", cryptoCount)

    // Поиск по code
    usdt := collectors.FindByCode(currencies, "USDTTRC20")
    btc := collectors.FindByCode(currencies, "BTC")
    
    if usdt != nil {
        fmt.Printf("\n💰 USDT TRC20: %s (ID=%d, Rank=%d)\n", usdt.Name, usdt.ID, usdt.Rank)
    }
    if btc != nil {
        fmt.Printf("₿ BTC: %s (ID=%d, Rank=%d)\n", btc.Name, btc.ID, btc.Rank)
    }

    // Топ-5 крипты по rank
    fmt.Println("\n🏆 Top 5 Crypto:")
    cryptoList := []models.Currency{}
    for _, c := range currencies {
        if c.Crypto {
            cryptoList = append(cryptoList, c)
        }
    }
    
    // Сортировка по rank (низкий rank = лучше)
    sort.Slice(cryptoList, func(i, j int) bool {
        return cryptoList[i].Rank < cryptoList[j].Rank
    })
    
    for i, c := range cryptoList[:5] {
        fmt.Printf("#%d %s (%s) rank=%d\n", i+1, c.Name, c.Code, c.Rank)
    }
}
