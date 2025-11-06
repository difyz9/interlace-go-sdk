package main

import (
	"context"
	"fmt"
	"log"
	"time"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	// Create client
	config := interlace.DefaultConfig()
	client := interlace.NewClient(config)

	// Authenticate
	ctx := context.Background()
	clientID := "your-client-id-here"

	tokenData, err := client.Authenticate(ctx, clientID)
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	fmt.Printf("✅ Authenticated successfully! Token expires in: %d seconds\n\n", tokenData.ExpiresIn)

	// Run Business Account demos
	fmt.Println("=== Business Account Management Demo ===\n")

	// Demo 1: Create Legal Entity
	runCreateLegalEntityDemo(ctx, client)

	// Demo 2: Get Legal Entity
	runGetLegalEntityDemo(ctx, client)

	// Demo 3: Update Legal Entity
	runUpdateLegalEntityDemo(ctx, client)

	// Demo 4: Create Virtual Account
	runCreateVirtualAccountDemo(ctx, client)

	// Demo 5: Get Business Accounts
	runGetBusinessAccountsDemo(ctx, client)

	// Demo 6: Get Account Balance
	runGetAccountBalanceDemo(ctx, client)

	// Demo 7: Get Account Transactions
	runGetAccountTransactionsDemo(ctx, client)
}

// Demo 1: Create Legal Entity
func runCreateLegalEntityDemo(ctx context.Context, client *interlace.Client) {
	fmt.Println("📝 Demo 1: Creating Legal Entity...")

	request := &interlace.CreateLegalEntityRequest{
		EntityType:         "COMPANY",
		CompanyName:        "Tech Innovations Inc",
		RegistrationNumber: "REG-123456789",
		TaxID:              "TAX-987654321",
		IncorporationDate:  "2020-01-15",
		Country:            "US",
		Address: &interlace.Address{
			AddressLine1: "123 Tech Boulevard",
			AddressLine2: "Suite 400",
			City:         "San Francisco",
			State:        "CA",
			PostalCode:   "94105",
			Country:      "US",
		},
		ContactPerson: &interlace.ContactPerson{
			FirstName:   "John",
			LastName:    "Smith",
			Email:       "john.smith@techinnovations.com",
			PhoneNumber: "+1-415-555-0123",
			Position:    "CEO",
		},
		Directors: []interlace.Director{
			{
				FirstName:   "John",
				LastName:    "Smith",
				DateOfBirth: "1980-05-15",
				Nationality: "US",
				Position:    "CEO",
			},
			{
				FirstName:   "Jane",
				LastName:    "Doe",
				DateOfBirth: "1985-08-22",
				Nationality: "US",
				Position:    "CFO",
			},
		},
		UltimateBeneficiaryOwners: []interlace.UBO{
			{
				FirstName:           "John",
				LastName:            "Smith",
				DateOfBirth:         "1980-05-15",
				Nationality:         "US",
				OwnershipPercentage: 60.0,
			},
			{
				FirstName:           "Jane",
				LastName:            "Doe",
				DateOfBirth:         "1985-08-22",
				Nationality:         "US",
				OwnershipPercentage: 40.0,
			},
		},
	}

	entity, err := client.BusinessAccount.CreateLegalEntity(ctx, request)
	if err != nil {
		log.Printf("❌ Failed to create legal entity: %v\n", err)
		return
	}

	fmt.Printf("✅ Legal Entity Created:\n")
	fmt.Printf("   Entity ID: %s\n", entity.EntityID)
	fmt.Printf("   Company Name: %s\n", entity.CompanyName)
	fmt.Printf("   Status: %s\n", entity.Status)
	fmt.Printf("   Country: %s\n", entity.Country)
	fmt.Printf("   Directors: %d\n", len(entity.Directors))
	fmt.Printf("   UBOs: %d\n", len(entity.UltimateBeneficiaryOwners))
	fmt.Printf("   Created At: %s\n\n", entity.CreatedAt)
}

