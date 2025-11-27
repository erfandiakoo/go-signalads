package signalads

import (
	"context"
	"fmt"
)

// MessagesService provides methods for sending and managing SMS messages.
type MessagesService struct {
	client *Client
}

// SendSingleMessage sends a single SMS message with optional document link.
func (s *MessagesService) SendSingleMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.To == "" {
		return nil, fmt.Errorf("recipient phone number is required")
	}
	if req.Message == "" {
		return nil, fmt.Errorf("message text is required")
	}

	var response SendMessageResponse
	if err := s.client.Post(ctx, "/send-message/single", req, &response); err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return &response, nil
}

// SendMessage sends a simple text message to the specified phone number.
func (s *MessagesService) SendMessage(ctx context.Context, to, message string) (*SendMessageResponse, error) {
	return s.SendSingleMessage(ctx, &SendMessageRequest{
		To:      to,
		Message: message,
	})
}

// SendMessageWithDocument sends a message with a document link.
func (s *MessagesService) SendMessageWithDocument(ctx context.Context, to, message, documentLink, caption string) (*SendMessageResponse, error) {
	return s.SendSingleMessage(ctx, &SendMessageRequest{
		To:              to,
		Message:         message,
		DocumentLink:    documentLink,
		DocumentCaption: caption,
	})
}

// SendBulkMessages sends multiple messages in a single request.
func (s *MessagesService) SendBulkMessages(ctx context.Context, req *SendBulkMessageRequest) (*SendBulkMessageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("at least one message is required")
	}

	var response SendBulkMessageResponse
	if err := s.client.Post(ctx, "/send-message/bulk", req, &response); err != nil {
		return nil, fmt.Errorf("failed to send bulk messages: %w", err)
	}

	return &response, nil
}

// SendBulkMessage is a convenience method for sending bulk messages.
func (s *MessagesService) SendBulkMessage(ctx context.Context, messages []BulkMessageItem, from string) (*SendBulkMessageResponse, error) {
	return s.SendBulkMessages(ctx, &SendBulkMessageRequest{
		Messages: messages,
		From:     from,
	})
}

// SendTemplateMessage sends a message using a predefined template.
func (s *MessagesService) SendTemplateMessage(ctx context.Context, req *SendTemplateMessageRequest) (*SendMessageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.To == "" {
		return nil, fmt.Errorf("recipient phone number is required")
	}
	if req.TemplateID == "" {
		return nil, fmt.Errorf("template ID is required")
	}

	var response SendMessageResponse
	if err := s.client.Post(ctx, "/send-message/template", req, &response); err != nil {
		return nil, fmt.Errorf("failed to send template message: %w", err)
	}

	return &response, nil
}

// SendTemplate is a convenience method for sending template messages.
func (s *MessagesService) SendTemplate(ctx context.Context, to, templateID string, params map[string]string) (*SendMessageResponse, error) {
	return s.SendTemplateMessage(ctx, &SendTemplateMessageRequest{
		To:             to,
		TemplateID:     templateID,
		TemplateParams: params,
	})
}

// SendVoiceMessage sends a voice or audio message.
func (s *MessagesService) SendVoiceMessage(ctx context.Context, req *SendVoiceMessageRequest) (*SendMessageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.To == "" {
		return nil, fmt.Errorf("recipient phone number is required")
	}
	if req.Message == "" && req.AudioURL == "" {
		return nil, fmt.Errorf("either message text or audio URL is required")
	}

	var response SendMessageResponse
	if err := s.client.Post(ctx, "/send-message/voice", req, &response); err != nil {
		return nil, fmt.Errorf("failed to send voice message: %w", err)
	}

	return &response, nil
}

// SendVoice is a convenience method for sending voice messages.
func (s *MessagesService) SendVoice(ctx context.Context, to, message, voiceType, language string) (*SendMessageResponse, error) {
	return s.SendVoiceMessage(ctx, &SendVoiceMessageRequest{
		To:        to,
		Message:   message,
		VoiceType: voiceType,
		Language:  language,
	})
}

// ListMessages retrieves a list of messages with optional pagination.
func (s *MessagesService) ListMessages(ctx context.Context, params *PaginationParams) (*ListMessagesResponse, error) {
	queryParams := make(map[string]string, 2)
	if params != nil {
		if params.Page > 0 {
			queryParams["page"] = fmt.Sprintf("%d", params.Page)
		}
		if params.PerPage > 0 {
			queryParams["per_page"] = fmt.Sprintf("%d", params.PerPage)
		}
	}

	var response ListMessagesResponse
	if err := s.client.Get(ctx, "/messages", &response, queryParams); err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	return &response, nil
}

// GetMessageStatus retrieves the status of a specific message by its ID.
func (s *MessagesService) GetMessageStatus(ctx context.Context, messageID string) (*MessageStatus, error) {
	if messageID == "" {
		return nil, fmt.Errorf("message ID is required")
	}

	var status MessageStatus
	if err := s.client.Get(ctx, "/messages/"+messageID+"/status", &status, nil); err != nil {
		return nil, fmt.Errorf("failed to get message status: %w", err)
	}

	return &status, nil
}

// GetUserInfo retrieves the current user's account information.
func (s *MessagesService) GetUserInfo(ctx context.Context) (*UserInfo, error) {
	var userInfo UserInfo
	if err := s.client.Get(ctx, "/user/info", &userInfo, nil); err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return &userInfo, nil
}
