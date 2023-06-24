package e2e

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func TestBTCRateE2E(t *testing.T) {
	var responseMap map[string]interface{}
	resp, err := http.Get("http://localhost:8080/api/rate")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	rate, ok := responseMap["rate"].(float64)
	if ok != true {
		t.Fatalf("Failed to test rate: %v", err)
	}

	expectedRate, err := getTestRate()
	if err != nil {
		t.Fatalf("Failed to load expected test rate: %v", err)
	}

	if rate != expectedRate {
		t.Errorf("Expected rate %f, but got %f", expectedRate, rate)
	}
}

func getTestRate() (float64, error) {
	err := godotenv.Load("../../services/currency/.env.test")
	if err != nil {
		panic(errors.Wrap(err, "Can not load config"))
	}

	rate, err := strconv.ParseFloat(os.Getenv("TEST_RATE"), 64)
	if err != nil {
		return 0, errors.Wrap(err, "Can not load TEST_RATE")
	}

	return rate, nil
}
