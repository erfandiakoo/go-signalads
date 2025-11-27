package signalads

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Code       string                 `json:"code,omitempty"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"status_code,omitempty"`
	ErrorMsg   string                 `json:"error,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.ErrorMsg != "" {
		return e.ErrorMsg
	}
	if e.Code != "" {
		return fmt.Sprintf("API error [%s]: status %d", e.Code, e.StatusCode)
	}
	return fmt.Sprintf("API error: status %d", e.StatusCode)
}

const (
	//nolint:gosec // This is an error code constant, not a credential
	ErrCodeInvalidCredentials  = "INVALID_CREDENTIALS"
	ErrCodeInvalidPhoneNumber  = "INVALID_PHONE_NUMBER"
	ErrCodeInvalidMessage      = "INVALID_MESSAGE"
	ErrCodeInsufficientBalance = "INSUFFICIENT_BALANCE"
	ErrCodeRateLimitExceeded   = "RATE_LIMIT_EXCEEDED"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeForbidden           = "FORBIDDEN"
	ErrCodeBadRequest          = "BAD_REQUEST"
	ErrCodeInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrCodeServiceUnavailable  = "SERVICE_UNAVAILABLE"
	ErrCodeInvalidTemplate     = "INVALID_TEMPLATE"
	ErrCodeTemplateNotApproved = "TEMPLATE_NOT_APPROVED"
	ErrCodeInvalidDocument     = "INVALID_DOCUMENT"
	ErrCodeInvalidVoiceFormat  = "INVALID_VOICE_FORMAT"
)

var (
	ErrInvalidCredentials = &APIError{
		Code:       ErrCodeInvalidCredentials,
		Message:    "Invalid API credentials",
		StatusCode: http.StatusUnauthorized,
	}

	ErrNotFound = &APIError{
		Code:       ErrCodeNotFound,
		Message:    "Resource not found",
		StatusCode: http.StatusNotFound,
	}

	ErrRateLimited = &APIError{
		Code:       ErrCodeRateLimitExceeded,
		Message:    "Rate limit exceeded",
		StatusCode: http.StatusTooManyRequests,
	}

	ErrInsufficientBalance = &APIError{
		Code:       ErrCodeInsufficientBalance,
		Message:    "Insufficient balance",
		StatusCode: http.StatusPaymentRequired,
	}

	ErrInvalidPhoneNumber = &APIError{
		Code:       ErrCodeInvalidPhoneNumber,
		Message:    "Invalid phone number",
		StatusCode: http.StatusBadRequest,
	}
)

func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}

func GetStatusCode(err error) int {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode
	}
	return 0
}

func GetErrorCode(err error) string {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code
	}
	return ""
}

func IsErrorCode(err error, code string) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code == code
	}
	return false
}

func NewAPIError(code, message string, statusCode int) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

func WrapError(err error, statusCode int) *APIError {
	if err == nil {
		return nil
	}

	if apiErr, ok := err.(*APIError); ok {
		return apiErr
	}

	return &APIError{
		Message:    err.Error(),
		StatusCode: statusCode,
	}
}

func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound || apiErr.Code == ErrCodeNotFound
	}
	return false
}

func IsUnauthorized(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusUnauthorized || apiErr.Code == ErrCodeInvalidCredentials
	}
	return false
}

func IsRateLimited(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusTooManyRequests || apiErr.Code == ErrCodeRateLimitExceeded
	}
	return false
}

func IsInsufficientBalance(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code == ErrCodeInsufficientBalance || apiErr.StatusCode == http.StatusPaymentRequired
	}
	return false
}

func IsBadRequest(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusBadRequest || apiErr.Code == ErrCodeBadRequest
	}
	return false
}
