package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

// API base URL and your API key
const baseURL = "https://api.exchangerate-api.com/v4/latest/"

var apiKey = os.Getenv("5a52c75e0188cd555b782264") // Set your API key as an environment variable

// ExchangeRateResponse struct to parse the API response
type ExchangeRateResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

func fetchExchangeRate(fromCurrency, toCurrency string) (float64, error) {
	url := fmt.Sprintf("%s%s", baseURL, fromCurrency)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var rateResponse ExchangeRateResponse
	err = json.Unmarshal(body, &rateResponse)
	if err != nil {
		return 0, err
	}

	rate, exists := rateResponse.Rates[toCurrency]
	if !exists {
		return 0, fmt.Errorf("rate not found for currency: %s", toCurrency)
	}

	return rate, nil
}

func exchangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		fromCurrency := r.FormValue("fromCurrency")
		toCurrency := r.FormValue("toCurrency")
		amount := r.FormValue("amount")

		// Convert amount to float64 for calculation
		var amountFloat float64
		fmt.Sscanf(amount, "%f", &amountFloat)

		rate, err := fetchExchangeRate(fromCurrency, toCurrency)
		if err != nil {
			http.Error(w, "Error fetching exchange rate", http.StatusInternalServerError)
			return
		}

		exchangedAmount := amountFloat * rate

		fmt.Fprintf(w, "%s %s is approximately %.2f %s", amount, fromCurrency, exchangedAmount, toCurrency)
		return
	}

	// If not POST, render the form
	tmpl, err := template.ParseFiles("index.html") // Ensure you have an "index.html" template
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func main() {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Handle exchange requests
	http.HandleFunc("/exchange", exchangeHandler)

	// Start the server
	http.ListenAndServe(":8080", nil)
}
