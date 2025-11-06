package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// OAuthClient handles OAuth operations
type OAuthClient struct {
	httpClient *HTTPClient
}

// NewOAuthClient creates a new OAuth client
func NewOAuthClient(httpClient *HTTPClient) *OAuthClient {
	return &OAuthClient{
		httpClient: httpClient,
	}
}

// Authorize initiates the OAuth authorization flow
// Returns the authorization code that can be used to obtain access token
func (c *OAuthClient) Authorize(ctx context.Context, clientID string) (*OAuthAuthorizeData, error) {
	// Add query parameters
	params := url.Values{}
	params.Add("clientId", clientID)

	var authorizeResp OAuthAuthorizeResponse
	err := c.httpClient.DoGetRequestNoAuth(ctx, "/open-api/v3/oauth/authorize", params, &authorizeResp)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if authorizeResp.Code != "000000" {
		return nil, &Error{
			Code:    authorizeResp.Code,
			Message: authorizeResp.Message,
		}
	}

	return &authorizeResp.Data, nil
}

// GetAccessToken exchanges authorization code for access token
func (c *OAuthClient) GetAccessToken(ctx context.Context, code, clientID string) (*OAuthTokenData, error) {
	// Prepare request body
	tokenReq := OAuthTokenRequest{
		Code:     code,
		ClientID: clientID,
	}

	var tokenResp OAuthTokenResponse
	err := c.httpClient.DoPostRequestNoAuth(ctx, "/open-api/v3/oauth/access-token", tokenReq, &tokenResp)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if tokenResp.Code != "000000" {
		return nil, &Error{
			Code:    tokenResp.Code,
			Message: tokenResp.Message,
		}
	}

	return &tokenResp.Data, nil
}

// AuthorizeAndGetToken is a convenience method that combines authorize and token retrieval
func (c *OAuthClient) AuthorizeAndGetToken(ctx context.Context, clientID string) (*OAuthTokenData, error) {
	// Step 1: Get authorization code
	authData, err := c.Authorize(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("authorization failed: %w", err)
	}

	// Step 2: Exchange code for token
	tokenData, err := c.GetAccessToken(ctx, authData.Code, clientID)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}

	return tokenData, nil
}