package signalads

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-key", "test-secret")

	if client == nil {
		t.Fatal("Expected client, got nil")
	}
	if client.apiKey != "test-key" {
		t.Errorf("Expected apiKey 'test-key', got '%s'", client.apiKey)
	}
	if client.apiSecret != "test-secret" {
		t.Errorf("Expected apiSecret 'test-secret', got '%s'", client.apiSecret)
	}
	if client.baseURL != DefaultBaseURL {
		t.Errorf("Expected baseURL '%s', got '%s'", DefaultBaseURL, client.baseURL)
	}
	if client.Messages == nil {
		t.Error("Expected Messages service, got nil")
	}
}

func TestClientOptions(t *testing.T) {
	customURL := "https://custom.example.com/api"
	customTimeout := 60 * time.Second

	client := NewClient(
		"test-key",
		"test-secret",
		WithBaseURL(customURL),
		WithTimeout(customTimeout),
	)

	if client.baseURL != customURL {
		t.Errorf("Expected baseURL '%s', got '%s'", customURL, client.baseURL)
	}
	if client.httpClient.Timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, client.httpClient.Timeout)
	}
}

func TestClient_Get(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		// Check authentication headers
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Error("Missing or incorrect X-API-Key header")
		}
		if r.Header.Get("X-API-Secret") != "test-secret" {
			t.Error("Missing or incorrect X-API-Secret header")
		}

		// Check query parameters
		param := r.URL.Query().Get("test")
		if param != "value" {
			t.Errorf("Expected query param 'value', got '%s'", param)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient("test-key", "test-secret", WithBaseURL(server.URL))

	var result map[string]string
	err := client.Get(
		context.Background(),
		"/test",
		&result,
		map[string]string{"test": "value"},
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result["status"])
	}
}

func TestClient_Post(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Missing Content-Type header")
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)

		if body["test"] != "value" {
			t.Errorf("Expected body test='value', got '%s'", body["test"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "123"})
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient("test-key", "test-secret", WithBaseURL(server.URL))

	var result map[string]string
	err := client.Post(
		context.Background(),
		"/test",
		map[string]string{"test": "value"},
		&result,
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result["id"] != "123" {
		t.Errorf("Expected id '123', got '%s'", result["id"])
	}
}

func TestClient_ParseResponse_Error(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIError{
			Message:    "Test error",
			StatusCode: http.StatusBadRequest,
		})
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient("test-key", "test-secret", WithBaseURL(server.URL))

	var result map[string]string
	err := client.Get(context.Background(), "/test", &result, nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Errorf("Expected APIError, got %T", err)
	} else {
		if apiErr.Message != "Test error" {
			t.Errorf("Expected error message 'Test error', got '%s'", apiErr.Message)
		}
		if apiErr.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
		}
	}
}

func TestClient_ParseResponse_NonJSONError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient("test-key", "test-secret", WithBaseURL(server.URL))

	var result map[string]string
	err := client.Get(context.Background(), "/test", &result, nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestClient_ContextCancellation(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient("test-key", "test-secret", WithBaseURL(server.URL))

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	var result map[string]string
	err := client.Get(ctx, "/test", &result, nil)

	if err == nil {
		t.Error("Expected error due to context cancellation, got nil")
	}
}
