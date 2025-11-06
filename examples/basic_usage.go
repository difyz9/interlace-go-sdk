package main

import (
	"context"
	"fmt"
	"log"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	// Example 1: Quick Setup (Recommended for most use cases)
	fmt.Println("=== Quick Setup Example ===")
	quickSetupExample()

	fmt.Println("\n=== Manual Setup Example ===")
	manualSetupExample()

	fmt.Println("\n=== Account Management Example ===")
	accountManagementExample()

	fmt.Println("\n=== File Upload Example ===")
	fileUploadExample()
}

func quickSetupExample() {
	clientID := "your-client-id-here"
	
	// Quick setup - handles authentication automatically
	client, tokenData, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Printf("Quick setup failed: %v", err)
		return
	}

	fmt.Printf("Authentication successful!\n")
	fmt.Printf("Access Token: %s\n", tokenData.AccessToken)
	fmt.Printf("Expires In: %d seconds\n", tokenData.ExpiresIn)
	fmt.Printf("Client is authenticated: %v\n", client.IsAuthenticated())
}

func manualSetupExample() {
	clientID := "your-client-id-here"
	
	// Create client with custom config
	config := interlace.DefaultConfig()
	config.ClientID = clientID
	
	client := interlace.NewClient(config)
	
	// Manual authentication
	ctx := context.Background()
	tokenData, err := client.Authenticate(ctx, clientID)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		return
	}

	fmt.Printf("Manual authentication successful!\n")
	fmt.Printf("Access Token: %s\n", tokenData.AccessToken)
	fmt.Printf("Refresh Token: %s\n", tokenData.RefreshToken)
}

func accountOperationsExample() {
	fmt.Println("\n=== Account Operations Example ===")
	clientID := "your-client-id-here"
	
	client, _, err := interlace.QuickSetup(clientID, nil)

func fileUploadExample() {
	clientID := "your-client-id-here"
	accountID := "your-account-id-here" // Use actual account ID
	
	// Setup client
	client, _, err := interlace.QuickSetup(clientID, nil)
	if err != nil {
		log.Printf("Setup failed: %v", err)
		return
	}

	ctx := context.Background()

	// Example 1: Upload single file by path
	// Note: Make sure the file exists or create a test file
	filePath := "/path/to/your/file.jpg"
	
	fmt.Printf("Attempting to upload file: %s\n", filePath)
	
	uploadResp, err := client.File.UploadFile(ctx, filePath, accountID)
	if err != nil {
		log.Printf("File upload failed: %v", err)
		// This is expected if the file doesn't exist
		fmt.Println("Note: Make sure the file path exists for this example to work")
		return
	}

	fmt.Printf("File upload successful!\n")
	fmt.Printf("Response: %+v\n", uploadResp)

	// Example 2: Upload multiple files
	filePaths := []string{
		"/path/to/file1.jpg",
		"/path/to/file2.pdf",
	}

	multiUploadResp, err := client.File.UploadMultipleFiles(ctx, filePaths, accountID)
	if err != nil {
		log.Printf("Multiple file upload failed: %v", err)
		return
	}

	fmt.Printf("Multiple files upload successful!\n")
	fmt.Printf("Response: %+v\n", multiUploadResp)
}