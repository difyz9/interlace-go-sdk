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

	// Example card ID
	cardID := "your-card-id"

	fmt.Println("=== Card Iframe Access Demo ===\n")

	// Get Card Access Token for Iframe
	fmt.Println("1. Getting Card Access Token for Iframe Display...")
	response, err := client.Iframe.GetCardAccessToken(ctx, cardID)
	if err != nil {
		log.Printf("Error getting card access token: %v\n", err)
	} else {
		fmt.Printf("Access token generated successfully!\n")
		fmt.Printf("  Card ID: %s\n", response.CardID)
		fmt.Printf("  Access Token: %s\n", response.AccessToken[:20]+"...") // Show partial token
		fmt.Printf("  Expires In: %d seconds\n", response.ExpiresIn)
		fmt.Printf("  Iframe URL: %s\n", response.IframeURL)
		fmt.Printf("  Created At: %s\n", response.CreatedAt)
	}

	// Usage Example
	fmt.Println("\n=== How to Use Card Iframe ===")
	fmt.Println("\n1. Backend: Get access token from API")
	fmt.Println("   - Call GetCardAccessToken(cardID)")
	fmt.Println("   - Receive temporary access token and iframe URL")
	fmt.Println("   - Token is valid for limited time (check ExpiresIn)")

	fmt.Println("\n2. Frontend: Embed iframe in your application")
	fmt.Println("   HTML Example:")
	fmt.Println(`   <iframe 
       src="{iframeURL}?token={accessToken}"
       width="400"
       height="250"
       style="border: none;"
       sandbox="allow-scripts allow-same-origin"
   ></iframe>`)

	fmt.Println("\n3. Security Considerations:")
	fmt.Println("   ✓ Access token expires quickly (check ExpiresIn)")
	fmt.Println("   ✓ Token is single-use or limited-use")
	fmt.Println("   ✓ Iframe displays card details securely")
	fmt.Println("   ✓ PCI DSS compliant - sensitive data stays on Interlace servers")
	fmt.Println("   ✓ No card details stored in your application")

	fmt.Println("\n4. Use Cases:")
	fmt.Println("   • Display card number, CVV, and expiry in web app")
	fmt.Println("   • Show card details in mobile app web view")
	fmt.Println("   • Allow users to view their card info securely")
	fmt.Println("   • Copy card details for online purchases")
	fmt.Println("   • Display virtual card information")

	fmt.Println("\n5. Best Practices:")
	fmt.Println("   • Generate new token each time user wants to view card")
	fmt.Println("   • Don't cache or store the access token")
	fmt.Println("   • Implement proper authentication before showing iframe")
	fmt.Println("   • Use HTTPS for all communications")
	fmt.Println("   • Verify user identity with 2FA before card display")
	fmt.Println("   • Log all card view requests for audit trail")
	fmt.Println("   • Set appropriate iframe sandbox attributes")
	fmt.Println("   • Monitor token expiration and refresh as needed")

	// Example: JavaScript Integration
	fmt.Println("\n=== JavaScript Integration Example ===")
	fmt.Println(`
async function showCardDetails(cardId) {
    // 1. Get access token from your backend
    const response = await fetch('/api/card/access-token', {
        method: 'POST',
        body: JSON.stringify({ cardId }),
        headers: { 'Content-Type': 'application/json' }
    });
    
    const { accessToken, iframeUrl, expiresIn } = await response.json();
    
    // 2. Create and display iframe
    const iframe = document.createElement('iframe');
    iframe.src = iframeUrl + '?token=' + accessToken;
    iframe.width = '400';
    iframe.height = '250';
    iframe.style.border = 'none';
    iframe.sandbox = 'allow-scripts allow-same-origin';
    
    document.getElementById('card-container').appendChild(iframe);
    
    // 3. Auto-remove iframe after token expires
    setTimeout(() => {
        iframe.remove();
        console.log('Card view session expired');
    }, expiresIn * 1000);
}
`)

	// Example: React Component
	fmt.Println("\n=== React Component Example ===")
	fmt.Println(`
import React, { useState, useEffect } from 'react';

function CardIframe({ cardId }) {
    const [tokenData, setTokenData] = useState(null);
    const [error, setError] = useState(null);
    
    useEffect(() => {
        async function fetchToken() {
            try {
                const response = await fetch('/api/card/access-token', {
                    method: 'POST',
                    body: JSON.stringify({ cardId }),
                    headers: { 'Content-Type': 'application/json' }
                });
                const data = await response.json();
                setTokenData(data);
                
                // Clear token after expiration
                setTimeout(() => {
                    setTokenData(null);
                }, data.expiresIn * 1000);
            } catch (err) {
                setError(err.message);
            }
        }
        
        fetchToken();
    }, [cardId]);
    
    if (error) return <div>Error: {error}</div>;
    if (!tokenData) return <div>Loading...</div>;
    
    return (
        <iframe
            src={tokenData.iframeUrl + '?token=' + tokenData.accessToken}
            width="400"
            height="250"
            style={{ border: 'none' }}
            sandbox="allow-scripts allow-same-origin"
            title="Card Details"
        />
    );
}
`)

	fmt.Println("\n=== Iframe Integration Complete ===")
	fmt.Println("\nThis secure approach ensures:")
	fmt.Println("• Card details never pass through your servers")
	fmt.Println("• PCI DSS compliance maintained")
	fmt.Println("• Temporary access prevents unauthorized viewing")
	fmt.Println("• Audit trail for all card access")
}
