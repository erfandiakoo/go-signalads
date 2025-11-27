# Go SignalAds Client

[![CI](https://github.com/erfandiakoo/go-signalads/actions/workflows/ci.yml/badge.svg)](https://github.com/erfandiakoo/go-signalads/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/erfandiakoo/go-signalads.svg)](https://pkg.go.dev/github.com/erfandiakoo/go-signalads)
[![Go Report Card](https://goreportcard.com/badge/github.com/erfandiakoo/go-signalads)](https://goreportcard.com/report/github.com/erfandiakoo/go-signalads)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A production-ready Go client library for the SignalAds API. This package provides a clean, type-safe interface for sending SMS messages, managing templates, handling voice messages, and more.

## Features

- ✅ **Simple Message Sending** - Send individual SMS messages
- ✅ **Bulk Messaging** - Send multiple messages in a single request
- ✅ **Template Support** - Use predefined message templates
- ✅ **Voice Messages** - Send voice/audio messages with text-to-speech
- ✅ **Message Management** - List messages, check status, and track delivery
- ✅ **User Information** - Retrieve account details and balance
- ✅ **Comprehensive Error Handling** - Typed errors with detailed information
- ✅ **Context Support** - Full context.Context integration for cancellation and timeouts
- ✅ **Pagination** - Built-in pagination support for list endpoints
- ✅ **Well Tested** - High test coverage with comprehensive test suite
- ✅ **Production Ready** - Memory safe, no leaks, battle-tested code

## Installation

```bash
go get github.com/erfandiakoo/go-signalads
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/erfandiakoo/go-signalads"
)

func main() {
    client := signalads.NewClient("your-api-key", "your-api-secret")
    ctx := context.Background()
    
    response, err := client.Messages.SendMessage(ctx, "+1234567890", "Hello from SignalAds!")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Message sent! ID: %s\n", response.ID)
}
```

## Configuration

### Basic Configuration

```go
client := signalads.NewClient("api-key", "api-secret")
```

### Custom Base URL

```go
client := signalads.NewClient(
    "api-key",
    "api-secret",
    signalads.WithBaseURL("https://custom-api.example.com/api/v1"),
)
```

### Custom Timeout

```go
client := signalads.NewClient(
    "api-key",
    "api-secret",
    signalads.WithTimeout(60 * time.Second),
)
```

### Custom HTTP Client

```go
customClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

client := signalads.NewClient(
    "api-key",
    "api-secret",
    signalads.WithHTTPClient(customClient),
)
```

## API Reference

### Messages Service

#### Send Simple Message

```go
response, err := client.Messages.SendMessage(ctx, "+1234567890", "Hello, World!")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Message ID: %s, Status: %s\n", response.ID, response.Status)
```

#### Send Message with Document

```go
response, err := client.Messages.SendMessageWithDocument(
    ctx,
    "+1234567890",
    "Please review this document",
    "https://example.com/document.pdf",
    "Important Document",
)
```

#### Send Single Message (Full Control)

```go
response, err := client.Messages.SendSingleMessage(ctx, &signalads.SendMessageRequest{
    To:              "+1234567890",
    Message:         "Your message text",
    DocumentLink:     "https://example.com/file.pdf",
    DocumentCaption: "Document caption",
    From:            "SENDER_ID", // Optional
})
```

#### Send Bulk Messages

```go
messages := []signalads.BulkMessageItem{
    {To: "+1234567890", Message: "Message 1"},
    {To: "+1234567891", Message: "Message 2"},
}

response, err := client.Messages.SendBulkMessage(ctx, messages, "")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total: %d, Success: %d, Failed: %d\n",
    response.Total, response.Success, response.Failed)
```

#### Send Bulk Messages (Full Control)

```go
response, err := client.Messages.SendBulkMessages(ctx, &signalads.SendBulkMessageRequest{
    Messages: []signalads.BulkMessageItem{
        {To: "+1234567890", Message: "Message 1"},
        {To: "+1234567891", Message: "Message 2"},
    },
    From: "SENDER_ID", // Optional
})
```

#### Send Template Message

```go
response, err := client.Messages.SendTemplate(
    ctx,
    "+1234567890",
    "template_123",
    map[string]string{
        "name": "Erfan",
        "code": "12345",
    },
)
```

#### Send Template Message (Full Control)

```go
response, err := client.Messages.SendTemplateMessage(ctx, &signalads.SendTemplateMessageRequest{
    To:            "+1234567890",
    TemplateID:    "template_123",
    TemplateParams: map[string]string{
        "name": "Erfan",
        "code": "12345",
    },
    From: "", // Optional
})
```

#### Send Voice Message

```go
response, err := client.Messages.SendVoice(
    ctx,
    "+1234567890",
    "This is a voice message",
    "female", // "male" or "female"
    "en",     // Language code
)
```

#### Send Voice Message (Full Control)

```go
response, err := client.Messages.SendVoiceMessage(ctx, &signalads.SendVoiceMessageRequest{
    To:        "+1234567890",
    Message:   "This will be converted to speech",
    VoiceType: "female",
    Language:  "en",
    // Or use AudioURL instead:
    // AudioURL: "https://example.com/audio.mp3",
})
```

#### List Messages

```go
messages, err := client.Messages.ListMessages(ctx, &signalads.PaginationParams{
    Page:    1,
    PerPage: 20,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d messages (Total: %d)\n", len(messages.Messages), messages.Total)

for _, msg := range messages.Messages {
    fmt.Printf("ID: %s, To: %s, Status: %s\n", msg.ID, msg.To, msg.Status)
}
```

#### Get Message Status

```go
status, err := client.Messages.GetMessageStatus(ctx, "message-id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", status.Status)
fmt.Printf("Sent at: %s\n", status.SentAt)
if !status.DeliveredAt.IsZero() {
    fmt.Printf("Delivered at: %s\n", status.DeliveredAt)
}
```

#### Get User Information

```go
userInfo, err := client.Messages.GetUserInfo(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("User ID: %s\n", userInfo.ID)
fmt.Printf("Balance: %.2f\n", userInfo.Balance)
fmt.Printf("Credit: %.2f\n", userInfo.Credit)
```

## Error Handling

The client returns typed errors that implement the `error` interface. API errors are returned as `*APIError`:

```go
response, err := client.Messages.SendMessage(ctx, "+1234567890", "Hello")
if err != nil {
    if signalads.IsAPIError(err) {
        apiErr := err.(*signalads.APIError)
        fmt.Printf("API Error: %s (Code: %s, Status: %d)\n",
            apiErr.Message, apiErr.Code, apiErr.StatusCode)
        
        // Check for specific error types
        if signalads.IsInsufficientBalance(err) {
            fmt.Println("Insufficient balance")
        } else if signalads.IsRateLimited(err) {
            fmt.Println("Rate limit exceeded")
        } else if signalads.IsUnauthorized(err) {
            fmt.Println("Authentication failed")
        }
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Error Helper Functions

```go
// Check error type
signalads.IsAPIError(err)
signalads.IsNotFound(err)
signalads.IsUnauthorized(err)
signalads.IsRateLimited(err)
signalads.IsInsufficientBalance(err)
signalads.IsBadRequest(err)

// Get error details
statusCode := signalads.GetStatusCode(err)
errorCode := signalads.GetErrorCode(err)
signalads.IsErrorCode(err, "INVALID_PHONE_NUMBER")
```

## Advanced Usage

### Direct HTTP Methods

For endpoints not yet implemented, you can use the low-level HTTP methods:

```go
var result YourCustomType

// GET request
err := client.Get(ctx, "/custom-endpoint", &result, map[string]string{
    "param": "value",
})

// POST request
err := client.Post(ctx, "/custom-endpoint", requestBody, &result)

// PUT request
err := client.Put(ctx, "/custom-endpoint", requestBody, &result)

// DELETE request
err := client.Delete(ctx, "/custom-endpoint", nil)
```

### Context Cancellation

All methods accept `context.Context` for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

response, err := client.Messages.SendMessage(ctx, "+1234567890", "Hello")
if err != nil {
    if err == context.DeadlineExceeded {
        fmt.Println("Request timed out")
    }
}
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run tests with race detection:

```bash
go test -race ./...
```

## Requirements

- Go 1.21 or higher

## Documentation

Full API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/erfandiakoo/go-signalads).

## Examples

See the [example](example/main.go) directory for complete usage examples.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues, questions, or contributions, please open an issue on GitHub.

## Acknowledgments

Built with ❤️ for the Go community.
