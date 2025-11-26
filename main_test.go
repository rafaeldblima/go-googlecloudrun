package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestIsValidCEP(t *testing.T) {
	tests := []struct {
		cep      string
		expected bool
	}{
		{"01310100", true},   // Valid CEP
		{"12345678", true},   // Valid CEP
		{"1234567", false},   // Too short
		{"123456789", false}, // Too long
		{"abcd1234", false},  // Contains letters
		{"1234-567", false},  // Contains dash
		{"", false},          // Empty
	}

	for _, test := range tests {
		result := isValidCEP(test.cep)
		if result != test.expected {
			t.Errorf("isValidCEP(%s) = %v; expected %v", test.cep, result, test.expected)
		}
	}
}

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		celsius  float64
		expected float64
	}{
		{0, 32},    // Freezing point
		{100, 212}, // Boiling point
		{25, 77},   // Room temperature
		{-40, -40}, // Same in both scales
	}

	for _, test := range tests {
		result := celsiusToFahrenheit(test.celsius)
		if result != test.expected {
			t.Errorf("celsiusToFahrenheit(%.1f) = %.1f; expected %.1f", test.celsius, result, test.expected)
		}
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		celsius  float64
		expected float64
	}{
		{0, 273},   // Freezing point
		{100, 373}, // Boiling point
		{25, 298},  // Room temperature
		{-273, 0},  // Absolute zero
	}

	for _, test := range tests {
		result := celsiusToKelvin(test.celsius)
		if result != test.expected {
			t.Errorf("celsiusToKelvin(%.1f) = %.1f; expected %.1f", test.celsius, result, test.expected)
		}
	}
}

func TestGetWeatherByCEPInvalidFormat(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/weather/{cep}", getWeatherByCEP).Methods("GET")

	// Test invalid CEP format
	req, err := http.NewRequest("GET", "/weather/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}

	var response ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Message != "invalid zipcode" {
		t.Errorf("handler returned wrong message: got %v want %v", response.Message, "invalid zipcode")
	}
}

func TestGetWeatherByCEPValidFormat(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/weather/{cep}", getWeatherByCEP).Methods("GET")

	// Test with a valid CEP format (01310100 - Av. Paulista, São Paulo)
	req, err := http.NewRequest("GET", "/weather/01310100", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Should return 200 OK since we have mock weather data
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response WeatherResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check if all temperature fields are present and make sense
	if response.TempC == 0 && response.TempF == 0 && response.TempK == 0 {
		t.Error("All temperature values are zero, which is unexpected")
	}

	// Verify temperature conversions are correct
	expectedF := celsiusToFahrenheit(response.TempC)
	expectedK := celsiusToKelvin(response.TempC)

	if response.TempF != expectedF {
		t.Errorf("Fahrenheit conversion incorrect: got %.1f want %.1f", response.TempF, expectedF)
	}

	if response.TempK != expectedK {
		t.Errorf("Kelvin conversion incorrect: got %.1f want %.1f", response.TempK, expectedK)
	}
}

func TestGetWeatherByCEPInvalidCEP(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/weather/{cep}", getWeatherByCEP).Methods("GET")

	// Test with an invalid CEP that has correct format but doesn't exist
	req, err := http.NewRequest("GET", "/weather/99999999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	var response ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Message != "can not find zipcode" {
		t.Errorf("handler returned wrong message: got %v want %v", response.Message, "can not find zipcode")
	}
}
