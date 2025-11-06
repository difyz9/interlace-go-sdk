package main

import (
	"context"
	"fmt"
	"log"
	"time"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	// Initialize client
	config := interlace.DefaultConfig()
	client := interlace.NewClient(config)

	// Set your access token
	accessToken := "your-access-token-here"
	client.SetAccessToken(accessToken)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Example card IDs
	cardID := "your-card-id"

	fmt.Println("=== Card Security Demo ===\n")

	// Update Card PIN
	fmt.Println("1. Updating Card PIN...")
	updatePINRequest := &interlace.UpdatePINRequest{
		CardID:     cardID,
		NewPIN:     "1234",
		ConfirmPIN: "1234",
	}

	response, err := client.Security.UpdateCardPIN(ctx, updatePINRequest)
	if err != nil {
		log.Printf("Error updating PIN: %v\n", err)
	} else {
		fmt.Printf("PIN updated successfully!\n")
		fmt.Printf("  Card ID: %s\n", response.CardID)
		fmt.Printf("  Status: %s\n", response.Status)
		fmt.Printf("  Updated At: %s\n", response.UpdatedAt)
		fmt.Printf("  Message: %s\n", response.Message)
	}

	// PIN validation examples
	fmt.Println("\n=== PIN Requirements ===")
	fmt.Println("✓ PIN must be 4-6 digits")
	fmt.Println("✓ NewPIN and ConfirmPIN must match")
	fmt.Println("✓ Card must be active")

	// Example: Invalid PIN (too short)
	fmt.Println("\n2. Testing Invalid PIN (too short)...")
	invalidRequest1 := &interlace.UpdatePINRequest{
		CardID:     cardID,
		NewPIN:     "123", // Only 3 digits
		ConfirmPIN: "123",
	}
	_, err = client.Security.UpdateCardPIN(ctx, invalidRequest1)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// Example: Mismatched PINs
	fmt.Println("\n3. Testing Mismatched PINs...")
	invalidRequest2 := &interlace.UpdatePINRequest{
		CardID:     cardID,
		NewPIN:     "1234",
		ConfirmPIN: "5678", // Different
	}
	_, err = client.Security.UpdateCardPIN(ctx, invalidRequest2)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// Best practices
	fmt.Println("\n=== PIN Security Best Practices ===")
	fmt.Println("1. Never store PINs in plain text")
	fmt.Println("2. Implement rate limiting for PIN changes")
	fmt.Println("3. Require additional authentication for PIN updates")
	fmt.Println("4. Notify cardholder via email/SMS after PIN change")
	fmt.Println("5. Use strong PINs (avoid common patterns like 1234, 0000)")
	fmt.Println("6. Consider implementing PIN retry limits")
	fmt.Println("7. Log all PIN change attempts for security auditing")

	fmt.Println("\n=== Security Operations Complete ===")
}
