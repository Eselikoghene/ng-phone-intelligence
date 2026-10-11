package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPhoneHandlerSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/phone?number=08031234567", nil)
	recorder := httptest.NewRecorder()

	phoneHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var got phoneSuccessResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if got.PhoneNumber != "+2348031234567" {
		t.Errorf("phone_number: expected %q, got %q", "+2348031234567", got.PhoneNumber)
	}
	if got.Network != "MTN" {
		t.Errorf("network: expected %q, got %q", "MTN", got.Network)
	}
	if !got.Valid {
		t.Errorf("valid: expected true, got false")
	}
}

func TestPhoneHandlerSuccessInternational(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/phone?number=%2B2348031234567", nil)
	recorder := httptest.NewRecorder()

	phoneHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var got phoneSuccessResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if got.PhoneNumber != "+2348031234567" || got.Network != "MTN" || !got.Valid {
		t.Fatalf("unexpected success body: %+v", got)
	}
}

func TestPhoneHandlerUnknownNetworkStillSucceeds(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/phone?number=08001234567", nil)
	recorder := httptest.NewRecorder()

	phoneHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var got phoneSuccessResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if got.PhoneNumber != "+2348001234567" {
		t.Errorf("phone_number: expected %q, got %q", "+2348001234567", got.PhoneNumber)
	}
	if got.Network != "Unknown" {
		t.Errorf("network: expected %q, got %q", "Unknown", got.Network)
	}
	if !got.Valid {
		t.Errorf("valid: expected true, got false")
	}
}

func TestPhoneHandlerIgnoresExtraQueryParameters(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/phone?number=08031234567&debug=1", nil)
	recorder := httptest.NewRecorder()

	phoneHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var got phoneSuccessResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if got.Network != "MTN" {
		t.Errorf("network: expected %q, got %q", "MTN", got.Network)
	}
}

func TestPhoneHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/phone?number=08031234567", nil)
	recorder := httptest.NewRecorder()

	phoneHandler(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, recorder.Code)
	}
	if got := recorder.Header().Get("Allow"); got != http.MethodGet {
		t.Fatalf("expected Allow %q, got %q", http.MethodGet, got)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", got)
	}

	var body apiErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if body.Error.Code != "METHOD_NOT_ALLOWED" {
		t.Errorf("error.code: expected %q, got %q", "METHOD_NOT_ALLOWED", body.Error.Code)
	}
}

func TestPhoneHandlerErrors(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		expectedCode string
		expectedMsg  string
	}{
		{
			name:         "missing number query parameter",
			url:          "/api/v1/phone",
			expectedCode: "MISSING_PHONE_NUMBER",
			expectedMsg:  "The number query parameter is required.",
		},
		{
			name:         "empty number",
			url:          "/api/v1/phone?number=",
			expectedCode: "EMPTY_PHONE_NUMBER",
			expectedMsg:  "The phone number must not be empty.",
		},
		{
			name:         "whitespace only number",
			url:          "/api/v1/phone?number=%20%20",
			expectedCode: "EMPTY_PHONE_NUMBER",
			expectedMsg:  "The phone number must not be empty.",
		},
		{
			name:         "invalid phone number format",
			url:          "/api/v1/phone?number=12345",
			expectedCode: "INVALID_PHONE_NUMBER",
			expectedMsg:  "The phone number format is invalid.",
		},
		{
			name:         "unsupported country code",
			url:          "/api/v1/phone?number=%2B18031234567",
			expectedCode: "UNSUPPORTED_COUNTRY_CODE",
			expectedMsg:  "Only Nigerian (+234) phone numbers are supported.",
		},
		{
			name:         "bare international without plus",
			url:          "/api/v1/phone?number=2348031234567",
			expectedCode: "INVALID_PHONE_NUMBER",
			expectedMsg:  "The phone number format is invalid.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.url, nil)
			recorder := httptest.NewRecorder()

			phoneHandler(recorder, req)

			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
			}

			contentType := recorder.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Fatalf("expected Content-Type application/json, got %q", contentType)
			}

			var got apiErrorResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
				t.Fatalf("response is not valid JSON: %v", err)
			}

			if got.Error.Code != test.expectedCode {
				t.Errorf("error.code: expected %q, got %q", test.expectedCode, got.Error.Code)
			}
			if got.Error.Message != test.expectedMsg {
				t.Errorf("error.message: expected %q, got %q", test.expectedMsg, got.Error.Message)
			}
		})
	}
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	healthHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", got)
	}

	var body map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status: expected %q, got %q", "ok", body["status"])
	}
}

func BenchmarkPhoneHandler(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/phone?number=08031234567", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		recorder := httptest.NewRecorder()
		phoneHandler(recorder, req)
		if recorder.Code != http.StatusOK {
			b.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
		}
	}
}
