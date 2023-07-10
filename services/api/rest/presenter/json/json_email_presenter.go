package json

import (
	"net/http"
)

type JSONEmailPresenter struct{}

func (p *JSONEmailPresenter) SuccessfulEmailsSending(w http.ResponseWriter) {
	DecodeJSONResponse(w, HTTPResponse{Description: "all emails were sent successfully"})
}

func (p *JSONEmailPresenter) SuccessfullyAddEmail(w http.ResponseWriter) {
	DecodeJSONResponse(w, HTTPResponse{Description: "email was successfully added"})
}
