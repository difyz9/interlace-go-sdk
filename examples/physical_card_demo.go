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

	// Example card and cardholder IDs
	cardID1 := "card-id-1"
	cardID2 := "card-id-2"
	cardholderID := "cardholder-id"

	fmt.Println("=== Physical Card Management Demo ===\n")

	// 1. List Physical Card Fees
	fmt.Println("1. Listing Physical Card Fees...")
	fees, err := client.PhysicalCard.ListPhysicalCardFees(ctx)
	if err != nil {
		log.Printf("Error listing physical card fees: %v\n", err)
	} else {
		fmt.Printf("Found %d fee structures:\n", len(fees))
		for _, fee := range fees {
			fmt.Printf("  - BIN: %s, Brand: %s, Currency: %s\n", fee.BinID, fee.CardBrand, fee.Currency)
			fmt.Printf("    Shipping Fee: %.2f, Production Fee: %.2f, Total: %.2f\n",
				fee.ShippingFee, fee.CardProductionFee, fee.TotalFee)
			fmt.Printf("    Estimated Delivery: %s\n", fee.EstimatedDelivery)
		}
	}

	// 2. Bulk Ship Physical Cards
	fmt.Println("\n2. Bulk Shipping Physical Cards...")
	shipRequest := &interlace.BulkShipRequest{
		CardIDs: []string{cardID1, cardID2},
		ShippingAddress: &interlace.ShippingAddress{
			RecipientName: "John Doe",
			AddressLine1:  "123 Main Street",
			AddressLine2:  "Apt 4B",
			City:          "New York",
			State:         "NY",
			PostalCode:    "10001",
			Country:       "US",
			PhoneNumber:   "+1-555-0123",
		},
		ShippingMethod: "EXPRESS",
		Notes:          "Please deliver during business hours",
	}

	shipResponse, err := client.PhysicalCard.BulkShipPhysicalCards(ctx, shipRequest)
	if err != nil {
		log.Printf("Error shipping physical cards: %v\n", err)
	} else {
		fmt.Printf("Shipment created successfully!\n")
		fmt.Printf("  Shipment ID: %s\n", shipResponse.ShipmentID)
		fmt.Printf("  Tracking Number: %s\n", shipResponse.TrackingNumber)
		fmt.Printf("  Status: %s\n", shipResponse.Status)
		fmt.Printf("  Cards: %v\n", shipResponse.CardIDs)
		fmt.Printf("  Estimated Delivery: %s\n", shipResponse.EstimatedDelivery)
		fmt.Printf("  Created At: %s\n", shipResponse.CreatedAt)
	}

	// 3. Generate Cardholder Identity URL
	fmt.Println("\n3. Generating Cardholder Identity Verification URL...")
	identityURL, err := client.PhysicalCard.GenerateCardholderIdentityURL(ctx, cardholderID)
	if err != nil {
		log.Printf("Error generating identity URL: %v\n", err)
	} else {
		fmt.Printf("Identity verification URL generated!\n")
		fmt.Printf("  Cardholder ID: %s\n", identityURL.CardholderID)
		fmt.Printf("  Verification ID: %s\n", identityURL.VerificationID)
		fmt.Printf("  Identity URL: %s\n", identityURL.IdentityURL)
		fmt.Printf("  Expires At: %s\n", identityURL.ExpiresAt)
		fmt.Printf("  Status: %s\n", identityURL.Status)
		fmt.Println("\n  Send this URL to the cardholder for identity verification.")
	}

	// 4. Confirm Cardholder Identity
	fmt.Println("\n4. Confirming Cardholder Identity...")
	confirmRequest := &interlace.ConfirmCardholderIdentityRequest{
		CardholderID:   cardholderID,
		VerificationID: "verification-id-from-previous-step",
		Verified:       true,
		Notes:          "Identity documents verified successfully",
	}

	confirmResponse, err := client.PhysicalCard.ConfirmCardholderIdentity(ctx, confirmRequest)
	if err != nil {
		log.Printf("Error confirming identity: %v\n", err)
	} else {
		fmt.Printf("Identity confirmed successfully!\n")
		fmt.Printf("  Cardholder ID: %s\n", confirmResponse.CardholderID)
		fmt.Printf("  Status: %s\n", confirmResponse.Status)
		fmt.Printf("  Verified At: %s\n", confirmResponse.VerifiedAt)
		fmt.Printf("  Message: %s\n", confirmResponse.Message)
	}

	// 5. Activate Physical Card
	fmt.Println("\n5. Activating Physical Card...")
	activateRequest := &interlace.ActivatePhysicalCardRequest{
		CardID:         cardID1,
		LastFourDigits: "1234",
		CVV:            "123",
		ExpiryMonth:    "12",
		ExpiryYear:     "2025",
	}

	activateResponse, err := client.PhysicalCard.ActivatePhysicalCard(ctx, activateRequest)
	if err != nil {
		log.Printf("Error activating physical card: %v\n", err)
	} else {
		fmt.Printf("Physical card activated successfully!\n")
		fmt.Printf("  Card ID: %s\n", activateResponse.CardID)
		fmt.Printf("  Status: %s\n", activateResponse.Status)
		fmt.Printf("  Activated At: %s\n", activateResponse.ActivatedAt)
		fmt.Printf("  Message: %s\n", activateResponse.Message)
	}

	fmt.Println("\n=== Physical Card Workflow Complete ===")
	fmt.Println("\nTypical Physical Card Flow:")
	fmt.Println("1. Check fees for different card types")
	fmt.Println("2. Create virtual cards (using Card API)")
	fmt.Println("3. Generate identity verification URL for cardholder")
	fmt.Println("4. Cardholder completes verification")
	fmt.Println("5. Confirm cardholder identity")
	fmt.Println("6. Bulk ship physical cards to address")
	fmt.Println("7. Cardholder receives card and activates it")
	fmt.Println("8. Card is ready for use!")
}