// Demo 2: Get Legal Entity
func runGetLegalEntityDemo(ctx context.Context, client *interlace.Client) {
	fmt.Println("🔍 Demo 2: Retrieving Legal Entity...")

	entityID := "your-entity-id-here"

	entity, err := client.BusinessAccount.GetLegalEntity(ctx, entityID)
	if err != nil {
		log.Printf("❌ Failed to get legal entity: %v\n", err)
		return
	}

	fmt.Printf("✅ Legal Entity Retrieved:\n")
	fmt.Printf("   Entity ID: %s\n", entity.EntityID)
	fmt.Printf("   Company Name: %s\n", entity.CompanyName)
	fmt.Printf("   Registration Number: %s\n", entity.RegistrationNumber)
	fmt.Printf("   Tax ID: %s\n", entity.TaxID)
	fmt.Printf("   Status: %s\n", entity.Status)
	fmt.Printf("   Address: %s, %s, %s %s\n",
		entity.Address.AddressLine1,
		entity.Address.City,
		entity.Address.State,
		entity.Address.PostalCode)

	if entity.ContactPerson != nil {
		fmt.Printf("   Contact: %s %s (%s)\n",
			entity.ContactPerson.FirstName,
			entity.ContactPerson.LastName,
			entity.ContactPerson.Email)
	}
	fmt.Println()
}

// Demo 3: Update Legal Entity
func runUpdateLegalEntityDemo(ctx context.Context, client *interlace.Client) {
	fmt.Println("✏️ Demo 3: Updating Legal Entity...")

	entityID := "your-entity-id-here"

	updateRequest := &interlace.UpdateLegalEntityRequest{
		CompanyName: "Tech Innovations International Inc",
		Address: &interlace.Address{
			AddressLine1: "456 Innovation Street",
			AddressLine2: "Floor 10",
			City:         "New York",
			State:        "NY",
			PostalCode:   "10001",
			Country:      "US",
		},
		ContactPerson: &interlace.ContactPerson{
			FirstName:   "John",
			LastName:    "Smith",
			Email:       "john.smith@techinnovations.com",
			PhoneNumber: "+1-212-555-0199",
			Position:    "CEO & President",
		},
		TaxID: "TAX-987654321-UPDATED",
	}

	entity, err := client.BusinessAccount.UpdateLegalEntity(ctx, entityID, updateRequest)
	if err != nil {
		log.Printf("❌ Failed to update legal entity: %v\n", err)
		return
	}

	fmt.Printf("✅ Legal Entity Updated:\n")
	fmt.Printf("   Company Name: %s\n", entity.CompanyName)
	fmt.Printf("   New Address: %s, %s\n", entity.Address.City, entity.Address.State)
	fmt.Printf("   Updated Tax ID: %s\n", entity.TaxID)
	fmt.Printf("   Updated At: %s\n\n", entity.UpdatedAt)
}

// Demo 4: Create Virtual Account
func runCreateVirtualAccountDemo(ctx context.Context, client *interlace.Client) {
	fmt.Println("🏦 Demo 4: Creating Virtual Account...")

	request := &interlace.CreateVirtualAccountRequest{
		LegalEntityID: "your-entity-id-here",
		Currency:      "USD",
		AccountType:   "CORPORATE",
		Reference:     "VA-2024-001",
	}

	account, err := client.BusinessAccount.CreateVirtualAccount(ctx, request)
	if err != nil {
		log.Printf("❌ Failed to create virtual account: %v\n", err)
		return
	}

	fmt.Printf("✅ Virtual Account Created:\n")
	fmt.Printf("   Account ID: %s\n", account.AccountID)
	fmt.Printf("   Account Number: %s\n", account.AccountNumber)
	fmt.Printf("   Currency: %s\n", account.Currency)
	fmt.Printf("   Status: %s\n", account.Status)
	fmt.Printf("   Initial Balance: %.2f\n", account.Balance)

	if account.IBAN != "" {
		fmt.Printf("   IBAN: %s\n", account.IBAN)
	}
	if account.SwiftCode != "" {
		fmt.Printf("   SWIFT: %s\n", account.SwiftCode)
	}
	if account.BankName != "" {
		fmt.Printf("   Bank: %s\n", account.BankName)
	}
	fmt.Printf("   Created At: %s\n\n", account.CreatedAt)
}

