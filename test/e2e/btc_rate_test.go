package e2e

import (
	"encoding/json"
	"net/http"
	"testing"
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
		t.Errorf("Failed to read response body: %v", err)
	}

	_, ok := responseMap["rate"].(float64)
	if ok != true {
		t.Errorf("Failed to test rate: %v", err)
	}
}
