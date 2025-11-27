package signalads

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		expected string
	}{
		{
			name: "with message",
			err: &APIError{
				Message:    "Test error",
				StatusCode: 400,
			},
			expected: "Test error",
		},
		{
			name: "with error message",
			err: &APIError{
				ErrorMsg:   "Error occurred",
				StatusCode: 500,
			},
			expected: "Error occurred",
		},
		{
			name: "with code",
			err: &APIError{
				Code:       "INVALID_PHONE",
				StatusCode: 400,
			},
			expected: "API error [INVALID_PHONE]: status 400",
		},
		{
			name: "with status code only",
			err: &APIError{
				StatusCode: 404,
			},
			expected: "API error: status 404",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestIsAPIError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "APIError",
			err:      &APIError{Message: "Test"},
			expected: true,
		},
		{
			name:     "regular error",
			err:      fmt.Errorf("regular error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAPIError(tt.err)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "APIError with status code",
			err:      &APIError{StatusCode: 400},
			expected: 400,
		},
		{
			name:     "regular error",
			err:      fmt.Errorf("regular error"),
			expected: 0,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStatusCode(tt.err)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGetErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "APIError with code",
			err:      &APIError{Code: "INVALID_PHONE"},
			expected: "INVALID_PHONE",
		},
		{
			name:     "APIError without code",
			err:      &APIError{Message: "Test"},
			expected: "",
		},
		{
			name:     "regular error",
			err:      fmt.Errorf("regular error"),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetErrorCode(tt.err)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestIsErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		code     string
		expected bool
	}{
		{
			name:     "matching code",
			err:      &APIError{Code: "INVALID_PHONE"},
			code:     "INVALID_PHONE",
			expected: true,
		},
		{
			name:     "non-matching code",
			err:      &APIError{Code: "INVALID_PHONE"},
			code:     "INVALID_MESSAGE",
			expected: false,
		},
		{
			name:     "regular error",
			err:      fmt.Errorf("regular error"),
			code:     "INVALID_PHONE",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsErrorCode(tt.err, tt.code)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNewAPIError(t *testing.T) {
	err := NewAPIError("TEST_CODE", "Test message", 400)

	if err.Code != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got '%s'", err.Code)
	}
	if err.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", err.Message)
	}
	if err.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", err.StatusCode)
	}
}

func TestWrapError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "regular error",
			err:      fmt.Errorf("test error"),
			expected: true,
		},
		{
			name:     "APIError",
			err:      &APIError{Message: "API error"},
			expected: true,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := WrapError(tt.err, 500)
			if tt.expected {
				if wrapped == nil {
					t.Error("Expected wrapped error, got nil")
				}
			} else {
				if wrapped != nil {
					t.Errorf("Expected nil, got %v", wrapped)
				}
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "404 status code",
			err:      &APIError{StatusCode: http.StatusNotFound},
			expected: true,
		},
		{
			name:     "NOT_FOUND code",
			err:      &APIError{Code: ErrCodeNotFound},
			expected: true,
		},
		{
			name:     "ErrNotFound variable",
			err:      ErrNotFound,
			expected: true,
		},
		{
			name:     "other error",
			err:      &APIError{StatusCode: 400},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFound(tt.err)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "401 status code",
			err:      &APIError{StatusCode: http.StatusUnauthorized},
			expected: true,
		},
		{
			name:     "INVALID_CREDENTIALS code",
			err:      &APIError{Code: ErrCodeInvalidCredentials},
			expected: true,
		},
		{
			name:     "ErrInvalidCredentials variable",
			err:      ErrInvalidCredentials,
			expected: true,
		},
		{
			name:     "other error",
			err:      &APIError{StatusCode: 400},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUnauthorized(tt.err)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsRateLimited(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "429 status code",
			err:      &APIError{StatusCode: http.StatusTooManyRequests},
			expected: true,
		},
		{
			name:     "RATE_LIMIT_EXCEEDED code",
			err:      &APIError{Code: ErrCodeRateLimitExceeded},
			expected: true,
		},
		{
			name:     "ErrRateLimited variable",
			err:      ErrRateLimited,
			expected: true,
		},
		{
			name:     "other error",
			err:      &APIError{StatusCode: 400},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRateLimited(tt.err)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsInsufficientBalance(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "402 status code",
			err:      &APIError{StatusCode: http.StatusPaymentRequired},
			expected: true,
		},
		{
			name:     "INSUFFICIENT_BALANCE code",
			err:      &APIError{Code: ErrCodeInsufficientBalance},
			expected: true,
		},
		{
			name:     "ErrInsufficientBalance variable",
			err:      ErrInsufficientBalance,
			expected: true,
		},
		{
			name:     "other error",
			err:      &APIError{StatusCode: 400},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInsufficientBalance(tt.err)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsBadRequest(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "400 status code",
			err:      &APIError{StatusCode: http.StatusBadRequest},
			expected: true,
		},
		{
			name:     "BAD_REQUEST code",
			err:      &APIError{Code: ErrCodeBadRequest},
			expected: true,
		},
		{
			name:     "other error",
			err:      &APIError{StatusCode: 500},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBadRequest(tt.err)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		code     string
		status   int
		hasError bool
	}{
		{
			name:     "ErrInvalidCredentials",
			err:      ErrInvalidCredentials,
			code:     ErrCodeInvalidCredentials,
			status:   http.StatusUnauthorized,
			hasError: true,
		},
		{
			name:     "ErrNotFound",
			err:      ErrNotFound,
			code:     ErrCodeNotFound,
			status:   http.StatusNotFound,
			hasError: true,
		},
		{
			name:     "ErrRateLimited",
			err:      ErrRateLimited,
			code:     ErrCodeRateLimitExceeded,
			status:   http.StatusTooManyRequests,
			hasError: true,
		},
		{
			name:     "ErrInsufficientBalance",
			err:      ErrInsufficientBalance,
			code:     ErrCodeInsufficientBalance,
			status:   http.StatusPaymentRequired,
			hasError: true,
		},
		{
			name:     "ErrInvalidPhoneNumber",
			err:      ErrInvalidPhoneNumber,
			code:     ErrCodeInvalidPhoneNumber,
			status:   http.StatusBadRequest,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("Expected error, got nil")
				return
			}
			if tt.err.Code != tt.code {
				t.Errorf("Expected code '%s', got '%s'", tt.code, tt.err.Code)
			}
			if tt.err.StatusCode != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, tt.err.StatusCode)
			}
			if tt.err.Error() == "" {
				t.Error("Expected error message, got empty string")
			}
		})
	}
}
