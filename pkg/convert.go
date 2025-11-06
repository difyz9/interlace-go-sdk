package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// ConvertClient handles cryptocurrency conversion operations
type ConvertClient struct {
	client *HTTPClient
}

// NewConvertClient creates a new convert client
func NewConvertClient(client *HTTPClient) *ConvertClient {
	return &ConvertClient{client: client}
}

// CurrencyPair represents a trading currency pair
type CurrencyPair struct {
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	MinAmount    float64 `json:"minAmount"`
	MaxAmount    float64 `json:"maxAmount"`
	Available    bool    `json:"available"`
}

// ConvertQuote represents a conversion quote
type ConvertQuote struct {
	QuoteID        string  `json:"quoteId"`
	FromCurrency   string  `json:"fromCurrency"`
	ToCurrency     string  `json:"toCurrency"`
	FromAmount     float64 `json:"fromAmount"`
	ToAmount       float64 `json:"toAmount"`
	ExchangeRate   float64 `json:"exchangeRate"`
	Fee            float64 `json:"fee"`
	ValidUntil     string  `json:"validUntil"`
	CreatedAt      string  `json:"createdAt"`
}

// ConvertTrade represents a conversion trade
type ConvertTrade struct {
	TradeID          string  `json:"tradeId"`
	WalletID         string  `json:"walletId"`
	FromCurrency     string  `json:"fromCurrency"`
	ToCurrency       string  `json:"toCurrency"`
	FromAmount       float64 `json:"fromAmount"`
	ToAmount         float64 `json:"toAmount"`
	ExchangeRate     float64 `json:"exchangeRate"`
	Fee              float64 `json:"fee"`
	Status           string  `json:"status"` // PENDING, COMPLETED, FAILED
	MerchantTradeNo  string  `json:"merchantTradeNo,omitempty"`
	CreatedAt        string  `json:"createdAt"`
	CompletedAt      string  `json:"completedAt,omitempty"`
}

// GetConvertQuoteRequest represents conversion quote request
type GetConvertQuoteRequest struct {
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	FromAmount   float64 `json:"fromAmount,omitempty"`
	ToAmount     float64 `json:"toAmount,omitempty"`
}

// CreateConvertTradeRequest represents conversion trade creation request
type CreateConvertTradeRequest struct {
	WalletID        string  `json:"walletId"`
	FromCurrency    string  `json:"fromCurrency"`
	ToCurrency      string  `json:"toCurrency"`
	FromAmount      float64 `json:"fromAmount,omitempty"`
	ToAmount        float64 `json:"toAmount,omitempty"`
	QuoteID         string  `json:"quoteId,omitempty"`
	MerchantTradeNo string  `json:"merchantTradeNo,omitempty"`
}

// ListConvertTradesOptions represents options for listing trades
type ListConvertTradesOptions struct {
	WalletID     string `json:"walletId,omitempty"`
	FromCurrency string `json:"fromCurrency,omitempty"`
	ToCurrency   string `json:"toCurrency,omitempty"`
	Status       string `json:"status,omitempty"`
	StartTime    string `json:"startTime,omitempty"`
	EndTime      string `json:"endTime,omitempty"`
	Page         int    `json:"page,omitempty"`
	Limit        int    `json:"limit,omitempty"`
}

// ConvertTradeListResponse represents the response for listing trades
type ConvertTradeListResponse struct {
	Trades     []ConvertTrade `json:"trades"`
	TotalCount int            `json:"totalCount"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
}

// GetCurrencyPairs retrieves all available trading currency pairs
func (c *ConvertClient) GetCurrencyPairs(ctx context.Context) ([]CurrencyPair, error) {
	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    "/open-api/v3/crypto/convert/currency-pairs",
		RequireAuth: true,
	}

	var pairs []CurrencyPair
	err := c.client.DoRequest(ctx, opts, &pairs)
	if err != nil {
		return nil, fmt.Errorf("failed to get currency pairs: %w", err)
	}
	return pairs, nil
}

// GetConvertQuote retrieves an estimate quote for conversion
func (c *ConvertClient) GetConvertQuote(ctx context.Context, req *GetConvertQuoteRequest) (*ConvertQuote, error) {
	if req == nil {
		return nil, fmt.Errorf("convert quote request is required")
	}
	if req.FromCurrency == "" {
		return nil, fmt.Errorf("from currency is required")
	}
	if req.ToCurrency == "" {
		return nil, fmt.Errorf("to currency is required")
	}
	if req.FromAmount <= 0 && req.ToAmount <= 0 {
		return nil, fmt.Errorf("either from amount or to amount must be specified")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/crypto/convert/quote",
		Body:        req,
		RequireAuth: true,
	}

	var quote ConvertQuote
	err := c.client.DoRequest(ctx, opts, &quote)
	if err != nil {
		return nil, fmt.Errorf("failed to get convert quote: %w", err)
	}
	return &quote, nil
}

// CreateConvertTrade creates a conversion trade
func (c *ConvertClient) CreateConvertTrade(ctx context.Context, req *CreateConvertTradeRequest) (*ConvertTrade, error) {
	if req == nil {
		return nil, fmt.Errorf("create trade request is required")
	}
	if req.WalletID == "" {
		return nil, fmt.Errorf("wallet ID is required")
	}
	if req.FromCurrency == "" {
		return nil, fmt.Errorf("from currency is required")
	}
	if req.ToCurrency == "" {
		return nil, fmt.Errorf("to currency is required")
	}
	if req.FromAmount <= 0 && req.ToAmount <= 0 {
		return nil, fmt.Errorf("either from amount or to amount must be specified")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/crypto/convert/trade",
		Body:        req,
		RequireAuth: true,
	}

	var trade ConvertTrade
	err := c.client.DoRequest(ctx, opts, &trade)
	if err != nil {
		return nil, fmt.Errorf("failed to create convert trade: %w", err)
	}
	return &trade, nil
}

// ListConvertTrades retrieves all conversion trades
func (c *ConvertClient) ListConvertTrades(ctx context.Context, options *ListConvertTradesOptions) (*ConvertTradeListResponse, error) {
	var queryParams url.Values
	if options != nil {
		queryParams = url.Values{}
		if options.WalletID != "" {
			queryParams.Set("walletId", options.WalletID)
		}
		if options.FromCurrency != "" {
			queryParams.Set("fromCurrency", options.FromCurrency)
		}
		if options.ToCurrency != "" {
			queryParams.Set("toCurrency", options.ToCurrency)
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
		Endpoint:    "/open-api/v3/crypto/convert/trades",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var resp ConvertTradeListResponse
	err := c.client.DoRequest(ctx, opts, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to list convert trades: %w", err)
	}
	return &resp, nil
}
