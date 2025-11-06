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

	// Test account ID
	accountID := "your-account-id-here"
	
	// 1. Create a Budget
	fmt.Println("\n=== Create Budget ===")
	createReq := &interlace.CreateBudgetRequest{
		AccountID:   accountID,
		Name:        "Marketing Budget 2024",
		Currency:    "USD",
		Description: "Budget for marketing team expenses",
		InitBalance: 10000.00,
	}
	budget, err := client.Budget.CreateBudget(ctx, createReq)
	if err != nil {
		log.Printf("Failed to create budget: %v", err)
	} else {
		fmt.Printf("Budget created successfully:\n")
		fmt.Printf("  ID: %s\n", budget.ID)
		fmt.Printf("  Name: %s\n", budget.Name)
		fmt.Printf("  Currency: %s\n", budget.Currency)
		fmt.Printf("  Balance: %.2f\n", budget.Balance)
		fmt.Printf("  Status: %s\n", budget.Status)
	}

	// For demonstration, we'll use a sample budget ID
	budgetID := "sample-budget-id"
	if budget != nil {
		budgetID = budget.ID
	}

	// 2. List Budgets
	fmt.Println("\n=== List Budgets ===")
	listOptions := &interlace.ListBudgetsOptions{
		AccountID: accountID,
		Limit:     10,
		Page:      1,
	}
	budgetList, err := client.Budget.ListBudgets(ctx, listOptions)
	if err != nil {
		log.Printf("Failed to list budgets: %v", err)
	} else {
		fmt.Printf("Total budgets: %d\n", budgetList.Total)
		for i, b := range budgetList.List {
			fmt.Printf("  %d. %s (ID: %s, Balance: %.2f %s)\n", 
				i+1, b.Name, b.ID, b.Balance, b.Currency)
		}
	}

	// 3. Get Budget Details
	fmt.Println("\n=== Get Budget Details ===")
	budgetDetail, err := client.Budget.GetBudget(ctx, budgetID)
	if err != nil {
		log.Printf("Failed to get budget: %v", err)
	} else {
		fmt.Printf("Budget Details:\n")
		fmt.Printf("  ID: %s\n", budgetDetail.ID)
		fmt.Printf("  Name: %s\n", budgetDetail.Name)
		fmt.Printf("  Balance: %.2f %s\n", budgetDetail.Balance, budgetDetail.Currency)
		fmt.Printf("  Available Balance: %.2f\n", budgetDetail.AvailableBalance)
		fmt.Printf("  Pending Balance: %.2f\n", budgetDetail.PendingBalance)
		fmt.Printf("  Card Count: %d\n", budgetDetail.CardCount)
		fmt.Printf("  Status: %s\n", budgetDetail.Status)
	}

	// 4. Update Budget
	fmt.Println("\n=== Update Budget ===")
	updateReq := &interlace.UpdateBudgetRequest{
		Name:        "Marketing Budget Q1 2024",
		Description: "Updated marketing budget for Q1",
	}
	updatedBudget, err := client.Budget.UpdateBudget(ctx, budgetID, updateReq)
	if err != nil {
		log.Printf("Failed to update budget: %v", err)
	} else {
		fmt.Printf("Budget updated successfully:\n")
		fmt.Printf("  New Name: %s\n", updatedBudget.Name)
		fmt.Printf("  Description: %s\n", updatedBudget.Description)
	}

	// 5. Increase Budget Balance (Top-up)
	fmt.Println("\n=== Increase Budget Balance ===")
	increaseReq := &interlace.IncreaseBudgetBalanceRequest{
		Amount:          5000.00,
		Currency:        "USD",
		MerchantTradeNo: fmt.Sprintf("TOPUP-%d", time.Now().Unix()),
		Description:     "Additional funding for Q1",
	}
	increaseResp, err := client.Budget.IncreaseBudgetBalance(ctx, budgetID, increaseReq)
	if err != nil {
		log.Printf("Failed to increase budget balance: %v", err)
	} else {
		fmt.Printf("Budget balance increased:\n")
		fmt.Printf("  Transaction ID: %s\n", increaseResp.ID)
		fmt.Printf("  Amount: %.2f %s\n", increaseResp.Amount, increaseResp.Currency)
		fmt.Printf("  Balance Before: %.2f\n", increaseResp.BalanceBefore)
		fmt.Printf("  Balance After: %.2f\n", increaseResp.BalanceAfter)
		fmt.Printf("  Status: %s\n", increaseResp.Status)
	}

	// 6. Decrease Budget Balance (Withdraw)
	fmt.Println("\n=== Decrease Budget Balance ===")
	decreaseReq := &interlace.DecreaseBudgetBalanceRequest{
		Amount:          1000.00,
		Currency:        "USD",
		MerchantTradeNo: fmt.Sprintf("WITHDRAW-%d", time.Now().Unix()),
		Description:     "Return unused funds",
	}
	decreaseResp, err := client.Budget.DecreaseBudgetBalance(ctx, budgetID, decreaseReq)
	if err != nil {
		log.Printf("Failed to decrease budget balance: %v", err)
	} else {
		fmt.Printf("Budget balance decreased:\n")
		fmt.Printf("  Transaction ID: %s\n", decreaseResp.ID)
		fmt.Printf("  Amount: %.2f %s\n", decreaseResp.Amount, decreaseResp.Currency)
		fmt.Printf("  Balance Before: %.2f\n", decreaseResp.BalanceBefore)
		fmt.Printf("  Balance After: %.2f\n", decreaseResp.BalanceAfter)
		fmt.Printf("  Status: %s\n", decreaseResp.Status)
	}

	// 7. List Budget Transactions
	fmt.Println("\n=== List Budget Transactions ===")
	txOptions := &interlace.ListBudgetTransactionsOptions{
		Limit: 10,
		Page:  1,
	}
	txList, err := client.Budget.ListBudgetTransactions(ctx, budgetID, txOptions)
	if err != nil {
		log.Printf("Failed to list budget transactions: %v", err)
	} else {
		fmt.Printf("Total transactions: %d\n", txList.Total)
		for i, tx := range txList.List {
			fmt.Printf("  %d. Type: %s, Amount: %.2f %s, Status: %s\n", 
				i+1, tx.Type, tx.Amount, tx.Currency, tx.Status)
			fmt.Printf("      Balance: %.2f -> %.2f\n", tx.BalanceBefore, tx.BalanceAfter)
		}
	}

	// 8. Get Budget Transaction Details
	fmt.Println("\n=== Get Budget Transaction Details ===")
	if txList != nil && len(txList.List) > 0 {
		transactionID := txList.List[0].ID
		txDetail, err := client.Budget.GetBudgetTransaction(ctx, budgetID, transactionID)
		if err != nil {
			log.Printf("Failed to get transaction details: %v", err)
		} else {
			fmt.Printf("Transaction Details:\n")
			fmt.Printf("  ID: %s\n", txDetail.ID)
			fmt.Printf("  Type: %s\n", txDetail.Type)
			fmt.Printf("  Amount: %.2f %s\n", txDetail.Amount, txDetail.Currency)
			fmt.Printf("  Status: %s\n", txDetail.Status)
			fmt.Printf("  Description: %s\n", txDetail.Description)
			fmt.Printf("  Created At: %s\n", txDetail.CreatedAt)
		}
	}

	// 9. Delete Budget (commented out to avoid accidental deletion)
	fmt.Println("\n=== Delete Budget (Commented Out) ===")
	fmt.Println("Budget deletion is commented out in the demo code.")
	/*
	deleteResp, err := client.Budget.DeleteBudget(ctx, budgetID)
	if err != nil {
		log.Printf("Failed to delete budget: %v", err)
	} else {
		fmt.Printf("Budget deleted successfully:\n")
		fmt.Printf("  ID: %s\n", deleteResp.ID)
		fmt.Printf("  Message: %s\n", deleteResp.Message)
		fmt.Printf("  Success: %v\n", deleteResp.Success)
	}
	*/

	fmt.Println("\n=== Budget Management Demo Completed ===")
}
