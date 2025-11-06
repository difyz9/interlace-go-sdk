package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// BusinessTransferClient handles business transfer operations
type BusinessTransferClient struct {
	httpClient *HTTPClient
}

// NewBusinessTransferClient creates a new business transfer client
func NewBusinessTransferClient(httpClient *HTTPClient) *BusinessTransferClient {
	return &BusinessTransferClient{
		httpClient: httpClient,
	}
}

// BusinessTransfer represents a business transfer transaction
type BusinessTransfer struct {
	TransferID       string  `json:"transferId"`
	FromAccountID    string  `json:"fromAccountId"`
	ToAccountID      string  `json:"toAccountId"`
	Currency         string  `json:"currency"`
	Amount           float64 `json:"amount"`
	TransferType     string  `json:"transferType"` // INTRA_ACCOUNT, DIFFERENT_ACCOUNT
	Fee              float64 `json:"fee"`
	Status           string  `json:"status"` // PENDING, PROCESSING, COMPLETED, FAILED
	MerchantTransferNo string `json:"merchantTransferNo,omitempty"`
	Description      string  `json:"description,omitempty"`
	CreatedAt        string  `json:"createdAt"`
	CompletedAt      string  `json:"completedAt,omitempty"`
	FailedReason     string  `json:"failedReason,omitempty"`
}

// IntraAccountTransferRequest represents intra-account transfer request
type IntraAccountTransferRequest struct {
	AccountID          string  `json:"accountId"`
	FromWalletID       string  `json:"fromWalletId"`
	ToWalletID         string  `json:"toWalletId"`
	Currency           string  `json:"currency"`
	Amount             float64 `json:"amount"`
	MerchantTransferNo string  `json:"merchantTransferNo,omitempty"`
	Description        string  `json:"description,omitempty"`
}

// DifferentAccountTransferRequest represents different-account transfer request
type DifferentAccountTransferRequest struct {
	FromAccountID      string  `json:"fromAccountId"`
	ToAccountID        string  `json:"toAccountId"`
	Currency           string  `json:"currency"`
	Amount             float64 `json:"amount"`
	MerchantTransferNo string  `json:"merchantTransferNo,omitempty"`
	Description        string  `json:"description,omitempty"`
}

// ListBusinessTransfersOptions represents options for listing transfers
type ListBusinessTransfersOptions struct {
	AccountID    string `json:"accountId,omitempty"`
	TransferType string `json:"transferType,omitempty"`
	Currency     string `json:"currency,omitempty"`
	Status       string `json:"status,omitempty"`
	StartTime    string `json:"startTime,omitempty"`
	EndTime      string `json:"endTime,omitempty"`
	Page         int    `json:"page,omitempty"`
	Limit        int    `json:"limit,omitempty"`
}

// BusinessTransferListResponse represents the response for listing transfers
type BusinessTransferListResponse struct {
	Transfers  []BusinessTransfer `json:"transfers"`
	TotalCount int                `json:"totalCount"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
}

// CreateIntraAccountTransfer creates an intra-account business transfer
// POST /open-api/v3/business/transfer/internal
func (c *BusinessTransferClient) CreateIntraAccountTransfer(ctx context.Context, req *IntraAccountTransferRequest) (*BusinessTransfer, error) {
	if req == nil {
		return nil, fmt.Errorf("transfer request is required")
	}
	if req.AccountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}
	if req.FromWalletID == "" {
		return nil, fmt.Errorf("from wallet ID is required")
	}
	if req.ToWalletID == "" {
		return nil, fmt.Errorf("to wallet ID is required")
	}
	if req.FromWalletID == req.ToWalletID {
		return nil, fmt.Errorf("from wallet and to wallet cannot be the same")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/business/transfer/internal",
		Body:        req,
		RequireAuth: true,
	}

	var transfer BusinessTransfer
	err := c.httpClient.DoRequest(ctx, opts, &transfer)
	if err != nil {
		return nil, fmt.Errorf("failed to create intra-account transfer: %w", err)
	}

	return &transfer, nil
}

// CreateDifferentAccountTransfer creates a different-account business transfer
// POST /open-api/v3/business/transfer/external
func (c *BusinessTransferClient) CreateDifferentAccountTransfer(ctx context.Context, req *DifferentAccountTransferRequest) (*BusinessTransfer, error) {
	if req == nil {
		return nil, fmt.Errorf("transfer request is required")
	}
	if req.FromAccountID == "" {
		return nil, fmt.Errorf("from account ID is required")
	}
	if req.ToAccountID == "" {
		return nil, fmt.Errorf("to account ID is required")
	}
	if req.FromAccountID == req.ToAccountID {
		return nil, fmt.Errorf("from account and to account cannot be the same")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/business/transfer/external",
		Body:        req,
		RequireAuth: true,
	}

	var transfer BusinessTransfer
	err := c.httpClient.DoRequest(ctx, opts, &transfer)
	if err != nil {
		return nil, fmt.Errorf("failed to create different-account transfer: %w", err)
	}

	return &transfer, nil
}

// ListBusinessTransfers retrieves all business transfers with optional filtering
// GET /open-api/v3/business/transfers
func (c *BusinessTransferClient) ListBusinessTransfers(ctx context.Context, options *ListBusinessTransfersOptions) (*BusinessTransferListResponse, error) {
	var queryParams url.Values
	if options != nil {
		queryParams = url.Values{}
		if options.AccountID != "" {
			queryParams.Set("accountId", options.AccountID)
		}
		if options.TransferType != "" {
			queryParams.Set("transferType", options.TransferType)
		}
		if options.Currency != "" {
			queryParams.Set("currency", options.Currency)
		}
		if options.Status != "" {
			queryParams.Set("status", options.Status)
		}
		if options.StartTime != "" {
			queryParams.Set("startTime", options.StartTime)
		}
		if options.EndTime != "" {
			queryParams.Set("endTime", options.EndTime)
		}
		if options.Page > 0 {
			queryParams.Set("page", fmt.Sprintf("%d", options.Page))
		}
		if options.Limit > 0 {
			queryParams.Set("limit", fmt.Sprintf("%d", options.Limit))
		}
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    "/open-api/v3/business/transfers",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var resp BusinessTransferListResponse
	err := c.httpClient.DoRequest(ctx, opts, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to list business transfers: %w", err)
	}

	return &resp, nil
}
