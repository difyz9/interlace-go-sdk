package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// PayoutClient handles payout-related operations
type PayoutClient struct {
	httpClient *HTTPClient
}

// NewPayoutClient creates a new payout client
func NewPayoutClient(httpClient *HTTPClient) *PayoutClient {
	return &PayoutClient{
		httpClient: httpClient,
	}
}

// ExchangeRateRequest contains the parameters for getting exchange rate
type ExchangeRateRequest struct {
	SourceCurrency string `json:"sourceCurrency"`
	TargetCurrency string `json:"targetCurrency"`
	Amount         float64 `json:"amount,omitempty"`
}

// ExchangeRateResponse represents the exchange rate response
type ExchangeRateResponse struct {
	SourceCurrency string  `json:"sourceCurrency"`
	TargetCurrency string  `json:"targetCurrency"`
	Rate           float64 `json:"rate"`
	InverseRate    float64 `json:"inverseRate"`
	Amount         float64 `json:"amount"`
	ConvertedAmount float64 `json:"convertedAmount"`
	Timestamp      string  `json:"timestamp"`
	ValidUntil     string  `json:"validUntil"`
}

// GetExchangeRate retrieves the current exchange rate between two currencies
// GET /open-api/v3/payment/rate
func (c *PayoutClient) GetExchangeRate(ctx context.Context, sourceCurrency, targetCurrency string, amount float64) (*ExchangeRateResponse, error) {
	if sourceCurrency == "" {
		return nil, fmt.Errorf("sourceCurrency is required")
	}
	if targetCurrency == "" {
		return nil, fmt.Errorf("targetCurrency is required")
	}

	queryParams := url.Values{}
	queryParams.Set("sourceCurrency", sourceCurrency)
	queryParams.Set("targetCurrency", targetCurrency)
	if amount > 0 {
		queryParams.Set("amount", fmt.Sprintf("%.2f", amount))
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    "/open-api/v3/payment/rate",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var response ExchangeRateResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	return &response, nil
}

// Payee represents a payee bank account
type Payee struct {
	ID                string `json:"id"`
	AccountID         string `json:"accountId"`
	BeneficiaryName   string `json:"beneficiaryName"`
	BankName          string `json:"bankName"`
	BankCode          string `json:"bankCode"`
	BankCountry       string `json:"bankCountry"`
	AccountNumber     string `json:"accountNumber"`
	IBAN              string `json:"iban"`
	SwiftCode         string `json:"swiftCode"`
	RoutingNumber     string `json:"routingNumber"`
	Currency          string `json:"currency"`
	BeneficiaryType   string `json:"beneficiaryType"`
	BeneficiaryAddress string `json:"beneficiaryAddress"`
	Status            string `json:"status"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
}

// CreatePayeeRequest contains the parameters for creating a payee
type CreatePayeeRequest struct {
	AccountID          string `json:"accountId"`
	BeneficiaryName    string `json:"beneficiaryName"`
	BankName           string `json:"bankName"`
	BankCode           string `json:"bankCode,omitempty"`
	BankCountry        string `json:"bankCountry"`
	AccountNumber      string `json:"accountNumber,omitempty"`
	IBAN               string `json:"iban,omitempty"`
	SwiftCode          string `json:"swiftCode,omitempty"`
	RoutingNumber      string `json:"routingNumber,omitempty"`
	Currency           string `json:"currency"`
	BeneficiaryType    string `json:"beneficiaryType,omitempty"`
	BeneficiaryAddress string `json:"beneficiaryAddress,omitempty"`
}

// CreatePayee creates a new payee bank account
// POST /open-api/v3/payee
func (c *PayoutClient) CreatePayee(ctx context.Context, req *CreatePayeeRequest) (*Payee, error) {
	if req.AccountID == "" {
		return nil, fmt.Errorf("accountId is required")
	}
	if req.BeneficiaryName == "" {
		return nil, fmt.Errorf("beneficiaryName is required")
	}
	if req.BankCountry == "" {
		return nil, fmt.Errorf("bankCountry is required")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/payee",
		Body:        req,
		RequireAuth: true,
	}

	var response Payee
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create payee: %w", err)
	}

	return &response, nil
}

// GetPayee retrieves details of a specific payee
// GET /open-api/v3/payee/{id}/detail
func (c *PayoutClient) GetPayee(ctx context.Context, payeeID string) (*Payee, error) {
	if payeeID == "" {
		return nil, fmt.Errorf("payee ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/payee/%s/detail", payeeID),
		RequireAuth: true,
	}

	var response Payee
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get payee %s: %w", payeeID, err)
	}

	return &response, nil
}

// ListPayeesOptions contains the query parameters for listing payees
type ListPayeesOptions struct {
	AccountID string `json:"accountId,omitempty"`
	Currency  string `json:"currency,omitempty"`
	Status    string `json:"status,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Page      int    `json:"page,omitempty"`
}

// PayeeListResponse represents the response from list payees API
type PayeeListResponse struct {
	List  []Payee `json:"list"`
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
}

// ListPayees retrieves a list of payees
// GET /open-api/v3/payees
func (c *PayoutClient) ListPayees(ctx context.Context, options *ListPayeesOptions) (*PayeeListResponse, error) {
	queryParams := url.Values{}
	
	if options != nil {
		if options.AccountID != "" {
			queryParams.Set("accountId", options.AccountID)
		}
		if options.Currency != "" {
			queryParams.Set("currency", options.Currency)
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
		Endpoint:    "/open-api/v3/payees",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var response PayeeListResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list payees: %w", err)
	}

	return &response, nil
}

// Payout represents a payout transaction
type Payout struct {
	ID                  string  `json:"id"`
	AccountID           string  `json:"accountId"`
	PayeeID             string  `json:"payeeId"`
	SourceCurrency      string  `json:"sourceCurrency"`
	SourceAmount        float64 `json:"sourceAmount"`
	TargetCurrency      string  `json:"targetCurrency"`
	TargetAmount        float64 `json:"targetAmount"`
	ExchangeRate        float64 `json:"exchangeRate"`
	Fee                 float64 `json:"fee"`
	Status              string  `json:"status"`
	PayoutMethod        string  `json:"payoutMethod"`
	Reference           string  `json:"reference"`
	MerchantTradeNo     string  `json:"merchantTradeNo"`
	BeneficiaryName     string  `json:"beneficiaryName"`
	BankName            string  `json:"bankName"`
	AccountNumber       string  `json:"accountNumber"`
	FailureReason       string  `json:"failureReason"`
	QuotationID         string  `json:"quotationId"`
	EstimatedArrival    string  `json:"estimatedArrival"`
	CompletedAt         string  `json:"completedAt"`
	CreatedAt           string  `json:"createdAt"`
	UpdatedAt           string  `json:"updatedAt"`
}

// CreatePayoutRequest contains the parameters for creating a payout
type CreatePayoutRequest struct {
	AccountID       string  `json:"accountId"`
	PayeeID         string  `json:"payeeId"`
	SourceCurrency  string  `json:"sourceCurrency"`
	SourceAmount    float64 `json:"sourceAmount"`
	TargetCurrency  string  `json:"targetCurrency"`
	TargetAmount    float64 `json:"targetAmount,omitempty"`
	MerchantTradeNo string  `json:"merchantTradeNo,omitempty"`
	Reference       string  `json:"reference,omitempty"`
	QuotationID     string  `json:"quotationId,omitempty"`
}

// CreatePayout creates a new payout transaction
// POST /open-api/v3/payment
func (c *PayoutClient) CreatePayout(ctx context.Context, req *CreatePayoutRequest) (*Payout, error) {
	if req.AccountID == "" {
		return nil, fmt.Errorf("accountId is required")
	}
	if req.PayeeID == "" {
		return nil, fmt.Errorf("payeeId is required")
	}
	if req.SourceCurrency == "" {
		return nil, fmt.Errorf("sourceCurrency is required")
	}
	if req.SourceAmount <= 0 {
		return nil, fmt.Errorf("sourceAmount must be greater than 0")
	}
	if req.TargetCurrency == "" {
		return nil, fmt.Errorf("targetCurrency is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/payment",
		Body:        req,
		RequireAuth: true,
	}

	var response Payout
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}

	return &response, nil
}

// GetPayout retrieves details of a specific payout
// GET /open-api/v3/payment/{id}/detail
func (c *PayoutClient) GetPayout(ctx context.Context, payoutID string) (*Payout, error) {
	if payoutID == "" {
		return nil, fmt.Errorf("payout ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/payment/%s/detail", payoutID),
		RequireAuth: true,
	}

	var response Payout
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get payout %s: %w", payoutID, err)
	}

	return &response, nil
}

// ListPayoutsOptions contains the query parameters for listing payouts
type ListPayoutsOptions struct {
	AccountID      string `json:"accountId,omitempty"`
	PayeeID        string `json:"payeeId,omitempty"`
	Status         string `json:"status,omitempty"`
	SourceCurrency string `json:"sourceCurrency,omitempty"`
	TargetCurrency string `json:"targetCurrency,omitempty"`
	StartTime      string `json:"startTime,omitempty"`
	EndTime        string `json:"endTime,omitempty"`
	Limit          int    `json:"limit,omitempty"`
	Page           int    `json:"page,omitempty"`
}

// PayoutListResponse represents the response from list payouts API
type PayoutListResponse struct {
	List  []Payout `json:"list"`
	Total int      `json:"total"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
}

// ListPayouts retrieves a list of payouts
// GET /open-api/v3/payments
func (c *PayoutClient) ListPayouts(ctx context.Context, options *ListPayoutsOptions) (*PayoutListResponse, error) {
	queryParams := url.Values{}
	
	if options != nil {
		if options.AccountID != "" {
			queryParams.Set("accountId", options.AccountID)
		}
		if options.PayeeID != "" {
			queryParams.Set("payeeId", options.PayeeID)
		}
		if options.Status != "" {
			queryParams.Set("status", options.Status)
		}
		if options.SourceCurrency != "" {
			queryParams.Set("sourceCurrency", options.SourceCurrency)
		}
		if options.TargetCurrency != "" {
			queryParams.Set("targetCurrency", options.TargetCurrency)
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
		Endpoint:    "/open-api/v3/payments",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var response PayoutListResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list payouts: %w", err)
	}

	return &response, nil
}

// Quotation represents a payout quotation
type Quotation struct {
	ID                string  `json:"id"`
	AccountID         string  `json:"accountId"`
	SourceCurrency    string  `json:"sourceCurrency"`
	SourceAmount      float64 `json:"sourceAmount"`
	TargetCurrency    string  `json:"targetCurrency"`
	TargetAmount      float64 `json:"targetAmount"`
	ExchangeRate      float64 `json:"rate"`
	Fee               float64 `json:"fee"`
	TotalAmount       float64 `json:"totalAmount"`
	Status            string  `json:"status"`
	ValidUntil        string  `json:"validUntil"`
	EstimatedArrival  string  `json:"estimatedArrival"`
	CreatedAt         string  `json:"createdAt"`
}

// CreateQuotationRequest contains the parameters for creating a quotation
type CreateQuotationRequest struct {
	AccountID      string  `json:"accountId"`
	SourceCurrency string  `json:"sourceCurrency"`
	SourceAmount   float64 `json:"sourceAmount"`
	TargetCurrency string  `json:"targetCurrency"`
	TargetAmount   float64 `json:"targetAmount,omitempty"`
	PayeeID        string  `json:"payeeId,omitempty"`
}

// CreateQuotation creates a payout quotation
// POST /open-api/v3/payment/quotation
func (c *PayoutClient) CreateQuotation(ctx context.Context, req *CreateQuotationRequest) (*Quotation, error) {
	if req.AccountID == "" {
		return nil, fmt.Errorf("accountId is required")
	}
	if req.SourceCurrency == "" {
		return nil, fmt.Errorf("sourceCurrency is required")
	}
	if req.SourceAmount <= 0 {
		return nil, fmt.Errorf("sourceAmount must be greater than 0")
	}
	if req.TargetCurrency == "" {
		return nil, fmt.Errorf("targetCurrency is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/payment/quotation",
		Body:        req,
		RequireAuth: true,
	}

	var response Quotation
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create quotation: %w", err)
	}

	return &response, nil
}

// GetQuotation retrieves details of a specific quotation
// GET /open-api/v3/payment/quotation/{id}
func (c *PayoutClient) GetQuotation(ctx context.Context, quotationID string) (*Quotation, error) {
	if quotationID == "" {
		return nil, fmt.Errorf("quotation ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/payment/quotation/%s", quotationID),
		RequireAuth: true,
	}

	var response Quotation
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotation %s: %w", quotationID, err)
	}

	return &response, nil
}

// AcceptQuotationRequest contains the parameters for accepting a quotation
type AcceptQuotationRequest struct {
	PayeeID         string `json:"payeeId"`
	MerchantTradeNo string `json:"merchantTradeNo,omitempty"`
	Reference       string `json:"reference,omitempty"`
}

// AcceptQuotation accepts a quotation and creates a payout
// POST /open-api/v3/payment/quotation/{id}/accept
func (c *PayoutClient) AcceptQuotation(ctx context.Context, quotationID string, req *AcceptQuotationRequest) (*Payout, error) {
	if quotationID == "" {
		return nil, fmt.Errorf("quotation ID cannot be empty")
	}
	if req.PayeeID == "" {
		return nil, fmt.Errorf("payeeId is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    fmt.Sprintf("/open-api/v3/payment/quotation/%s/accept", quotationID),
		Body:        req,
		RequireAuth: true,
	}

	var response Payout
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to accept quotation %s: %w", quotationID, err)
	}

	return &response, nil
}

// CancelPayoutResponse represents the response from canceling a payout
type CancelPayoutResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// CancelPayout cancels a pending payout
// POST /open-api/v3/payment/{id}/cancel
func (c *PayoutClient) CancelPayout(ctx context.Context, payoutID string) (*CancelPayoutResponse, error) {
	if payoutID == "" {
		return nil, fmt.Errorf("payout ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    fmt.Sprintf("/open-api/v3/payment/%s/cancel", payoutID),
		RequireAuth: true,
	}

	var response CancelPayoutResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel payout %s: %w", payoutID, err)
	}

	return &response, nil
}
