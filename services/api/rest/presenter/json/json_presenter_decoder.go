package json

import (
	"api/logger"
	"encoding/json"
	"net/http"
)

type HTTPResponse struct {
	Description string `json:"description"`
}

func EncodeJSONResponse(w http.ResponseWriter, toDecode any) {
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(toDecode)
	if err != nil {
		logger.DefaultLog(logger.ERROR, "can not decode response")
		http.Error(w, "can not decode response", http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(jsonData); err != nil {
		logger.DefaultLog(logger.ERROR, "can not decode response")
		http.Error(w, "can not decode response", http.StatusInternalServerError)
		return
	}
}
