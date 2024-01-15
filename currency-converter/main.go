package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const exchangeRatesAPI = "https://open.er-api.com/v6/latest"

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
}

func getExchangeRates() (ExchangeRates, error) {
	client := resty.New()
	resp, err := client.R().Get(exchangeRatesAPI)
	if err != nil {
		return ExchangeRates{}, err
	}

	var rates ExchangeRates
	if err := json.Unmarshal(resp.Body(), &rates); err != nil {
		return ExchangeRates{}, err
	}

	return rates, nil
}

func convertCurrency(amount float64, fromCurrency string, toCurrency string) (float64, error) {
	rates, err := getExchangeRates()
	if err != nil {
		return 0, err
	}

	fromRate, fromExists := rates.Rates[fromCurrency]
	toRate, toExists := rates.Rates[toCurrency]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("currency not found")
	}

	convertedAmount := (amount / fromRate) * toRate
	return convertedAmount, nil
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Uso: go run main.go <quantidade> <moeda_origem> <moeda_destino>")
		os.Exit(1)
	}

	amount := os.Args[1]
	fromCurrency := os.Args[2]
	toCurrency := os.Args[3]

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Fatal("Erro ao converter a quantidade para float64")
	}

	convertedAmount, err := convertCurrency(amountFloat, fromCurrency, toCurrency)
	if err != nil {
		log.Fatal("Erro ao converter a moeda:", err)
	}

	fmt.Printf("%.2f %s é equivalente a %.2f %s\n", amountFloat, fromCurrency, convertedAmount, toCurrency)
}