// Demo 5: Get Business Accounts
func runGetBusinessAccountsDemo(ctx context.Context, client *interlace.Client) {
	fmt.Println("📋 Demo 5: Listing Business Accounts...")

	legalEntityID := "your-entity-id-here"

	accounts, err := client.BusinessAccount.GetBusinessAccounts(ctx, legalEntityID)
	if err != nil {
		log.Printf("❌ Failed to get business accounts: %v\n", err)
		return
	}

	fmt.Printf("✅ Found %d Business Account(s):\n", len(accounts))

	for i, account := range accounts {
		fmt.Printf("\n   Account %d:\n", i+1)
		fmt.Printf("   ├─ ID: %s\n", account.AccountID)
		fmt.Printf("   ├─ Type: %s\n", account.AccountType)
		fmt.Printf("   ├─ Currency: %s\n", account.Currency)
		fmt.Printf("   ├─ Balance: %.2f\n", account.Balance)
		fmt.Printf("   ├─ Available: %.2f\n", account.AvailableBalance)
		fmt.Printf("   ├─ Frozen: %.2f\n", account.FrozenBalance)
		fmt.Printf("   ├─ Status: %s\n", account.Status)
		fmt.Printf("   └─ Account Number: %s\n", account.AccountNumber)
	}
	fmt.Println()
}

// Demo 6: Get Account Balance
func runGetAccountBalanceDemo(ctx context.Context, client *interlace.Client) {
	fmt.Println("💰 Demo 6: Checking Account Balance...")

	accountID := "your-account-id-here"

	balance, err := client.BusinessAccount.GetAccountBalance(ctx, accountID)
	if err != nil {
		log.Printf("❌ Failed to get account balance: %v\n", err)
		return
	}

	fmt.Printf("✅ Account Balance:\n")
	fmt.Printf("   Account ID: %s\n", balance.AccountID)
	fmt.Printf("   Currency: %s\n", balance.Currency)
	fmt.Printf("   ┌─ Total Balance: %.2f\n", balance.Balance)
	fmt.Printf("   ├─ Available Balance: %.2f\n", balance.AvailableBalance)
	fmt.Printf("   ├─ Frozen Balance: %.2f\n", balance.FrozenBalance)
	fmt.Printf("   └─ Pending Balance: %.2f\n", balance.PendingBalance)
	fmt.Printf("   Last Updated: %s\n\n", balance.LastUpdated)
}

// Demo 7: Get Account Transactions
func runGetAccountTransactionsDemo(ctx context.Context, client *interlace.Client) {
	fmt.Println("📊 Demo 7: Listing Account Transactions...")

	// Get transactions from the last 30 days
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -30)

	options := &interlace.ListBusinessAccountTransactionsOptions{
		AccountID: "your-account-id-here",
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   endTime.Format(time.RFC3339),
		Page:      1,
		Limit:     10,
	}

	response, err := client.BusinessAccount.GetAccountTransactions(ctx, options)
	if err != nil {
		log.Printf("❌ Failed to get account transactions: %v\n", err)
		return
	}

	fmt.Printf("✅ Found %d Transaction(s) (Page %d of %d):\n\n",
		response.TotalCount, response.Page, (response.TotalCount+response.Limit-1)/response.Limit)

	for i, tx := range response.Transactions {
		fmt.Printf("   Transaction %d:\n", i+1)
		fmt.Printf("   ├─ ID: %s\n", tx.TransactionID)
		fmt.Printf("   ├─ Type: %s\n", tx.Type)
		fmt.Printf("   ├─ Category: %s\n", tx.Category)
		fmt.Printf("   ├─ Amount: %.2f %s\n", tx.Amount, tx.Currency)
		fmt.Printf("   ├─ Balance Before: %.2f\n", tx.BalanceBefore)
		fmt.Printf("   ├─ Balance After: %.2f\n", tx.BalanceAfter)
		fmt.Printf("   ├─ Status: %s\n", tx.Status)

		if tx.Description != "" {
			fmt.Printf("   ├─ Description: %s\n", tx.Description)
		}
		if tx.CounterpartyName != "" {
			fmt.Printf("   ├─ Counterparty: %s\n", tx.CounterpartyName)
		}
		if tx.Reference != "" {
			fmt.Printf("   ├─ Reference: %s\n", tx.Reference)
		}

		fmt.Printf("   └─ Created At: %s\n\n", tx.CreatedAt)
	}
}

