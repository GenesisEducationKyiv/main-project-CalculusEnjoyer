package json

import (
	"encoding/json"
	"net/http"
)

type HTTPResponse struct {
	Description string `json:"description"`
}

func DecodeJSONResponse(w http.ResponseWriter, toDecode any) {
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(toDecode)
	if err != nil {
		http.Error(w, "can not decode response", http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(jsonData); err != nil {
		http.Error(w, "can not decode response", http.StatusInternalServerError)
		return
	}
}
