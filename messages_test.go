package signalads

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestClient(handler http.HandlerFunc) *Client {
	server := httptest.NewServer(handler)
	client := NewClient("test-api-key", "test-api-secret", WithBaseURL(server.URL))
	return client
}

func TestSendMessage(t *testing.T) {
	tests := []struct {
		name           string
		to             string
		message        string
		responseStatus int
		responseBody   interface{}
		expectError    bool
	}{
		{
			name:           "successful send",
			to:             "+989123456789",
			message:        "Test message",
			responseStatus: http.StatusOK,
			responseBody: SendMessageResponse{
				ID:     "msg-123",
				Status: "sent",
				To:     "+989123456789",
			},
			expectError: false,
		},
		{
			name:           "missing phone number",
			to:             "",
			message:        "Test message",
			responseStatus: http.StatusOK,
			responseBody:   nil,
			expectError:    true,
		},
		{
			name:           "missing message",
			to:             "+989123456789",
			message:        "",
			responseStatus: http.StatusOK,
			responseBody:   nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST, got %s", r.Method)
				}
				if r.URL.Path != "/send-message/single" {
					t.Errorf("Expected /send-message/single, got %s", r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseStatus)
				if tt.responseBody != nil {
					json.NewEncoder(w).Encode(tt.responseBody)
				}
			}

			client := setupTestClient(handler)
			ctx := context.Background()

			response, err := client.Messages.SendMessage(ctx, tt.to, tt.message)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if response != nil {
					t.Errorf("Expected nil response, got %v", response)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if response == nil {
					t.Errorf("Expected response, got nil")
				} else if response.ID != "msg-123" {
					t.Errorf("Expected ID 'msg-123', got '%s'", response.ID)
				}
			}
		})
	}
}

func TestSendMessageWithDocument(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var req SendMessageRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.DocumentLink == "" {
			t.Error("Expected document link in request")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SendMessageResponse{
			ID:     "msg-456",
			Status: "sent",
		})
	}

	client := setupTestClient(handler)
	ctx := context.Background()

	response, err := client.Messages.SendMessageWithDocument(
		ctx,
		"+989123456789",
		"Check this document",
		"https://example.com/doc.pdf",
		"Document",
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if response.ID != "msg-456" {
		t.Errorf("Expected ID 'msg-456', got '%s'", response.ID)
	}
}

func TestSendBulkMessages(t *testing.T) {
	tests := []struct {
		name        string
		messages    []BulkMessageItem
		expectError bool
	}{
		{
			name: "successful bulk send",
			messages: []BulkMessageItem{
				{To: "+989123456789", Message: "Message 1"},
				{To: "+989123456790", Message: "Message 2"},
			},
			expectError: false,
		},
		{
			name:        "empty messages",
			messages:    []BulkMessageItem{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/send-message/bulk" {
					t.Errorf("Expected /send-message/bulk, got %s", r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(SendBulkMessageResponse{
					Total:   len(tt.messages),
					Success: len(tt.messages),
					Failed:  0,
					Status:  "success",
				})
			}

			client := setupTestClient(handler)
			ctx := context.Background()

			response, err := client.Messages.SendBulkMessage(ctx, tt.messages, "")

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if response.Total != len(tt.messages) {
					t.Errorf("Expected total %d, got %d", len(tt.messages), response.Total)
				}
			}
		})
	}
}

func TestSendTemplateMessage(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var req SendTemplateMessageRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.TemplateID == "" {
			t.Error("Expected template ID in request")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SendMessageResponse{
			ID:     "msg-template-123",
			Status: "sent",
		})
	}

	client := setupTestClient(handler)
	ctx := context.Background()

	response, err := client.Messages.SendTemplate(
		ctx,
		"+989123456789",
		"template_123",
		map[string]string{"name": "علی"},
	)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if response.ID != "msg-template-123" {
		t.Errorf("Expected ID 'msg-template-123', got '%s'", response.ID)
	}
}

