package currency

import (
	"encoding/json"
	"net/http"
	"time"
)

type CurrencyResponse struct {
	Success bool                   `json:"success"`
	Base    string                 `json:"base"`
	Date    string                 `json:"date"`
	Rates   map[string]interface{} `json:"rates"`
}

func GetResponse(base string) (*CurrencyResponse, error) {
	urlString := "https://api.exchangerate.host/latest?base=" + base

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	request, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CurrencyResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