// Advanced Demo: Complete Business Account Workflow
func runCompleteWorkflowDemo(ctx context.Context, client *interlace.Client) {
	fmt.Println("\n=== Complete Business Account Workflow ===\n")

	// Step 1: Create legal entity
	fmt.Println("Step 1: Creating legal entity...")
	entityReq := &interlace.CreateLegalEntityRequest{
		EntityType:  "COMPANY",
		CompanyName: "Global Payments Ltd",
		Country:     "GB",
		Address: &interlace.Address{
			AddressLine1: "1 Financial Street",
			City:         "London",
			PostalCode:   "EC1A 1AA",
			Country:      "GB",
		},
		ContactPerson: &interlace.ContactPerson{
			FirstName:   "Alice",
			LastName:    "Johnson",
			Email:       "alice@globalpayments.com",
			PhoneNumber: "+44-20-7946-0958",
		},
	}

	entity, err := client.BusinessAccount.CreateLegalEntity(ctx, entityReq)
	if err != nil {
		log.Fatalf("Failed to create entity: %v", err)
	}
	fmt.Printf("✅ Entity created: %s\n\n", entity.EntityID)

	// Step 2: Create virtual accounts for multiple currencies
	fmt.Println("Step 2: Creating virtual accounts...")
	currencies := []string{"USD", "EUR", "GBP"}

	for _, currency := range currencies {
		accountReq := &interlace.CreateVirtualAccountRequest{
			LegalEntityID: entity.EntityID,
			Currency:      currency,
			Reference:     fmt.Sprintf("VA-%s-2024", currency),
		}

		account, err := client.BusinessAccount.CreateVirtualAccount(ctx, accountReq)
		if err != nil {
			log.Printf("Failed to create %s account: %v", currency, err)
			continue
		}

		fmt.Printf("   ✅ %s account: %s\n", currency, account.AccountNumber)
	}

	fmt.Println("\n✅ Workflow completed successfully!")
	fmt.Println("   - Legal entity established")
	fmt.Println("   - Multi-currency accounts created")
	fmt.Println("   - Ready for business operations")
}

// Best Practices Demo
func showBestPractices() {
	fmt.Println("\n=== Business Account Best Practices ===\n")

	fmt.Println("1️⃣ Legal Entity Management:")
	fmt.Println("   - Ensure all required documents are uploaded")
	fmt.Println("   - Keep UBO information up to date (ownership > 25%)")
	fmt.Println("   - Monitor entity status (PENDING → VERIFIED)")

	fmt.Println("\n2️⃣ Virtual Account Usage:")
	fmt.Println("   - Create separate accounts for different purposes")
	fmt.Println("   - Use meaningful references for tracking")
	fmt.Println("   - Monitor frozen/pending balances regularly")

	fmt.Println("\n3️⃣ Transaction Monitoring:")
	fmt.Println("   - Query transactions with date range filters")
	fmt.Println("   - Track balance changes (before/after)")
	fmt.Println("   - Use pagination for large result sets")

	fmt.Println("\n4️⃣ Security:")
	fmt.Println("   - Store entity IDs and account IDs securely")
	fmt.Println("   - Implement proper access controls")
	fmt.Println("   - Log all account operations for audit")

	fmt.Println("\n5️⃣ Error Handling:")
	fmt.Println("   - Always check error returns")
	fmt.Println("   - Implement retry logic for transient failures")
	fmt.Println("   - Handle KYC verification delays gracefully")
}
