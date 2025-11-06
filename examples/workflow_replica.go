package main

import (
	"context"
	"fmt"
	"log"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	// This example demonstrates the exact same workflow as shown in your dev01.md
	
	clientID := "your-client-id-here"
	
	fmt.Println("=== Replicating dev01.md Workflow ===")
	
	// Step 1: OAuth Authorization (equivalent to the first curl command)
	config := interlace.DefaultConfig()
	client := interlace.NewClient(config)
	
	ctx := context.Background()
	
	// Get authorization code
	authData, err := client.OAuth.Authorize(ctx, clientID)
	if err != nil {
		log.Fatalf("Authorization failed: %v", err)
	}
	
	fmt.Printf("1. Authorization successful!\n")
	fmt.Printf("   Code: %s\n", authData.Code)
	fmt.Printf("   Timestamp: %d\n", authData.Timestamp)
	
	// Step 2: Get Access Token (equivalent to the second curl command)
	tokenData, err := client.OAuth.GetAccessToken(ctx, authData.Code, clientID)
	if err != nil {
		log.Fatalf("Token exchange failed: %v", err)
	}
	
	fmt.Printf("\n2. Token exchange successful!\n")
	fmt.Printf("   Access Token: %s\n", tokenData.AccessToken)
	fmt.Printf("   Refresh Token: %s\n", tokenData.RefreshToken)
	fmt.Printf("   Expires In: %d seconds\n", tokenData.ExpiresIn)
	
	// Set the access token for authenticated requests
	client.SetAccessToken(tokenData.AccessToken)
	
	// Step 3: Register Account (equivalent to the third curl command)
	registerReq := &interlace.AccountRegisterRequest{
		PhoneNumber:      "15900000042",
		Email:            "15900000042@qq.com",
		Name:             "saf",
		PhoneCountryCode: "86",
	}
	
	account, err := client.Account.Register(ctx, registerReq)
	if err != nil {
		log.Fatalf("Account registration failed: %v", err)
	}
	
	fmt.Printf("\n3. Account registration successful!\n")
	fmt.Printf("   Account ID: %s\n", account.ID)
	fmt.Printf("   Display ID: %s\n", account.DisplayID)
	fmt.Printf("   Name: %s\n", account.VerifiedName)
	fmt.Printf("   Status: %s\n", account.Status)
	fmt.Printf("   Create Time: %s\n", account.CreateTime)
	
	// Step 4: List/Get Account (equivalent to the fourth curl command)
	listOpts := &interlace.AccountListOptions{
		AccountID: account.ID,
		Limit:     10,
		Page:      1,
	}
	
	accountList, err := client.Account.List(ctx, listOpts)
	if err != nil {
		log.Fatalf("Account listing failed: %v", err)
	}
	
	fmt.Printf("\n4. Account listing successful!\n")
	fmt.Printf("   Total accounts: %s\n", accountList.Total)
	for _, acc := range accountList.List {
		fmt.Printf("   - ID: %s\n", acc.ID)
		fmt.Printf("     Display ID: %s\n", acc.DisplayID)
		fmt.Printf("     Name: %s\n", acc.VerifiedName)
		fmt.Printf("     Status: %s\n", acc.Status)
		if acc.ParentAccountID != nil {
			fmt.Printf("     Parent Account ID: %s\n", *acc.ParentAccountID)
		}
	}
	
	// Step 5: File Upload (equivalent to the fifth curl command)
	// Note: This requires an actual file to exist
	fmt.Printf("\n5. File upload (simulated - requires actual file):\n")
	fmt.Printf("   Use client.File.UploadFile(ctx, filePath, accountID)\n")
	fmt.Printf("   Where filePath = path to your file (e.g., '20131022160002_JJwfv.jpeg')\n")
	fmt.Printf("   And accountID = '%s'\n", account.ID)
	
	// Example of how to upload if you have a file:
	/*
	uploadResp, err := client.File.UploadFile(ctx, "20131022160002_JJwfv.jpeg", account.ID)
	if err != nil {
		log.Printf("File upload failed: %v", err)
	} else {
		fmt.Printf("File upload successful: %+v\n", uploadResp)
	}
	*/
	
	fmt.Printf("\n=== Workflow completed successfully! ===\n")
}