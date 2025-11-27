package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/erfandiakoo/go-signalads"
)

func main() {
	apiKey := os.Getenv("SIGNALADS_API_KEY")
	apiSecret := os.Getenv("SIGNALADS_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatal("SIGNALADS_API_KEY and SIGNALADS_API_SECRET environment variables must be set")
	}

	client := signalads.NewClient(apiKey, apiSecret)
	ctx := context.Background()

	phoneNumber := os.Getenv("TEST_PHONE_NUMBER")
	if phoneNumber == "" {
		phoneNumber = "+989123456789"
		fmt.Println("Using default phone number. Set TEST_PHONE_NUMBER env var to use your number.")
	}

	// Example 1: Send simple message
	fmt.Println("=== Send Simple Message ===")
	msgResponse, err := client.Messages.SendMessage(ctx, phoneNumber, "Hello from SignalAds Go client!")
	if err != nil {
		log.Printf("Error sending message: %v\n", err)
	} else {
		fmt.Printf("Message sent successfully! ID: %s, Status: %s\n", msgResponse.ID, msgResponse.Status)
	}

	// Example 2: Send message with document link
	fmt.Println("\n=== Send Message with Document ===")
	docResponse, err := client.Messages.SendMessageWithDocument(
		ctx,
		phoneNumber,
		"Please check this important document",
		"https://example.com/document.pdf",
		"Important Document",
	)
	if err != nil {
		log.Printf("Error sending message with document: %v\n", err)
	} else {
		fmt.Printf("Message with document sent successfully! ID: %s\n", docResponse.ID)
	}

	// Example 3: Send bulk messages
	fmt.Println("\n=== Send Bulk Messages ===")
	bulkMessages := []signalads.BulkMessageItem{
		{To: phoneNumber, Message: "Message 1"},
		{To: phoneNumber, Message: "Message 2"},
	}
	bulkResponse, err := client.Messages.SendBulkMessage(ctx, bulkMessages, "")
	if err != nil {
		log.Printf("Error sending bulk messages: %v\n", err)
	} else {
		fmt.Printf("Bulk messages sent! Total: %d, Success: %d, Failed: %d\n",
			bulkResponse.Total, bulkResponse.Success, bulkResponse.Failed)
	}

	// Example 4: Send template message
	fmt.Println("\n=== Send Template Message ===")
	templateResponse, err := client.Messages.SendTemplate(
		ctx,
		phoneNumber,
		"template_123",
		map[string]string{
			"name": "Erfan",
			"code": "12345",
		},
	)
	if err != nil {
		log.Printf("Error sending template message: %v\n", err)
	} else {
		fmt.Printf("Template message sent successfully! ID: %s\n", templateResponse.ID)
	}

	// Example 5: Send voice message
	fmt.Println("\n=== Send Voice Message ===")
	voiceResponse, err := client.Messages.SendVoice(
		ctx,
		phoneNumber,
		"This is a test voice message",
		"female",
		"fa",
	)
	if err != nil {
		log.Printf("Error sending voice message: %v\n", err)
	} else {
		fmt.Printf("Voice message sent successfully! ID: %s\n", voiceResponse.ID)
	}

	// Example 6: List messages
	fmt.Println("\n=== List Messages ===")
	messagesList, err := client.Messages.ListMessages(ctx, &signalads.PaginationParams{
		Page:    1,
		PerPage: 10,
	})
	if err != nil {
		log.Printf("Error listing messages: %v\n", err)
	} else {
		fmt.Printf("Found %d messages (Total: %d)\n", len(messagesList.Messages), messagesList.Total)
		for i, msg := range messagesList.Messages {
			if i >= 5 {
				break
			}
			fmt.Printf("  - ID: %s, To: %s, Status: %s\n", msg.ID, msg.To, msg.Status)
		}
	}

	// Example 7: Get message status
	if msgResponse != nil && msgResponse.ID != "" {
		fmt.Println("\n=== Get Message Status ===")
		status, err := client.Messages.GetMessageStatus(ctx, msgResponse.ID)
		if err != nil {
			log.Printf("Error getting message status: %v\n", err)
		} else {
			fmt.Printf("Message status: %s\n", status.Status)
			if status.Cost > 0 {
				fmt.Printf("Cost: %.2f\n", status.Cost)
			}
		}
	}

	// Example 8: Get user info
	fmt.Println("\n=== Get User Info ===")
	userInfo, err := client.Messages.GetUserInfo(ctx)
	if err != nil {
		log.Printf("Error getting user info: %v\n", err)
	} else {
		fmt.Printf("User ID: %s\n", userInfo.ID)
		if userInfo.Balance > 0 {
			fmt.Printf("Balance: %.2f\n", userInfo.Balance)
		}
		if userInfo.Credit > 0 {
			fmt.Printf("Credit: %.2f\n", userInfo.Credit)
		}
	}
}
