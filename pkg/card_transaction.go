package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// CardTransactionClient handles card transaction-related operations
type CardTransactionClient struct {
	httpClient *HTTPClient
}

// NewCardTransactionClient creates a new card transaction client
func NewCardTransactionClient(httpClient *HTTPClient) *CardTransactionClient {
	return &CardTransactionClient{
		httpClient: httpClient,
	}
}

// CardTransferInRequest contains the parameters for transferring funds into a prepaid card
type CardTransferInRequest struct {
	CardID              string  `json:"cardId"`
	Amount              float64 `json:"amount"`
	Currency            string  `json:"currency"`
	MerchantTradeNo     string  `json:"merchantTradeNo,omitempty"`
	InfinityAccountID   string  `json:"infinityAccountId,omitempty"`
	InfinityAccountType string  `json:"infinityAccountType,omitempty"`
}

// CardTransferInResponse represents the response from card transfer in
type CardTransferInResponse struct {
	ID                  string  `json:"id"`
	CardID              string  `json:"cardId"`
	Amount              float64 `json:"amount"`
	Currency            string  `json:"currency"`
	Status              string  `json:"status"`
	MerchantTradeNo     string  `json:"merchantTradeNo"`
	InfinityAccountID   string  `json:"infinityAccountId"`
	InfinityAccountType string  `json:"infinityAccountType"`
	CreatedAt           string  `json:"createdAt"`
	UpdatedAt           string  `json:"updatedAt"`
}

// CardTransferIn transfers funds from Quantum account to prepaid card
// POST /open-api/v3/cards/transfer-in
func (c *CardTransactionClient) CardTransferIn(ctx context.Context, req *CardTransferInRequest) (*CardTransferInResponse, error) {
	if req.CardID == "" {
		return nil, fmt.Errorf("cardId is required")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/cards/transfer-in",
		Body:        req,
		RequireAuth: true,
	}

	var response CardTransferInResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer in to card: %w", err)
	}

	return &response, nil
}

// CardTransferOutRequest contains the parameters for transferring funds out of a prepaid card
type CardTransferOutRequest struct {
	CardID              string  `json:"cardId"`
	Amount              float64 `json:"amount"`
	Currency            string  `json:"currency"`
	MerchantTradeNo     string  `json:"merchantTradeNo,omitempty"`
	InfinityAccountID   string  `json:"infinityAccountId,omitempty"`
	InfinityAccountType string  `json:"infinityAccountType,omitempty"`
}

// CardTransferOutResponse represents the response from card transfer out
type CardTransferOutResponse struct {
	ID                  string  `json:"id"`
	CardID              string  `json:"cardId"`
	Amount              float64 `json:"amount"`
	Currency            string  `json:"currency"`
	Status              string  `json:"status"`
	MerchantTradeNo     string  `json:"merchantTradeNo"`
	InfinityAccountID   string  `json:"infinityAccountId"`
	InfinityAccountType string  `json:"infinityAccountType"`
	CreatedAt           string  `json:"createdAt"`
	UpdatedAt           string  `json:"updatedAt"`
}

// CardTransferOut transfers funds from prepaid card to Quantum account
// POST /open-api/v3/cards/transfer-out
func (c *CardTransactionClient) CardTransferOut(ctx context.Context, req *CardTransferOutRequest) (*CardTransferOutResponse, error) {
	if req.CardID == "" {
		return nil, fmt.Errorf("cardId is required")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/cards/transfer-out",
		Body:        req,
		RequireAuth: true,
	}

	var response CardTransferOutResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer out from card: %w", err)
	}

	return &response, nil
}

// CardTransaction represents a card transaction
type CardTransaction struct {
	ID                     string  `json:"id"`
	CardID                 string  `json:"cardId"`
	Type                   string  `json:"type"`
	Amount                 float64 `json:"amount"`
	Currency               string  `json:"currency"`
	Status                 string  `json:"status"`
	MerchantName           string  `json:"merchantName"`
	MerchantCategoryCode   string  `json:"merchantCategoryCode"`
	MerchantCountry        string  `json:"merchantCountry"`
	AuthorizationCode      string  `json:"authorizationCode"`
	SettlementAmount       float64 `json:"settlementAmount"`
	SettlementCurrency     string  `json:"settlementCurrency"`
	BillingAmount          float64 `json:"billingAmount"`
	BillingCurrency        string  `json:"billingCurrency"`
	ExchangeRate           float64 `json:"exchangeRate"`
	CardholderID           string  `json:"cardholderId"`
	TransactionTime        string  `json:"transactionTime"`
	SettlementTime         string  `json:"settlementTime"`
	Description            string  `json:"description"`
	DeclineReason          string  `json:"declineReason"`
	IsInternational        bool    `json:"isInternational"`
	IsOnline               bool    `json:"isOnline"`
	MerchantTradeNo        string  `json:"merchantTradeNo"`
	InfinityAccountID      string  `json:"infinityAccountId"`
	InfinityAccountType    string  `json:"infinityAccountType"`
	CreatedAt              string  `json:"createdAt"`
	UpdatedAt              string  `json:"updatedAt"`
}

// ListCardTransactionsOptions contains the query parameters for listing card transactions
type ListCardTransactionsOptions struct {
	AccountID   string `json:"accountId,omitempty"`
	CardID      string `json:"cardId,omitempty"`
	Type        string `json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	StartTime   string `json:"startTime,omitempty"`
	EndTime     string `json:"endTime,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	Page        int    `json:"page,omitempty"`
}

// CardTransactionListResponse represents the response from list card transactions API
type CardTransactionListResponse struct {
	List  []CardTransaction `json:"list"`
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}

// ListCardTransactions retrieves a list of card transactions with optional filtering
// GET /open-api/v3/cards/transaction-list
func (c *CardTransactionClient) ListCardTransactions(ctx context.Context, options *ListCardTransactionsOptions) (*CardTransactionListResponse, error) {
	queryParams := url.Values{}
	
	if options != nil {
		if options.AccountID != "" {
			queryParams.Set("accountId", options.AccountID)
		}
		if options.CardID != "" {
			queryParams.Set("cardId", options.CardID)
		}
		if options.Type != "" {
			queryParams.Set("type", options.Type)
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
		if options.Limit > 0 {
			queryParams.Set("limit", fmt.Sprintf("%d", options.Limit))
		}
		if options.Page > 0 {
			queryParams.Set("page", fmt.Sprintf("%d", options.Page))
		}
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    "/open-api/v3/cards/transaction-list",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var response CardTransactionListResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list card transactions: %w", err)
	}

	return &response, nil
}
