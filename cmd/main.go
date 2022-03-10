package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Error struct {
	Message string `json:"message"`
}

type CurrencyResponse struct {
	Success bool                   `json:"success"`
	Base    string                 `json:"base"`
	Date    string                 `json:"date"`
	Rates   map[string]interface{} `json:"rates"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		// parse x from query parameters
		x, err := castToFloat64("x", values.Get("x"))
		if err != nil {
			json.NewEncoder(w).Encode(&Error{err.Error()})
			return
		}

		// parse y from query parameters
		y, err := castToFloat64("y", values.Get("y"))
		if err != nil {
			json.NewEncoder(w).Encode(&Error{err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// load env parameters from file
		err = godotenv.Load()
		if err != nil {
			json.NewEncoder(w).Encode(&Error{err.Error()})
			return
		}

		// load x0 parameter from env file
		x0, err := castToFloat64("X0", os.Getenv("X0"))
		if err != nil {
			json.NewEncoder(w).Encode(&Error{err.Error()})
			return
		}

		// load y0 parameter from env file
		y0, err := castToFloat64("Y0", os.Getenv("Y0"))
		if err != nil {
			json.NewEncoder(w).Encode(&Error{err.Error()})
			return
		}

		// load radius parameter from env file
		radius, err := castToFloat64("R", os.Getenv("R"))
		if err != nil {
			json.NewEncoder(w).Encode(&Error{err.Error()})
			return
		}

		// calculate formula
		formula := math.Pow(x-x0, 2)+math.Pow(y-y0, 2) <= math.Pow(radius, 2)
		if formula {
			response, err := getCurrency("USD")
			if err != nil {
				json.NewEncoder(w).Encode(&Error{err.Error()})
				return
			}

			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				json.NewEncoder(w).Encode(&Error{err.Error()})
				return
			}
			return
		}

		response, err := getCurrency("EUR")
		if err != nil {
			json.NewEncoder(w).Encode(&Error{err.Error()})
			return
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			json.NewEncoder(w).Encode(&Error{err.Error()})
			return
		}

	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("server is listening on %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func castToFloat64(name, value string) (float64, error) {
	if len(value) == 0 {
		return 0, fmt.Errorf("the %v arguments must not be empty", name)
	}

	casted, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to case to float64: %w", err)
	}
	return casted, nil
}

func getCurrency(base string) (*CurrencyResponse, error) {
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
