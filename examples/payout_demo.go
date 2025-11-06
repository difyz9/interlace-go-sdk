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

	accountID := "your-account-id-here"

	// 1. Get Exchange Rate
	fmt.Println("\n=== Get Exchange Rate ===")
	rate, err := client.Payout.GetExchangeRate(ctx, "USD", "EUR", 1000.00)
	if err != nil {
		log.Printf("Failed to get exchange rate: %v", err)
	} else {
		fmt.Printf("Exchange Rate USD to EUR:\n")
		fmt.Printf("  Rate: %.6f\n", rate.Rate)
		fmt.Printf("  Inverse Rate: %.6f\n", rate.InverseRate)
		fmt.Printf("  Amount: %.2f %s\n", rate.Amount, rate.SourceCurrency)
		fmt.Printf("  Converted Amount: %.2f %s\n", rate.ConvertedAmount, rate.TargetCurrency)
		fmt.Printf("  Valid Until: %s\n", rate.ValidUntil)
	}

	// 2. Create a Payee
	fmt.Println("\n=== Create Payee ===")
	createPayeeReq := &interlace.CreatePayeeRequest{
		AccountID:          accountID,
		BeneficiaryName:    "John Doe",
		BankName:           "Example Bank",
		BankCountry:        "US",
		AccountNumber:      "1234567890",
		SwiftCode:          "EXAMPUS33",
		RoutingNumber:      "021000021",
		Currency:           "USD",
		BeneficiaryType:    "INDIVIDUAL",
		BeneficiaryAddress: "123 Main St, New York, NY 10001",
	}
	payee, err := client.Payout.CreatePayee(ctx, createPayeeReq)
	if err != nil {
		log.Printf("Failed to create payee: %v", err)
	} else {
		fmt.Printf("Payee created successfully:\n")
		fmt.Printf("  ID: %s\n", payee.ID)
		fmt.Printf("  Beneficiary: %s\n", payee.BeneficiaryName)
		fmt.Printf("  Bank: %s (%s)\n", payee.BankName, payee.BankCountry)
		fmt.Printf("  Account: %s\n", payee.AccountNumber)
		fmt.Printf("  Currency: %s\n", payee.Currency)
		fmt.Printf("  Status: %s\n", payee.Status)
	}

	payeeID := "sample-payee-id"
	if payee != nil {
		payeeID = payee.ID
	}

	// 3. List Payees
	fmt.Println("\n=== List Payees ===")
	listPayeesOptions := &interlace.ListPayeesOptions{
		AccountID: accountID,
		Limit:     10,
		Page:      1,
	}
	payeeList, err := client.Payout.ListPayees(ctx, listPayeesOptions)
	if err != nil {
		log.Printf("Failed to list payees: %v", err)
	} else {
		fmt.Printf("Total payees: %d\n", payeeList.Total)
		for i, p := range payeeList.List {
			fmt.Printf("  %d. %s - %s (%s)\n", 
				i+1, p.BeneficiaryName, p.BankName, p.Currency)
		}
	}

	// 4. Get Payee Details
	fmt.Println("\n=== Get Payee Details ===")
	payeeDetail, err := client.Payout.GetPayee(ctx, payeeID)
	if err != nil {
		log.Printf("Failed to get payee: %v", err)
	} else {
		fmt.Printf("Payee Details:\n")
		fmt.Printf("  ID: %s\n", payeeDetail.ID)
		fmt.Printf("  Beneficiary: %s\n", payeeDetail.BeneficiaryName)
		fmt.Printf("  Bank: %s\n", payeeDetail.BankName)
		fmt.Printf("  Account Number: %s\n", payeeDetail.AccountNumber)
		fmt.Printf("  SWIFT: %s\n", payeeDetail.SwiftCode)
		fmt.Printf("  Currency: %s\n", payeeDetail.Currency)
	}

	// 5. Create Quotation
	fmt.Println("\n=== Create Quotation ===")
	createQuotationReq := &interlace.CreateQuotationRequest{
		AccountID:      accountID,
		SourceCurrency: "USD",
		SourceAmount:   1000.00,
		TargetCurrency: "EUR",
		PayeeID:        payeeID,
	}
	quotation, err := client.Payout.CreateQuotation(ctx, createQuotationReq)
	if err != nil {
		log.Printf("Failed to create quotation: %v", err)
	} else {
		fmt.Printf("Quotation created:\n")
		fmt.Printf("  ID: %s\n", quotation.ID)
		fmt.Printf("  Source: %.2f %s\n", quotation.SourceAmount, quotation.SourceCurrency)
		fmt.Printf("  Target: %.2f %s\n", quotation.TargetAmount, quotation.TargetCurrency)
		fmt.Printf("  Exchange Rate: %.6f\n", quotation.ExchangeRate)
		fmt.Printf("  Fee: %.2f\n", quotation.Fee)
		fmt.Printf("  Total Amount: %.2f\n", quotation.TotalAmount)
		fmt.Printf("  Valid Until: %s\n", quotation.ValidUntil)
		fmt.Printf("  Estimated Arrival: %s\n", quotation.EstimatedArrival)
	}

	quotationID := "sample-quotation-id"
	if quotation != nil {
		quotationID = quotation.ID
	}

	// 6. Get Quotation Details
	fmt.Println("\n=== Get Quotation Details ===")
	quotationDetail, err := client.Payout.GetQuotation(ctx, quotationID)
	if err != nil {
		log.Printf("Failed to get quotation: %v", err)
	} else {
		fmt.Printf("Quotation Details:\n")
		fmt.Printf("  ID: %s\n", quotationDetail.ID)
		fmt.Printf("  Rate: %.6f\n", quotationDetail.ExchangeRate)
		fmt.Printf("  Status: %s\n", quotationDetail.Status)
	}

	// 7. Accept Quotation (creates payout)
	fmt.Println("\n=== Accept Quotation ===")
	acceptReq := &interlace.AcceptQuotationRequest{
		PayeeID:         payeeID,
		MerchantTradeNo: fmt.Sprintf("PAYOUT-%d", time.Now().Unix()),
		Reference:       "Payment for services",
	}
	payout, err := client.Payout.AcceptQuotation(ctx, quotationID, acceptReq)
	if err != nil {
		log.Printf("Failed to accept quotation: %v", err)
	} else {
		fmt.Printf("Quotation accepted, payout created:\n")
		fmt.Printf("  Payout ID: %s\n", payout.ID)
		fmt.Printf("  Source: %.2f %s\n", payout.SourceAmount, payout.SourceCurrency)
		fmt.Printf("  Target: %.2f %s\n", payout.TargetAmount, payout.TargetCurrency)
		fmt.Printf("  Status: %s\n", payout.Status)
		fmt.Printf("  Estimated Arrival: %s\n", payout.EstimatedArrival)
	}

	// 8. Create Payout Directly (without quotation)
	fmt.Println("\n=== Create Payout Directly ===")
	createPayoutReq := &interlace.CreatePayoutRequest{
		AccountID:       accountID,
		PayeeID:         payeeID,
		SourceCurrency:  "USD",
		SourceAmount:    500.00,
		TargetCurrency:  "EUR",
		MerchantTradeNo: fmt.Sprintf("DIRECT-%d", time.Now().Unix()),
		Reference:       "Direct payment",
	}
	directPayout, err := client.Payout.CreatePayout(ctx, createPayoutReq)
	if err != nil {
		log.Printf("Failed to create payout: %v", err)
	} else {
		fmt.Printf("Payout created directly:\n")
		fmt.Printf("  ID: %s\n", directPayout.ID)
		fmt.Printf("  Beneficiary: %s\n", directPayout.BeneficiaryName)
		fmt.Printf("  Amount: %.2f %s -> %.2f %s\n", 
			directPayout.SourceAmount, directPayout.SourceCurrency,
			directPayout.TargetAmount, directPayout.TargetCurrency)
		fmt.Printf("  Rate: %.6f\n", directPayout.ExchangeRate)
		fmt.Printf("  Fee: %.2f\n", directPayout.Fee)
		fmt.Printf("  Status: %s\n", directPayout.Status)
	}

	payoutID := "sample-payout-id"
	if payout != nil {
		payoutID = payout.ID
	} else if directPayout != nil {
		payoutID = directPayout.ID
	}

	// 9. Get Payout Details
	fmt.Println("\n=== Get Payout Details ===")
	payoutDetail, err := client.Payout.GetPayout(ctx, payoutID)
	if err != nil {
		log.Printf("Failed to get payout: %v", err)
	} else {
		fmt.Printf("Payout Details:\n")
		fmt.Printf("  ID: %s\n", payoutDetail.ID)
		fmt.Printf("  Payee ID: %s\n", payoutDetail.PayeeID)
		fmt.Printf("  Beneficiary: %s\n", payoutDetail.BeneficiaryName)
		fmt.Printf("  Bank: %s\n", payoutDetail.BankName)
		fmt.Printf("  Account: %s\n", payoutDetail.AccountNumber)
		fmt.Printf("  Amount: %.2f %s -> %.2f %s\n",
			payoutDetail.SourceAmount, payoutDetail.SourceCurrency,
			payoutDetail.TargetAmount, payoutDetail.TargetCurrency)
		fmt.Printf("  Rate: %.6f\n", payoutDetail.ExchangeRate)
		fmt.Printf("  Status: %s\n", payoutDetail.Status)
		fmt.Printf("  Reference: %s\n", payoutDetail.Reference)
		if payoutDetail.CompletedAt != "" {
			fmt.Printf("  Completed At: %s\n", payoutDetail.CompletedAt)
		}
	}

	// 10. List Payouts
	fmt.Println("\n=== List Payouts ===")
	listPayoutsOptions := &interlace.ListPayoutsOptions{
		AccountID: accountID,
		Limit:     20,
		Page:      1,
	}
	payoutList, err := client.Payout.ListPayouts(ctx, listPayoutsOptions)
	if err != nil {
		log.Printf("Failed to list payouts: %v", err)
	} else {
		fmt.Printf("Total payouts: %d\n", payoutList.Total)
		
		// Group by status
		statusCount := make(map[string]int)
		var totalAmount float64
		
		for i, p := range payoutList.List {
			fmt.Printf("  %d. %s -> %s: %.2f %s (Status: %s)\n",
				i+1, p.SourceCurrency, p.TargetCurrency,
				p.SourceAmount, p.SourceCurrency, p.Status)
			
			statusCount[p.Status]++
			if p.Status == "COMPLETED" || p.Status == "PROCESSING" {
				totalAmount += p.SourceAmount
			}
		}
		
		fmt.Println("\nSummary:")
		for status, count := range statusCount {
			fmt.Printf("  %s: %d payouts\n", status, count)
		}
		fmt.Printf("  Total amount (completed/processing): %.2f\n", totalAmount)
	}

	// 11. Cancel Payout (only for pending payouts)
	fmt.Println("\n=== Cancel Payout (Example) ===")
	fmt.Println("Payout cancellation is demonstrated here (only works for pending payouts)")
	/*
	cancelResp, err := client.Payout.CancelPayout(ctx, payoutID)
	if err != nil {
		log.Printf("Failed to cancel payout: %v", err)
	} else {
		fmt.Printf("Payout cancelled:\n")
		fmt.Printf("  ID: %s\n", cancelResp.ID)
		fmt.Printf("  Status: %s\n", cancelResp.Status)
		fmt.Printf("  Message: %s\n", cancelResp.Message)
		fmt.Printf("  Success: %v\n", cancelResp.Success)
	}
	*/

	// 12. Filter Payouts by Status
	fmt.Println("\n=== List Completed Payouts ===")
	completedOptions := &interlace.ListPayoutsOptions{
		AccountID: accountID,
		Status:    "COMPLETED",
		Limit:     10,
		Page:      1,
	}
	completedList, err := client.Payout.ListPayouts(ctx, completedOptions)
	if err != nil {
		log.Printf("Failed to list completed payouts: %v", err)
	} else {
		fmt.Printf("Total completed payouts: %d\n", completedList.Total)
		for i, p := range completedList.List {
			fmt.Printf("  %d. %.2f %s -> %.2f %s (%s)\n",
				i+1, p.SourceAmount, p.SourceCurrency,
				p.TargetAmount, p.TargetCurrency,
				p.CompletedAt)
		}
	}

	fmt.Println("\n=== Payout Management Demo Completed ===")
	fmt.Println("\nPayout Flow:")
	fmt.Println("1. Get exchange rate")
	fmt.Println("2. Create payee bank account")
	fmt.Println("3. Create quotation (optional)")
	fmt.Println("4. Accept quotation OR create payout directly")
	fmt.Println("5. Monitor payout status")
	fmt.Println("6. Cancel if needed (only pending payouts)")
}