func TestSendVoiceMessage(t *testing.T) {
	tests := []struct {
		name        string
		req         *SendVoiceMessageRequest
		expectError bool
	}{
		{
			name: "successful voice send with text",
			req: &SendVoiceMessageRequest{
				To:        "+989123456789",
				Message:   "Test voice message",
				VoiceType: "female",
				Language:  "fa",
			},
			expectError: false,
		},
		{
			name: "successful voice send with audio URL",
			req: &SendVoiceMessageRequest{
				To:       "+989123456789",
				AudioURL: "https://example.com/audio.mp3",
			},
			expectError: false,
		},
		{
			name: "missing phone number",
			req: &SendVoiceMessageRequest{
				Message: "Test",
			},
			expectError: true,
		},
		{
			name: "missing both message and audio URL",
			req: &SendVoiceMessageRequest{
				To: "+989123456789",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(SendMessageResponse{
					ID:     "msg-voice-123",
					Status: "sent",
				})
			}

			client := setupTestClient(handler)
			ctx := context.Background()

			response, err := client.Messages.SendVoiceMessage(ctx, tt.req)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if response == nil {
					t.Errorf("Expected response, got nil")
				}
			}
		})
	}
}

func TestListMessages(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/messages" {
			t.Errorf("Expected /messages, got %s", r.URL.Path)
		}

		page := r.URL.Query().Get("page")
		if page == "" {
			t.Error("Expected page parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ListMessagesResponse{
			Messages: []Message{
				{
					ID:      "msg-1",
					To:      "+989123456789",
					Message: "Test message 1",
					Status:  "sent",
				},
				{
					ID:      "msg-2",
					To:      "+989123456790",
					Message: "Test message 2",
					Status:  "delivered",
				},
			},
			Page:    1,
			PerPage: 10,
			Total:   2,
		})
	}

	client := setupTestClient(handler)
	ctx := context.Background()

	response, err := client.Messages.ListMessages(ctx, &PaginationParams{
		Page:    1,
		PerPage: 10,
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(response.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(response.Messages))
	}
	if response.Total != 2 {
		t.Errorf("Expected total 2, got %d", response.Total)
	}
}

func TestGetMessageStatus(t *testing.T) {
	tests := []struct {
		name        string
		messageID   string
		expectError bool
	}{
		{
			name:        "successful status retrieval",
			messageID:   "msg-123",
			expectError: false,
		},
		{
			name:        "empty message ID",
			messageID:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(MessageStatus{
					ID:     tt.messageID,
					Status: "delivered",
					To:     "+989123456789",
				})
			}

			client := setupTestClient(handler)
			ctx := context.Background()

			status, err := client.Messages.GetMessageStatus(ctx, tt.messageID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if status.ID != tt.messageID {
					t.Errorf("Expected ID '%s', got '%s'", tt.messageID, status.ID)
				}
			}
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/user/info" {
			t.Errorf("Expected /user/info, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UserInfo{
			ID:       "user-123",
			Username: "testuser",
			Email:    "test@example.com",
			Balance:  1000.50,
			Credit:   500.25,
			Status:   "active",
		})
	}

	client := setupTestClient(handler)
	ctx := context.Background()

	userInfo, err := client.Messages.GetUserInfo(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if userInfo.ID != "user-123" {
		t.Errorf("Expected ID 'user-123', got '%s'", userInfo.ID)
	}
	if userInfo.Balance != 1000.50 {
		t.Errorf("Expected balance 1000.50, got %.2f", userInfo.Balance)
	}
}

func TestSendMessage_ErrorHandling(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIError{
			Message:    "Invalid phone number",
			StatusCode: http.StatusBadRequest,
		})
	}

	client := setupTestClient(handler)
	ctx := context.Background()

	response, err := client.Messages.SendMessage(ctx, "+989123456789", "Test")

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if response != nil {
		t.Errorf("Expected nil response, got %v", response)
	}

	// Check if error is wrapped APIError
	var apiErr *APIError
	if !IsAPIError(err) {
		// Try to unwrap
		if unwrapped := errors.Unwrap(err); unwrapped != nil {
			if ae, ok := unwrapped.(*APIError); ok {
				apiErr = ae
			}
		}
		if apiErr == nil {
			t.Errorf("Expected APIError, got %T: %v", err, err)
			return
		}
	} else {
		apiErr = err.(*APIError)
	}

	if apiErr.Message != "Invalid phone number" {
		t.Errorf("Expected error message 'Invalid phone number', got '%s'", apiErr.Message)
	}
}
