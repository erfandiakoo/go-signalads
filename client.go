package signalads

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultBaseURL = "https://panel.signalads.com/api/v1"
	DefaultTimeout = 30 * time.Second
)

// ClientOption is a function type for configuring a Client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API client.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client for the API client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets a custom timeout for API requests.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if c.httpClient != nil {
			c.httpClient.Timeout = timeout
		}
	}
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}, queryParams map[string]string) (*http.Response, error) {
	reqURL := c.baseURL + endpoint
	if len(queryParams) > 0 {
		u, err := url.Parse(reqURL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}
		q := u.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		reqURL = u.String()
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("X-API-Secret", c.apiSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func (c *Client) parseResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err == nil {
			if apiErr.StatusCode == 0 {
				apiErr.StatusCode = resp.StatusCode
			}
			if apiErr.Message != "" || apiErr.Code != "" || apiErr.ErrorMsg != "" {
				return &apiErr
			}
		}
		return NewAPIError(
			getErrorCodeFromStatusCode(resp.StatusCode),
			fmt.Sprintf("API error: status %d, body: %s", resp.StatusCode, string(body)),
			resp.StatusCode,
		)
	}

	if v != nil {
		if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func getErrorCodeFromStatusCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return ErrCodeBadRequest
	case http.StatusUnauthorized:
		return ErrCodeInvalidCredentials
	case http.StatusForbidden:
		return ErrCodeForbidden
	case http.StatusNotFound:
		return ErrCodeNotFound
	case http.StatusPaymentRequired:
		return ErrCodeInsufficientBalance
	case http.StatusTooManyRequests:
		return ErrCodeRateLimitExceeded
	case http.StatusInternalServerError:
		return ErrCodeInternalServerError
	case http.StatusServiceUnavailable:
		return ErrCodeServiceUnavailable
	default:
		return ""
	}
}

// Get performs a GET request to the specified endpoint.
func (c *Client) Get(ctx context.Context, endpoint string, result interface{}, queryParams map[string]string) error {
	resp, err := c.doRequest(ctx, http.MethodGet, endpoint, nil, queryParams)
	if err != nil {
		return err
	}
	return c.parseResponse(resp, result)
}

// Post performs a POST request to the specified endpoint.
func (c *Client) Post(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodPost, endpoint, body, nil)
	if err != nil {
		return err
	}
	return c.parseResponse(resp, result)
}

// Put performs a PUT request to the specified endpoint.
func (c *Client) Put(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodPut, endpoint, body, nil)
	if err != nil {
		return err
	}
	return c.parseResponse(resp, result)
}

// Delete performs a DELETE request to the specified endpoint.
func (c *Client) Delete(ctx context.Context, endpoint string, result interface{}) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}
	return c.parseResponse(resp, result)
}
