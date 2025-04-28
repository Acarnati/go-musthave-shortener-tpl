package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid POST request",
			method:         http.MethodPost,
			body:           "https://example.com",
			expectedStatus: http.StatusCreated,
			expectedBody:   "http://localhost:8080/",
		},
		{
			name:           "Wrong method",
			method:         http.MethodGet,
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "",
		},
		{
			name:           "Empty body",
			method:         http.MethodPost,
			body:           "",
			expectedStatus: http.StatusCreated,
			expectedBody:   "http://localhost:8080/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/shorten", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(shortHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedBody != "" && !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want it to contain %v",
					rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	// Pre-populate urlStore with test data
	testKey := "abc123"
	testURL := "https://example.com"
	urlStore[testKey] = testURL

	tests := []struct {
		name           string
		method         string
		id             string
		expectedStatus int
		expectedURL    string
	}{
		{
			name:           "Valid redirect",
			method:         http.MethodGet,
			id:             testKey,
			expectedStatus: http.StatusTemporaryRedirect,
			expectedURL:    testURL,
		},
		{
			name:           "Wrong method",
			method:         http.MethodPost,
			id:             testKey,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedURL:    "",
		},
		{
			name:           "Non-existent ID",
			method:         http.MethodGet,
			id:             "nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedURL:    "",
		},
		{
			name:           "Empty ID",
			method:         http.MethodGet,
			id:             "",
			expectedStatus: http.StatusBadRequest,
			expectedURL:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Set path parameter
			if tt.id != "" {
				req.SetPathValue("id", tt.id)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(redirectHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedURL != "" && rr.Header().Get("Location") != tt.expectedURL {
				t.Errorf("handler returned wrong Location header: got %v want %v",
					rr.Header().Get("Location"), tt.expectedURL)
			}
		})
	}
}

func TestGenerateShortKey(t *testing.T) {
	key1 := generateShortKey()
	key2 := generateShortKey()

	if len(key1) != 8 {
		t.Errorf("generateShortKey() returned key with wrong length: got %v want 8", len(key1))
	}

	if key1 == key2 {
		t.Errorf("generateShortKey() returned duplicate keys: %v", key1)
	}
}
