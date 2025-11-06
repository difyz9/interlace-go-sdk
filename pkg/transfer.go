package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// TransferListOptions contains the query parameters for listing transfers
type TransferListOptions struct {
	WalletID  string `json:"walletId,omitempty"`
	Currency  string `json:"currency,omitempty"`
	Chain     string `json:"chain,omitempty"`
	Status    string `json:"status,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Page      int    `json:"page,omitempty"`
}

// CreateTransferRequest represents the request to create a blockchain transfer
type CreateTransferRequest struct {
	WalletID       string `json:"walletId"`
	Currency       string `json:"currency"`
	Chain          string `json:"chain"`
	Amount         string `json:"amount"`
	ToAddress      string `json:"toAddress"`
	Tag            string `json:"tag,omitempty"`
	IdempotencyKey string `json:"idempotencyKey"`
}

// FeeAndQuotaRequest represents the request to get fee and quota
type FeeAndQuotaRequest struct {
	WalletID  string `json:"walletId"`
	Currency  string `json:"currency"`
	Chain     string `json:"chain"`
	Amount    string `json:"amount"`
	ToAddress string `json:"toAddress"`
}

// TransferClient handles blockchain transfer operations
type TransferClient struct {
	httpClient *HTTPClient
}

// NewTransferClient creates a new transfer client
func NewTransferClient(httpClient *HTTPClient) *TransferClient {
	return &TransferClient{
		httpClient: httpClient,
	}
}

// CreateTransfer creates a new blockchain transfer
// POST /open-api/v3/cryptoconnect/transfers
func (c *TransferClient) CreateTransfer(ctx context.Context, req *CreateTransferRequest) (*BlockchainTransfer, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.WalletID == "" || req.Currency == "" || req.Chain == "" || req.Amount == "" || req.ToAddress == "" || req.IdempotencyKey == "" {
		return nil, fmt.Errorf("walletId, currency, chain, amount, toAddress and idempotencyKey are required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/cryptoconnect/transfers",
		Body:        req,
		RequireAuth: true,
	}

	var response BlockchainTransfer
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer: %w", err)
	}

	return &response, nil
}

// ListTransfers retrieves a list of blockchain transfers
// GET /open-api/v3/cryptoconnect/transfers
func (c *TransferClient) ListTransfers(ctx context.Context, options *TransferListOptions) (*TransferListResponse, error) {
	queryParams := url.Values{}
	
	if options != nil {
		if options.WalletID != "" {
			queryParams.Set("walletId", options.WalletID)
		}
		if options.Currency != "" {
			queryParams.Set("currency", options.Currency)
		}
		if options.Chain != "" {
			queryParams.Set("chain", options.Chain)
		}
		if options.Status != "" {
			queryParams.Set("status", options.Status)
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
		Endpoint:    "/open-api/v3/cryptoconnect/transfers",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var response TransferListResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list transfers: %w", err)
	}

	return &response, nil
}

// GetTransfer retrieves a specific blockchain transfer by ID
// GET /open-api/v3/cryptoconnect/transfers/{id}
func (c *TransferClient) GetTransfer(ctx context.Context, transferID string) (*BlockchainTransfer, error) {
	if transferID == "" {
		return nil, fmt.Errorf("transfer ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/cryptoconnect/transfers/%s", transferID),
		RequireAuth: true,
	}

	var response BlockchainTransfer
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfer %s: %w", transferID, err)
	}

	return &response, nil
}

// GetTransferKYT retrieves KYT (Know Your Transaction) information for a transfer
// GET /open-api/v3/cryptoconnect/transfers/{id}/kyt
func (c *TransferClient) GetTransferKYT(ctx context.Context, transferID string) (*TransferKYT, error) {
	if transferID == "" {
		return nil, fmt.Errorf("transfer ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/cryptoconnect/transfers/%s/kyt", transferID),
		RequireAuth: true,
	}

	var response TransferKYT
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfer KYT for %s: %w", transferID, err)
	}

	return &response, nil
}

// GetFeeAndQuota retrieves transfer fee and cross-chain quota information
// POST /open-api/v3/cryptoconnect/transfers/fee-and-quota
func (c *TransferClient) GetFeeAndQuota(ctx context.Context, req *FeeAndQuotaRequest) (*FeeAndQuota, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.WalletID == "" || req.Currency == "" || req.Chain == "" || req.Amount == "" || req.ToAddress == "" {
		return nil, fmt.Errorf("walletId, currency, chain, amount and toAddress are required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/cryptoconnect/transfers/fee-and-quota",
		Body:        req,
		RequireAuth: true,
	}

	var response FeeAndQuota
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get fee and quota: %w", err)
	}

	return &response, nil
}
