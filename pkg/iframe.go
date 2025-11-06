package interlace

import (
	"context"
	"fmt"
)

// IframeClient handles iframe card access operations
type IframeClient struct {
	client *HTTPClient
}

// NewIframeClient creates a new iframe client
func NewIframeClient(client *HTTPClient) *IframeClient {
	return &IframeClient{client: client}
}

// CardAccessTokenRequest represents card access token request
type CardAccessTokenRequest struct {
	CardID string `json:"cardId"`
}

// CardAccessTokenResponse represents card access token response
type CardAccessTokenResponse struct {
	CardID      string `json:"cardId"`
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
	IframeURL   string `json:"iframeUrl"`
	CreatedAt   string `json:"createdAt"`
}

// GetCardAccessToken retrieves an access token for iframe card information display
func (c *IframeClient) GetCardAccessToken(ctx context.Context, cardID string) (*CardAccessTokenResponse, error) {
	if cardID == "" {
		return nil, fmt.Errorf("card ID is required")
	}

	req := &CardAccessTokenRequest{
		CardID: cardID,
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/card/access-token",
		Body:        req,
		RequireAuth: true,
	}

	var resp CardAccessTokenResponse
	err := c.client.DoRequest(ctx, opts, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get card access token: %w", err)
	}
	return &resp, nil
}
