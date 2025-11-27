package signalads

import "net/http"

// Client represents a SignalAds API client.
// It provides methods to interact with the SignalAds API services.
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
	apiSecret  string
	Messages   *MessagesService
}

// NewClient creates a new SignalAds API client with the provided credentials.
// apiKey and apiSecret are required for authentication.
// Additional configuration can be provided using ClientOption functions.
func NewClient(apiKey, apiSecret string, opts ...ClientOption) *Client {
	client := &Client{
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}

	for _, opt := range opts {
		opt(client)
	}

	client.Messages = &MessagesService{client: client}

	return client
}
