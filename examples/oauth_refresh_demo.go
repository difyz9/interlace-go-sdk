package main

import (
	"context"
	"fmt"
	"log"
	"time"

	interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
	// Create a new client
	config := interlace.DefaultConfig()
	// For OAuth flow, we don't have an access token yet, so pass empty string
	httpClient := interlace.NewHTTPClient(config, "")
	oauthClient := interlace.NewOAuthClient(httpClient)

	ctx := context.Background()

	// Your client ID from Interlace
	
	clientID := "your-client-id-here"

	// Example 1: Get initial access token and refresh token
	fmt.Println("Step 1: Getting initial access token...")
	tokenData, err := oauthClient.AuthorizeAndGetToken(ctx, clientID)
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	fmt.Printf("Access Token: %s\n", tokenData.AccessToken)
	fmt.Printf("Refresh Token: %s\n", tokenData.RefreshToken)
	fmt.Printf("Expires In: %d seconds\n", tokenData.ExpiresIn)
	fmt.Printf("Timestamp: %d\n", tokenData.Timestamp)

	// Store the refresh token securely for later use
	refreshToken := tokenData.RefreshToken

	// Simulate waiting until token is about to expire or has expired
	fmt.Println("\nStep 2: Simulating token expiration...")
	time.Sleep(2 * time.Second) // In production, check if token is actually expired

	// Example 2: Refresh the access token using refresh token
	fmt.Println("\nStep 3: Refreshing access token...")
	newTokenData, err := oauthClient.RefreshToken(ctx, clientID, refreshToken)
	if err != nil {
		// This might fail in sandbox if the refresh token is not yet valid
		fmt.Printf("Note: Refresh token failed (this might be expected in sandbox): %v\n", err)
		fmt.Println("\nIn production, you would:")
		fmt.Println("1. Store the refresh token securely")
		fmt.Println("2. Use it when the access token expires (after ~24 hours)")
		fmt.Println("3. Get a new access token without re-authenticating")
	} else {
		fmt.Printf("New Access Token: %s\n", newTokenData.AccessToken)
		fmt.Printf("Expires In: %d seconds\n", newTokenData.ExpiresIn)
		fmt.Printf("Timestamp: %d\n", newTokenData.Timestamp)
		fmt.Println("\nToken refresh successful!")
	}

	fmt.Println("\nDemo completed!")
	fmt.Println("\nRefresh Token Usage:")
	fmt.Println("- Store the refresh token securely after initial authentication")
	fmt.Println("- Use it to get a new access token when the current one expires")
	fmt.Println("- This avoids requiring the user to re-authenticate")
	fmt.Println("- Typical use case: Long-lived background services or scheduled tasks")
}

// Helper function to check if token needs refresh
func shouldRefreshToken(expiresIn int, timestamp int64) bool {
	expiryTime := time.Unix(timestamp, 0).Add(time.Duration(expiresIn) * time.Second)
	// Refresh if less than 5 minutes until expiry
	return time.Until(expiryTime) < 5*time.Minute
}

// Example of a token manager that automatically refreshes tokens
type TokenManager struct {
	oauthClient  *interlace.OAuthClient
	clientID     string
	accessToken  string
	refreshToken string
	expiresIn    int
	timestamp    int64
}

func NewTokenManager(oauthClient *interlace.OAuthClient, clientID string) *TokenManager {
	return &TokenManager{
		oauthClient: oauthClient,
		clientID:    clientID,
	}
}

// GetAccessToken returns a valid access token, refreshing if necessary
func (tm *TokenManager) GetAccessToken(ctx context.Context) (string, error) {
	// Check if we need to refresh
	if shouldRefreshToken(tm.expiresIn, tm.timestamp) {
		newTokenData, err := tm.oauthClient.RefreshToken(ctx, tm.clientID, tm.refreshToken)
		if err != nil {
			return "", fmt.Errorf("failed to refresh token: %w", err)
		}

		tm.accessToken = newTokenData.AccessToken
		tm.expiresIn = newTokenData.ExpiresIn
		tm.timestamp = newTokenData.Timestamp
	}

	return tm.accessToken, nil
}

// Initialize performs initial OAuth flow and stores tokens
func (tm *TokenManager) Initialize(ctx context.Context) error {
	tokenData, err := tm.oauthClient.AuthorizeAndGetToken(ctx, tm.clientID)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	tm.accessToken = tokenData.AccessToken
	tm.refreshToken = tokenData.RefreshToken
	tm.expiresIn = tokenData.ExpiresIn
	tm.timestamp = tokenData.Timestamp

	return nil
}
