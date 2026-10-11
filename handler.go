package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	phoneintel "github.com/eselikoghene/ng-phone-intelligence/phone"
)

type phoneSuccessResponse struct {
	PhoneNumber string `json:"phone_number"`
	Network     string `json:"network"`
	Valid       bool   `json:"valid"`
}

type apiErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiErrorResponse struct {
	Error apiErrorBody `json:"error"`
}

func phoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET is supported for this endpoint.")
		return
	}

	query := r.URL.Query()
	if !query.Has("number") {
		writeAPIError(w, http.StatusBadRequest, "MISSING_PHONE_NUMBER", "The number query parameter is required.")
		return
	}

	// Extra query parameters are ignored. Only "number" is used.
	number := query.Get("number")
	normalized, err := phoneintel.NormalizePhoneNumber(number)
	if err != nil {
		switch {
		case errors.Is(err, phoneintel.ErrEmptyInput):
			writeAPIError(w, http.StatusBadRequest, "EMPTY_PHONE_NUMBER", "The phone number must not be empty.")
		case errors.Is(err, phoneintel.ErrUnsupportedCountryCode):
			writeAPIError(w, http.StatusBadRequest, "UNSUPPORTED_COUNTRY_CODE", "Only Nigerian (+234) phone numbers are supported.")
		default:
			writeAPIError(w, http.StatusBadRequest, "INVALID_PHONE_NUMBER", "The phone number format is invalid.")
		}
		return
	}

	response := phoneSuccessResponse{
		PhoneNumber: normalized,
		Network:     phoneintel.IdentifyNetwork(normalized),
		Valid:       true,
	}

	writeJSON(w, http.StatusOK, response)
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, apiErrorResponse{
		Error: apiErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		// Do not include request input or internal details in the client response.
		log.Printf("failed to encode JSON response: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"code":"INTERNAL_ERROR","message":"Failed to encode response."}}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(payload)
	_, _ = w.Write([]byte("\n"))
}
