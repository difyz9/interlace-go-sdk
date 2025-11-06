package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// WalletListOptions contains the query parameters for listing wallets
type WalletListOptions struct {
	AccountID string `json:"accountId,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Page      int    `json:"page,omitempty"`
}

// CreateWalletRequest represents the request to create a wallet
type CreateWalletRequest struct {
	AccountID      string `json:"accountId"`
	Nickname       string `json:"nickname,omitempty"`
	IdempotencyKey string `json:"idempotencyKey"`
}

// UpdateWalletRequest represents the request to update a wallet nickname
type UpdateWalletRequest struct {
	Nickname string `json:"nickname"`
}

// CreateAddressRequest represents the request to create a blockchain address
type CreateAddressRequest struct {
	Currency string `json:"currency"`
	Chain    string `json:"chain"`
}

// WalletClient handles crypto wallet operations
type WalletClient struct {
	httpClient *HTTPClient
}

// NewWalletClient creates a new wallet client
func NewWalletClient(httpClient *HTTPClient) *WalletClient {
	return &WalletClient{
		httpClient: httpClient,
	}
}

// CreateWallet creates a new crypto wallet
// POST /open-api/v3/cryptoconnect/wallets
func (c *WalletClient) CreateWallet(ctx context.Context, req *CreateWalletRequest) (*Wallet, error) {
	if req == nil || req.AccountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/cryptoconnect/wallets",
		Body:        req,
		RequireAuth: true,
	}

	var response Wallet
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return &response, nil
}

// ListWallets retrieves a list of wallets
// GET /open-api/v3/cryptoconnect/wallets
func (c *WalletClient) ListWallets(ctx context.Context, options *WalletListOptions) (*WalletListResponse, error) {
	queryParams := url.Values{}
	
	if options != nil {
		if options.AccountID != "" {
			queryParams.Set("accountId", options.AccountID)
		}
		if options.Limit > 0 {
			queryParams.Set("limit", fmt.Sprintf("%d", options.Limit))
		}
		if options.Page > 0 {
			queryParams.Set("page", fmt.Sprintf("%d", options.Page))
		}
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    "/open-api/v3/cryptoconnect/wallets",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var response WalletListResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list wallets: %w", err)
	}

	return &response, nil
}

// GetWallet retrieves a specific wallet by ID
// GET /open-api/v3/cryptoconnect/wallets/{id}
func (c *WalletClient) GetWallet(ctx context.Context, walletID string) (*Wallet, error) {
	if walletID == "" {
		return nil, fmt.Errorf("wallet ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/cryptoconnect/wallets/%s", walletID),
		RequireAuth: true,
	}

	var response Wallet
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet %s: %w", walletID, err)
	}

	return &response, nil
}

// UpdateWallet updates a wallet's nickname
// PATCH /open-api/v3/cryptoconnect/wallets/{id}
func (c *WalletClient) UpdateWallet(ctx context.Context, walletID string, req *UpdateWalletRequest) (*Wallet, error) {
	if walletID == "" {
		return nil, fmt.Errorf("wallet ID cannot be empty")
	}
	if req == nil || req.Nickname == "" {
		return nil, fmt.Errorf("nickname is required")
	}

	opts := &RequestOptions{
		Method:      "PATCH",
		Endpoint:    fmt.Sprintf("/open-api/v3/cryptoconnect/wallets/%s", walletID),
		Body:        req,
		RequireAuth: true,
	}

	var response Wallet
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update wallet %s: %w", walletID, err)
	}

	return &response, nil
}

// CreateWalletAddress creates a new blockchain address for a wallet
// POST /open-api/v3/cryptoconnect/wallets/{id}/addresses
func (c *WalletClient) CreateWalletAddress(ctx context.Context, walletID string, req *CreateAddressRequest) (*WalletAddress, error) {
	if walletID == "" {
		return nil, fmt.Errorf("wallet ID cannot be empty")
	}
	if req == nil || req.Currency == "" || req.Chain == "" {
		return nil, fmt.Errorf("currency and chain are required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    fmt.Sprintf("/open-api/v3/cryptoconnect/wallets/%s/addresses", walletID),
		Body:        req,
		RequireAuth: true,
	}

	var response WalletAddress
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create address for wallet %s: %w", walletID, err)
	}

	return &response, nil
}
