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
	httpClient := interlace.NewHTTPClient(config)
	oauthClient := interlace.NewOAuthClient(httpClient)

	ctx := context.Background()

	// Your client ID from Interlace
	clientID := "qbit22c4571c943240d5"

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
		log.Fatalf("Failed to refresh token: %v", err)
	}

	fmt.Printf("New Access Token: %s\n", newTokenData.AccessToken)
	fmt.Printf("Expires In: %d seconds\n", newTokenData.ExpiresIn)
	fmt.Printf("Timestamp: %d\n", newTokenData.Timestamp)

	// Example 3: Token refresh with error handling
	fmt.Println("\nStep 4: Demonstrating error handling...")
	invalidRefreshToken := "invalid_token_123"
	_, err = oauthClient.RefreshToken(ctx, clientID, invalidRefreshToken)
	if err != nil {
		if interlaceErr, ok := err.(*interlace.Error); ok {
			fmt.Printf("Error Code: %s, Message: %s\n", interlaceErr.Code, interlaceErr.Message)
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}

	fmt.Println("\nDemo completed!")
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
