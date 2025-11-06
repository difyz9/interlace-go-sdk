package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// InfinityAccountClient handles Infinity Account operations
type InfinityAccountClient struct {
	httpClient *HTTPClient
}

// NewInfinityAccountClient creates a new Infinity Account client
func NewInfinityAccountClient(httpClient *HTTPClient) *InfinityAccountClient {
	return &InfinityAccountClient{
		httpClient: httpClient,
	}
}

// InfinityAccountTransaction represents an Infinity Account transaction
type InfinityAccountTransaction struct {
	TransactionID   string  `json:"transactionId"`
	AccountID       string  `json:"accountId"`
	Type            string  `json:"type"` // DEBIT, CREDIT
	Category        string  `json:"category"` // CARD_LOAD, CARD_UNLOAD, TRANSFER_IN, TRANSFER_OUT, FEE, REFUND
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Balance         float64 `json:"balance"`
	BalanceBefore   float64 `json:"balanceBefore"`
	BalanceAfter    float64 `json:"balanceAfter"`
	RelatedID       string  `json:"relatedId,omitempty"` // Related card, transfer, or payment ID
	Status          string  `json:"status"`
	Description     string  `json:"description,omitempty"`
	MerchantTradeNo string  `json:"merchantTradeNo,omitempty"`
	CreatedAt       string  `json:"createdAt"`
}

// ListInfinityAccountTransactionsOptions represents options for listing transactions
type ListInfinityAccountTransactionsOptions struct {
	AccountID   string `json:"accountId,omitempty"`
	Type        string `json:"type,omitempty"`
	Category    string `json:"category,omitempty"`
	Currency    string `json:"currency,omitempty"`
	Status      string `json:"status,omitempty"`
	StartTime   string `json:"startTime,omitempty"`
	EndTime     string `json:"endTime,omitempty"`
	MinAmount   float64 `json:"minAmount,omitempty"`
	MaxAmount   float64 `json:"maxAmount,omitempty"`
	Page        int    `json:"page,omitempty"`
	Limit       int    `json:"limit,omitempty"`
}

// InfinityAccountTransactionListResponse represents the response for listing transactions
type InfinityAccountTransactionListResponse struct {
	Transactions []InfinityAccountTransaction `json:"transactions"`
	TotalCount   int                          `json:"totalCount"`
	Page         int                          `json:"page"`
	Limit        int                          `json:"limit"`
}

// ListInfinityAccountTransactions retrieves all Infinity Account transactions with optional filtering
// GET /open-api/v3/infinity-account/transactions
func (c *InfinityAccountClient) ListInfinityAccountTransactions(ctx context.Context, options *ListInfinityAccountTransactionsOptions) (*InfinityAccountTransactionListResponse, error) {
	var queryParams url.Values
	if options != nil {
		queryParams = url.Values{}
		if options.AccountID != "" {
			queryParams.Set("accountId", options.AccountID)
		}
		if options.Type != "" {
			queryParams.Set("type", options.Type)
		}
		if options.Category != "" {
			queryParams.Set("category", options.Category)
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
		if options.MinAmount > 0 {
			queryParams.Set("minAmount", fmt.Sprintf("%.2f", options.MinAmount))
		}
		if options.MaxAmount > 0 {
			queryParams.Set("maxAmount", fmt.Sprintf("%.2f", options.MaxAmount))
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
		Endpoint:    "/open-api/v3/infinity-account/transactions",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var resp InfinityAccountTransactionListResponse
	err := c.httpClient.DoRequest(ctx, opts, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to list Infinity Account transactions: %w", err)
	}

	return &resp, nil
}
