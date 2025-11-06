package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	// Initialize client
	clientID := os.Getenv("INTERLACE_CLIENT_ID")
	clientSecret := os.Getenv("INTERLACE_CLIENT_SECRET")
	
	if clientID == "" {
		clientID = "your-client-id-here"
	}
	if clientSecret == "" {
		clientSecret = "your-client-secret"
	}

	client := interlace.NewClient(clientID, clientSecret, true)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test account and card IDs (replace with real ones)
	accountID := "your-account-id-here"
	cardID := "your-card-id"
	binID := "your-bin-id"
	cardholderID := "your-cardholder-id"
	budgetID := "your-budget-id"
	walletID := "your-wallet-id"

	// 1. Freeze Card
	fmt.Println("\n=== Freeze Card ===")
	frozenCard, err := client.Card.FreezeCard(ctx, cardID)
	if err != nil {
		log.Printf("Failed to freeze card: %v", err)
	} else {
		fmt.Printf("Card frozen successfully: %s (Status: %s)\n", frozenCard.ID, frozenCard.Status)
	}

	// 2. Unfreeze Card
	fmt.Println("\n=== Unfreeze Card ===")
	unfrozenCard, err := client.Card.UnfreezeCard(ctx, cardID)
	if err != nil {
		log.Printf("Failed to unfreeze card: %v", err)
	} else {
		fmt.Printf("Card unfrozen successfully: %s (Status: %s)\n", unfrozenCard.ID, unfrozenCard.Status)
	}

	// 3. Set Velocity Control
	fmt.Println("\n=== Set Velocity Control ===")
	dailyLimit := 1000.00
	singleLimit := 500.00
	velocityReq := &interlace.VelocityControlRequest{
		DailySpendingLimit: &dailyLimit,
		SingleTransLimit:   &singleLimit,
	}
	controlledCard, err := client.Card.SetCardVelocityControl(ctx, cardID, velocityReq)
	if err != nil {
		log.Printf("Failed to set velocity control: %v", err)
	} else {
		fmt.Printf("Velocity control set successfully for card: %s\n", controlledCard.ID)
		fmt.Printf("Daily Limit: %.2f, Single Transaction Limit: %.2f\n", 
			controlledCard.DailySpendingLimit, controlledCard.SingleTransactionLimit)
	}

	// 4. Create Prepaid Card
	fmt.Println("\n=== Create Prepaid Card ===")
	prepaidReq := &interlace.CreatePrepaidCardRequest{
		BinID:                    binID,
		CardholderID:             cardholderID,
		Label:                    "Demo Prepaid Card",
		DailySpendingLimit:       500.00,
		SingleTransLimit:         100.00,
		MonthlySpendingLimit:     5000.00,
		ThreeDSecureAuthRequired: true,
	}
	prepaidCard, err := client.Card.CreatePrepaidCard(ctx, prepaidReq)
	if err != nil {
		log.Printf("Failed to create prepaid card: %v", err)
	} else {
		fmt.Printf("Prepaid card created successfully: %s (Type: %s)\n", 
			prepaidCard.ID, prepaidCard.Type)
	}

	// 5. Batch Create Prepaid Cards
	fmt.Println("\n=== Batch Create Prepaid Cards ===")
	batchPrepaidCards := []interlace.CreatePrepaidCardRequest{
		{
			BinID:                binID,
			CardholderID:         cardholderID,
			Label:                "Batch Prepaid Card 1",
			DailySpendingLimit:   300.00,
			SingleTransLimit:     50.00,
		},
		{
			BinID:                binID,
			CardholderID:         cardholderID,
			Label:                "Batch Prepaid Card 2",
			DailySpendingLimit:   400.00,
			SingleTransLimit:     75.00,
		},
	}
	batchPrepaidResp, err := client.Card.BatchCreatePrepaidCards(ctx, batchPrepaidCards)
	if err != nil {
		log.Printf("Failed to batch create prepaid cards: %v", err)
	} else {
		fmt.Printf("Batch created %d prepaid cards successfully\n", batchPrepaidResp.Total)
		for i, card := range batchPrepaidResp.List {
			fmt.Printf("  %d. Card ID: %s, Label: %s\n", i+1, card.ID, card.Label)
		}
	}

	// 6. Create Budget Card
	fmt.Println("\n=== Create Budget Card ===")
	budgetCardReq := &interlace.CreateBudgetCardRequest{
		BinID:                    binID,
		CardholderID:             cardholderID,
		BudgetID:                 budgetID,
		Label:                    "Demo Budget Card",
		DailySpendingLimit:       200.00,
		SingleTransLimit:         50.00,
		MonthlySpendingLimit:     2000.00,
		ThreeDSecureAuthRequired: true,
	}
	budgetCard, err := client.Card.CreateBudgetCard(ctx, budgetCardReq)
	if err != nil {
		log.Printf("Failed to create budget card: %v", err)
	} else {
		fmt.Printf("Budget card created successfully: %s (Type: %s)\n", 
			budgetCard.ID, budgetCard.Type)
	}

	// 7. Batch Create Budget Cards
	fmt.Println("\n=== Batch Create Budget Cards ===")
	batchBudgetCards := []interlace.CreateBudgetCardRequest{
		{
			BinID:              binID,
			CardholderID:       cardholderID,
			BudgetID:           budgetID,
			Label:              "Batch Budget Card 1",
			DailySpendingLimit: 150.00,
			SingleTransLimit:   30.00,
		},
		{
			BinID:              binID,
			CardholderID:       cardholderID,
			BudgetID:           budgetID,
			Label:              "Batch Budget Card 2",
			DailySpendingLimit: 250.00,
			SingleTransLimit:   60.00,
		},
	}
	batchBudgetResp, err := client.Card.BatchCreateBudgetCards(ctx, batchBudgetCards)
	if err != nil {
		log.Printf("Failed to batch create budget cards: %v", err)
	} else {
		fmt.Printf("Batch created %d budget cards successfully\n", batchBudgetResp.Total)
		for i, card := range batchBudgetResp.List {
			fmt.Printf("  %d. Card ID: %s, Label: %s\n", i+1, card.ID, card.Label)
		}
	}

	// 8. Get Card Summary
	fmt.Println("\n=== Get Card Summary ===")
	summary, err := client.Card.GetCardSummary(ctx, cardID)
	if err != nil {
		log.Printf("Failed to get card summary: %v", err)
	} else {
		fmt.Printf("Card Summary for %s:\n", summary.CardID)
		fmt.Printf("  Available Balance: %.2f\n", summary.AvailableBalance)
		fmt.Printf("  Current Balance: %.2f\n", summary.CurrentBalance)
		fmt.Printf("  Pending Transactions: %.2f\n", summary.PendingTransactions)
		fmt.Printf("  Daily Spending Limit: %.2f (Spent Today: %.2f)\n", 
			summary.DailySpendingLimit, summary.SpentToday)
		fmt.Printf("  Monthly Spending Limit: %.2f (Spent This Month: %.2f)\n", 
			summary.MonthlySpendingLimit, summary.SpentThisMonth)
		fmt.Printf("  Single Transaction Limit: %.2f\n", summary.SingleTransLimit)
	}

	// 9. Update Card
	fmt.Println("\n=== Update Card ===")
	newDailyLimit := 800.00
	newLabel := "Updated Card Label"
	updateReq := &interlace.UpdateCardRequest{
		CardID:             cardID,
		Label:              newLabel,
		DailySpendingLimit: &newDailyLimit,
	}
	updatedCard, err := client.Card.UpdateCard(ctx, updateReq)
	if err != nil {
		log.Printf("Failed to update card: %v", err)
	} else {
		fmt.Printf("Card updated successfully: %s\n", updatedCard.ID)
		fmt.Printf("  New Label: %s\n", updatedCard.Label)
		fmt.Printf("  New Daily Spending Limit: %.2f\n", updatedCard.DailySpendingLimit)
	}

	// 10. Bind Wallet
	fmt.Println("\n=== Bind Wallet ===")
	bindReq := &interlace.BindWalletRequest{
		WalletID: walletID,
	}
	walletBoundCard, err := client.Card.BindWallet(ctx, cardID, bindReq)
	if err != nil {
		log.Printf("Failed to bind wallet: %v", err)
	} else {
		fmt.Printf("Wallet bound successfully to card: %s\n", walletBoundCard.ID)
	}

	// List all cards to verify
	fmt.Println("\n=== List All Cards ===")
	listOptions := &interlace.CardListOptions{
		AccountID: accountID,
		Limit:     10,
		Page:      1,
	}
	cardList, err := client.Card.ListCards(ctx, listOptions)
	if err != nil {
		log.Printf("Failed to list cards: %v", err)
	} else {
		fmt.Printf("Total cards: %d\n", cardList.Total)
		for i, card := range cardList.List {
			fmt.Printf("  %d. ID: %s, Type: %s, Status: %s, Label: %s\n", 
				i+1, card.ID, card.Type, card.Status, card.Label)
		}
	}

	fmt.Println("\n=== Card Management Demo Completed ===")
}
