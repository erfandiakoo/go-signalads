package signalads

import (
	"time"
)

// APIError is defined in errors.go

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// Message types for SMS functionality

// SendMessageRequest represents a request to send a single SMS message
type SendMessageRequest struct {
	// Recipient phone number (required)
	To string `json:"to"`

	// Message text (required)
	Message string `json:"message"`

	// Document link/URL to attach (optional)
	DocumentLink string `json:"document_link,omitempty"`

	// Document caption (optional, used with document_link)
	DocumentCaption string `json:"document_caption,omitempty"`

	// Sender ID or phone number (optional)
	From string `json:"from,omitempty"`

	// Additional parameters that may be supported by the API
	Params map[string]interface{} `json:"params,omitempty"`
}

// SendMessageResponse represents the response from sending a message
type SendMessageResponse struct {
	// Message ID if successful
	ID string `json:"id,omitempty"`

	// Status of the message
	Status string `json:"status"`

	// Status message
	Message string `json:"message,omitempty"`

	// Recipient phone number
	To string `json:"to,omitempty"`

	// Cost of the message (if available)
	Cost float64 `json:"cost,omitempty"`

	// Timestamp when message was sent
	SentAt string `json:"sent_at,omitempty"`

	// Additional response data
	Data map[string]interface{} `json:"data,omitempty"`
}

// BulkMessageItem represents a single message in a bulk send request
type BulkMessageItem struct {
	To      string            `json:"to"`
	Message string            `json:"message"`
	Params  map[string]string `json:"params,omitempty"`
}

// SendBulkMessageRequest represents a request to send bulk/group messages
type SendBulkMessageRequest struct {
	// List of messages to send
	Messages []BulkMessageItem `json:"messages"`

	// Sender ID or phone number (optional)
	From string `json:"from,omitempty"`

	// Additional parameters
	Params map[string]interface{} `json:"params,omitempty"`
}

// SendBulkMessageResponse represents the response from bulk sending
type SendBulkMessageResponse struct {
	// Total number of messages sent
	Total int `json:"total"`

	// Number of successful messages
	Success int `json:"success"`

	// Number of failed messages
	Failed int `json:"failed"`

	// List of message IDs
	MessageIDs []string `json:"message_ids,omitempty"`

	// Detailed results for each message
	Results []SendMessageResponse `json:"results,omitempty"`

	// Status
	Status string `json:"status"`

	// Additional data
	Data map[string]interface{} `json:"data,omitempty"`
}

// SendTemplateMessageRequest represents a request to send a template message
type SendTemplateMessageRequest struct {
	// Recipient phone number (required)
	To string `json:"to"`

	// Template ID (required)
	TemplateID string `json:"template_id"`

	// Template parameters (optional, for dynamic templates)
	TemplateParams map[string]string `json:"template_params,omitempty"`

	// Sender ID or phone number (optional)
	From string `json:"from,omitempty"`

	// Additional parameters
	Params map[string]interface{} `json:"params,omitempty"`
}

// SendVoiceMessageRequest represents a request to send a voice/audio message
type SendVoiceMessageRequest struct {
	// Recipient phone number (required)
	To string `json:"to"`

	// Voice message text (required) - will be converted to speech
	Message string `json:"message"`

	// Audio file URL (optional, if provided, this will be used instead of text-to-speech)
	AudioURL string `json:"audio_url,omitempty"`

	// Voice type (optional, e.g., "male", "female")
	VoiceType string `json:"voice_type,omitempty"`

	// Language (optional, e.g., "fa", "en")
	Language string `json:"language,omitempty"`

	// Sender ID or phone number (optional)
	From string `json:"from,omitempty"`

	// Additional parameters
	Params map[string]interface{} `json:"params,omitempty"`
}

// Message represents a message in the list
type Message struct {
	ID          string    `json:"id"`
	To          string    `json:"to"`
	From        string    `json:"from,omitempty"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	Cost        float64   `json:"cost,omitempty"`
	SentAt      time.Time `json:"sent_at,omitempty"`
	DeliveredAt time.Time `json:"delivered_at,omitempty"`
	ReadAt      time.Time `json:"read_at,omitempty"`
	Error       string    `json:"error,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

// ListMessagesResponse represents the response from listing messages
type ListMessagesResponse struct {
	Messages []Message `json:"messages"`
	Page     int       `json:"page"`
	PerPage  int       `json:"per_page"`
	Total    int       `json:"total"`
}

// MessageStatus represents the status of a message
type MessageStatus struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"` // e.g., "sent", "delivered", "failed", "pending"
	To          string    `json:"to"`
	SentAt      time.Time `json:"sent_at,omitempty"`
	DeliveredAt time.Time `json:"delivered_at,omitempty"`
	ReadAt      time.Time `json:"read_at,omitempty"`
	Error       string    `json:"error,omitempty"`
	Cost        float64   `json:"cost,omitempty"`
}

// UserInfo represents user account information
type UserInfo struct {
	ID          string    `json:"id"`
	Username    string    `json:"username,omitempty"`
	Email       string    `json:"email,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Balance     float64   `json:"balance,omitempty"`
	Credit      float64   `json:"credit,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
}
