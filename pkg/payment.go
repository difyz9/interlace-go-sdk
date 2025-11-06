package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// CreatePaymentRequest represents the request to create a payment order
type CreatePaymentRequest struct {
	MerchantTradeNo string `json:"merchantTradeNo"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	Country         string `json:"country"`
	Description     string `json:"description,omitempty"`
	PaymentMethod   string `json:"paymentMethod,omitempty"`
	ReturnURL       string `json:"returnUrl,omitempty"`
	NotifyURL       string `json:"notifyUrl,omitempty"`
}

// CancelPaymentRequest represents the request to cancel a payment
type CancelPaymentRequest struct {
	OrderNo string `json:"orderNo"` // Can be system order ID or merchant order no
}

// CreateRefundRequest represents the request to create a refund
type CreateRefundRequest struct {
	SourceMerchantTradeNo string `json:"sourceMerchantTradeNo"` // The original payment merchant trade number
	MerchantTradeNo       string `json:"merchantTradeNo"`        // Refund merchant trade number
	Amount                string `json:"amount"`
	Reason                string `json:"reason,omitempty"`
}

// QueryPaymentRequest represents the request to query a payment
type QueryPaymentRequest struct {
	OrderNo string `json:"orderNo"` // Can be system order ID or merchant order no
}

// QueryRefundRequest represents the request to query a refund
type QueryRefundRequest struct {
	OrderNo string `json:"orderNo"` // Can be system refund ID or merchant order no
}

// SearchRequest represents the request to search payments and refunds
type SearchRequest struct {
	OrderNos []string `json:"orderNos"` // System order IDs or merchant order nos
}

// PaymentClient handles payment and refund operations
type PaymentClient struct {
	httpClient *HTTPClient
}

// NewPaymentClient creates a new payment client
func NewPaymentClient(httpClient *HTTPClient) *PaymentClient {
	return &PaymentClient{
		httpClient: httpClient,
	}
}

// CreatePayment creates a new payment order
// POST /open-api/v3/acquiring/payments
func (c *PaymentClient) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*Payment, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.MerchantTradeNo == "" || req.Amount == "" || req.Currency == "" || req.Country == "" {
		return nil, fmt.Errorf("merchantTradeNo, amount, currency and country are required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/acquiring/payments",
		Body:        req,
		RequireAuth: true,
	}

	var response Payment
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return &response, nil
}

// CancelPayment cancels an existing payment order
// POST /open-api/v3/acquiring/payments/cancel
func (c *PaymentClient) CancelPayment(ctx context.Context, req *CancelPaymentRequest) (*Payment, error) {
	if req == nil || req.OrderNo == "" {
		return nil, fmt.Errorf("orderNo is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/acquiring/payments/cancel",
		Body:        req,
		RequireAuth: true,
	}

	var response Payment
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel payment: %w", err)
	}

	return &response, nil
}

// CreateRefund creates a new refund for a payment
// POST /open-api/v3/acquiring/refunds
func (c *PaymentClient) CreateRefund(ctx context.Context, req *CreateRefundRequest) (*Refund, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.SourceMerchantTradeNo == "" || req.MerchantTradeNo == "" || req.Amount == "" {
		return nil, fmt.Errorf("sourceMerchantTradeNo, merchantTradeNo and amount are required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/acquiring/refunds",
		Body:        req,
		RequireAuth: true,
	}

	var response Refund
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	return &response, nil
}

// QueryPayment queries a single payment order
// GET /open-api/v3/acquiring/payments
func (c *PaymentClient) QueryPayment(ctx context.Context, orderNo string) (*Payment, error) {
	if orderNo == "" {
		return nil, fmt.Errorf("orderNo cannot be empty")
	}

	queryParams := url.Values{}
	queryParams.Set("orderNo", orderNo)

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    "/open-api/v3/acquiring/payments",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var response Payment
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query payment: %w", err)
	}

	return &response, nil
}

// QueryRefund queries a single refund order
// GET /open-api/v3/acquiring/refunds
func (c *PaymentClient) QueryRefund(ctx context.Context, orderNo string) (*Refund, error) {
	if orderNo == "" {
		return nil, fmt.Errorf("orderNo cannot be empty")
	}

	queryParams := url.Values{}
	queryParams.Set("orderNo", orderNo)

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    "/open-api/v3/acquiring/refunds",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var response Refund
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query refund: %w", err)
	}

	return &response, nil
}

// Search searches for multiple payments and refunds by order numbers
// POST /open-api/v3/acquiring/search
func (c *PaymentClient) Search(ctx context.Context, orderNos []string) (*SearchResult, error) {
	if len(orderNos) == 0 {
		return nil, fmt.Errorf("orderNos cannot be empty")
	}

	req := &SearchRequest{
		OrderNos: orderNos,
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/acquiring/search",
		Body:        req,
		RequireAuth: true,
	}

	var response SearchResult
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to search orders: %w", err)
	}

	return &response, nil
}
