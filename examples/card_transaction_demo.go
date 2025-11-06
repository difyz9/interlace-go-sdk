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

	// Test IDs (replace with real ones)
	accountID := "your-account-id-here"
	cardID := "your-card-id"

	// 1. Card Transfer In (从 Quantum 账户向卡转入资金)
	fmt.Println("\n=== Card Transfer In ===")
	transferInReq := &interlace.CardTransferInRequest{
		CardID:              cardID,
		Amount:              1000.00,
		Currency:            "USD",
		MerchantTradeNo:     fmt.Sprintf("TIN-%d", time.Now().Unix()),
		InfinityAccountID:   accountID,
		InfinityAccountType: "MAIN",
	}
	transferInResp, err := client.CardTransaction.CardTransferIn(ctx, transferInReq)
	if err != nil {
		log.Printf("Failed to transfer in: %v", err)
	} else {
		fmt.Printf("Transfer in successful:\n")
		fmt.Printf("  Transaction ID: %s\n", transferInResp.ID)
		fmt.Printf("  Card ID: %s\n", transferInResp.CardID)
		fmt.Printf("  Amount: %.2f %s\n", transferInResp.Amount, transferInResp.Currency)
		fmt.Printf("  Status: %s\n", transferInResp.Status)
		fmt.Printf("  Merchant Trade No: %s\n", transferInResp.MerchantTradeNo)
		fmt.Printf("  Created At: %s\n", transferInResp.CreatedAt)
	}

	// 2. Card Transfer Out (从卡向 Quantum 账户转出资金)
	fmt.Println("\n=== Card Transfer Out ===")
	transferOutReq := &interlace.CardTransferOutRequest{
		CardID:              cardID,
		Amount:              500.00,
		Currency:            "USD",
		MerchantTradeNo:     fmt.Sprintf("TOUT-%d", time.Now().Unix()),
		InfinityAccountID:   accountID,
		InfinityAccountType: "MAIN",
	}
	transferOutResp, err := client.CardTransaction.CardTransferOut(ctx, transferOutReq)
	if err != nil {
		log.Printf("Failed to transfer out: %v", err)
	} else {
		fmt.Printf("Transfer out successful:\n")
		fmt.Printf("  Transaction ID: %s\n", transferOutResp.ID)
		fmt.Printf("  Card ID: %s\n", transferOutResp.CardID)
		fmt.Printf("  Amount: %.2f %s\n", transferOutResp.Amount, transferOutResp.Currency)
		fmt.Printf("  Status: %s\n", transferOutResp.Status)
		fmt.Printf("  Merchant Trade No: %s\n", transferOutResp.MerchantTradeNo)
		fmt.Printf("  Created At: %s\n", transferOutResp.CreatedAt)
	}

	// 3. List Card Transactions (查询所有卡交易)
	fmt.Println("\n=== List Card Transactions (All) ===")
	allTxOptions := &interlace.ListCardTransactionsOptions{
		AccountID: accountID,
		Limit:     20,
		Page:      1,
	}
	allTxList, err := client.CardTransaction.ListCardTransactions(ctx, allTxOptions)
	if err != nil {
		log.Printf("Failed to list transactions: %v", err)
	} else {
		fmt.Printf("Total transactions: %d\n", allTxList.Total)
		fmt.Printf("Showing page %d (limit: %d)\n", allTxList.Page, allTxList.Limit)
		for i, tx := range allTxList.List {
			fmt.Printf("\n  %d. Transaction ID: %s\n", i+1, tx.ID)
			fmt.Printf("     Type: %s\n", tx.Type)
			fmt.Printf("     Amount: %.2f %s\n", tx.Amount, tx.Currency)
			fmt.Printf("     Status: %s\n", tx.Status)
			if tx.MerchantName != "" {
				fmt.Printf("     Merchant: %s (%s)\n", tx.MerchantName, tx.MerchantCountry)
			}
			fmt.Printf("     Transaction Time: %s\n", tx.TransactionTime)
		}
	}

	// 4. List Card Transactions by Card ID
	fmt.Println("\n=== List Card Transactions (By Card) ===")
	cardTxOptions := &interlace.ListCardTransactionsOptions{
		CardID: cardID,
		Limit:  10,
		Page:   1,
	}
	cardTxList, err := client.CardTransaction.ListCardTransactions(ctx, cardTxOptions)
	if err != nil {
		log.Printf("Failed to list card transactions: %v", err)
	} else {
		fmt.Printf("Total transactions for card %s: %d\n", cardID, cardTxList.Total)
		for i, tx := range cardTxList.List {
			fmt.Printf("  %d. %s - %.2f %s (%s)\n", 
				i+1, tx.Type, tx.Amount, tx.Currency, tx.Status)
		}
	}

	// 5. Filter by Transaction Type
	fmt.Println("\n=== List Purchase Transactions ===")
	purchaseTxOptions := &interlace.ListCardTransactionsOptions{
		AccountID: accountID,
		Type:      "PURCHASE",
		Status:    "APPROVED",
		Limit:     10,
		Page:      1,
	}
	purchaseTxList, err := client.CardTransaction.ListCardTransactions(ctx, purchaseTxOptions)
	if err != nil {
		log.Printf("Failed to list purchase transactions: %v", err)
	} else {
		fmt.Printf("Total purchase transactions: %d\n", purchaseTxList.Total)
		for i, tx := range purchaseTxList.List {
			fmt.Printf("\n  %d. Purchase at %s\n", i+1, tx.MerchantName)
			fmt.Printf("     Amount: %.2f %s\n", tx.Amount, tx.Currency)
			fmt.Printf("     Settlement: %.2f %s\n", tx.SettlementAmount, tx.SettlementCurrency)
			if tx.BillingCurrency != tx.Currency {
				fmt.Printf("     Billing: %.2f %s (Rate: %.4f)\n", 
					tx.BillingAmount, tx.BillingCurrency, tx.ExchangeRate)
			}
			fmt.Printf("     MCC: %s\n", tx.MerchantCategoryCode)
			fmt.Printf("     Country: %s\n", tx.MerchantCountry)
			fmt.Printf("     International: %v, Online: %v\n", tx.IsInternational, tx.IsOnline)
			if tx.AuthorizationCode != "" {
				fmt.Printf("     Auth Code: %s\n", tx.AuthorizationCode)
			}
		}
	}

	// 6. Filter by Date Range
	fmt.Println("\n=== List Transactions by Date Range ===")
	startDate := time.Now().AddDate(0, -1, 0).Format(time.RFC3339) // Last month
	endDate := time.Now().Format(time.RFC3339)
	
	dateTxOptions := &interlace.ListCardTransactionsOptions{
		AccountID: accountID,
		StartTime: startDate,
		EndTime:   endDate,
		Limit:     15,
		Page:      1,
	}
	dateTxList, err := client.CardTransaction.ListCardTransactions(ctx, dateTxOptions)
	if err != nil {
		log.Printf("Failed to list transactions by date: %v", err)
	} else {
		fmt.Printf("Transactions from %s to %s: %d\n", 
			startDate, endDate, dateTxList.Total)
		
		// Group by type
		typeCount := make(map[string]int)
		var totalAmount float64
		for _, tx := range dateTxList.List {
			typeCount[tx.Type]++
			if tx.Status == "APPROVED" {
				totalAmount += tx.Amount
			}
		}
		
		fmt.Println("\nSummary:")
		for txType, count := range typeCount {
			fmt.Printf("  %s: %d transactions\n", txType, count)
		}
		fmt.Printf("  Total approved amount: %.2f\n", totalAmount)
	}

	// 7. List Declined Transactions
	fmt.Println("\n=== List Declined Transactions ===")
	declinedTxOptions := &interlace.ListCardTransactionsOptions{
		AccountID: accountID,
		Status:    "DECLINED",
		Limit:     10,
		Page:      1,
	}
	declinedTxList, err := client.CardTransaction.ListCardTransactions(ctx, declinedTxOptions)
	if err != nil {
		log.Printf("Failed to list declined transactions: %v", err)
	} else {
		fmt.Printf("Total declined transactions: %d\n", declinedTxList.Total)
		for i, tx := range declinedTxList.List {
			fmt.Printf("  %d. Type: %s, Amount: %.2f %s\n", 
				i+1, tx.Type, tx.Amount, tx.Currency)
			if tx.DeclineReason != "" {
				fmt.Printf("     Decline Reason: %s\n", tx.DeclineReason)
			}
			fmt.Printf("     Time: %s\n", tx.TransactionTime)
		}
	}

	fmt.Println("\n=== Card Transaction Demo Completed ===")
	fmt.Println("\nNote: Replace test IDs with real card IDs to see actual transaction data.")
	fmt.Println("Transaction types may include: PURCHASE, REFUND, TRANSFER_IN, TRANSFER_OUT, etc.")
	fmt.Println("Transaction statuses may include: APPROVED, DECLINED, PENDING, SETTLED, etc.")
}
