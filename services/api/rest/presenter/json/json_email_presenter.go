package json

import (
	"api/logger"
	"net/http"
)

type JSONEmailPresenter struct{}

func (p *JSONEmailPresenter) SuccessfulEmailsSending(w http.ResponseWriter) {
	logger.DefaultLog(logger.DEBUG, "decoding emails send successfully response")
	EncodeJSONResponse(w, HTTPResponse{Description: "all emails were sent successfully"})
}

func (p *JSONEmailPresenter) SuccessfullyAddEmail(w http.ResponseWriter) {
	logger.DefaultLog(logger.DEBUG, "decoding email was added successfully response")
	EncodeJSONResponse(w, HTTPResponse{Description: "email was successfully added"})
}
