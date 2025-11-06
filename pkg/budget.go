package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// BudgetClient handles budget-related operations
type BudgetClient struct {
	httpClient *HTTPClient
}

// NewBudgetClient creates a new budget client
func NewBudgetClient(httpClient *HTTPClient) *BudgetClient {
	return &BudgetClient{
		httpClient: httpClient,
	}
}

// Budget represents a budget resource
type Budget struct {
	ID                  string  `json:"id"`
	AccountID           string  `json:"accountId"`
	Name                string  `json:"name"`
	Currency            string  `json:"currency"`
	Balance             float64 `json:"balance"`
	AvailableBalance    float64 `json:"availableBalance"`
	PendingBalance      float64 `json:"pendingBalance"`
	Status              string  `json:"status"`
	Description         string  `json:"description"`
	CardCount           int     `json:"cardCount"`
	CreatedAt           string  `json:"createdAt"`
	UpdatedAt           string  `json:"updatedAt"`
}

// CreateBudgetRequest contains the parameters for creating a budget
type CreateBudgetRequest struct {
	AccountID   string  `json:"accountId"`
	Name        string  `json:"name"`
	Currency    string  `json:"currency"`
	Description string  `json:"description,omitempty"`
	InitBalance float64 `json:"initBalance,omitempty"`
}

// CreateBudget creates a new budget
// POST /open-api/v3/budgets
func (c *BudgetClient) CreateBudget(ctx context.Context, req *CreateBudgetRequest) (*Budget, error) {
	if req.AccountID == "" {
		return nil, fmt.Errorf("accountId is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/budgets",
		Body:        req,
		RequireAuth: true,
	}

	var response Budget
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	return &response, nil
}

// ListBudgetsOptions contains the query parameters for listing budgets
type ListBudgetsOptions struct {
	AccountID string `json:"accountId,omitempty"`
	Status    string `json:"status,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Page      int    `json:"page,omitempty"`
}

// BudgetListResponse represents the response from list budgets API
type BudgetListResponse struct {
	List  []Budget `json:"list"`
	Total int      `json:"total"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
}

// ListBudgets retrieves a list of budgets with optional filtering
// GET /open-api/v3/budgets
func (c *BudgetClient) ListBudgets(ctx context.Context, options *ListBudgetsOptions) (*BudgetListResponse, error) {
	queryParams := url.Values{}
	
	if options != nil {
		if options.AccountID != "" {
			queryParams.Set("accountId", options.AccountID)
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
		Endpoint:    "/open-api/v3/budgets",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var response BudgetListResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list budgets: %w", err)
	}

	return &response, nil
}

// GetBudget retrieves details of a specific budget
// GET /open-api/v3/budgets/{id}
func (c *BudgetClient) GetBudget(ctx context.Context, budgetID string) (*Budget, error) {
	if budgetID == "" {
		return nil, fmt.Errorf("budget ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/budgets/%s", budgetID),
		RequireAuth: true,
	}

	var response Budget
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget %s: %w", budgetID, err)
	}

	return &response, nil
}

// UpdateBudgetRequest contains the parameters for updating a budget
type UpdateBudgetRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

// UpdateBudget updates budget information
// PATCH /open-api/v3/budgets/{id}
func (c *BudgetClient) UpdateBudget(ctx context.Context, budgetID string, req *UpdateBudgetRequest) (*Budget, error) {
	if budgetID == "" {
		return nil, fmt.Errorf("budget ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "PATCH",
		Endpoint:    fmt.Sprintf("/open-api/v3/budgets/%s", budgetID),
		Body:        req,
		RequireAuth: true,
	}

	var response Budget
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget %s: %w", budgetID, err)
	}

	return &response, nil
}

// DeleteBudgetResponse represents the response from delete budget API
type DeleteBudgetResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// DeleteBudget deletes a budget
// DELETE /open-api/v3/budgets/{id}
func (c *BudgetClient) DeleteBudget(ctx context.Context, budgetID string) (*DeleteBudgetResponse, error) {
	if budgetID == "" {
		return nil, fmt.Errorf("budget ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "DELETE",
		Endpoint:    fmt.Sprintf("/open-api/v3/budgets/%s", budgetID),
		RequireAuth: true,
	}

	var response DeleteBudgetResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget %s: %w", budgetID, err)
	}

	return &response, nil
}

// IncreaseBudgetBalanceRequest contains the parameters for increasing budget balance
type IncreaseBudgetBalanceRequest struct {
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	MerchantTradeNo string  `json:"merchantTradeNo,omitempty"`
	Description     string  `json:"description,omitempty"`
}

// BudgetBalanceResponse represents the response from budget balance operations
type BudgetBalanceResponse struct {
	ID              string  `json:"id"`
	BudgetID        string  `json:"budgetId"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Type            string  `json:"type"`
	Status          string  `json:"status"`
	MerchantTradeNo string  `json:"merchantTradeNo"`
	Description     string  `json:"description"`
	BalanceBefore   float64 `json:"balanceBefore"`
	BalanceAfter    float64 `json:"balanceAfter"`
	CreatedAt       string  `json:"createdAt"`
}

// IncreaseBudgetBalance increases the budget balance (top-up)
// POST /open-api/v3/budgets/{id}/increase
func (c *BudgetClient) IncreaseBudgetBalance(ctx context.Context, budgetID string, req *IncreaseBudgetBalanceRequest) (*BudgetBalanceResponse, error) {
	if budgetID == "" {
		return nil, fmt.Errorf("budget ID cannot be empty")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    fmt.Sprintf("/open-api/v3/budgets/%s/increase", budgetID),
		Body:        req,
		RequireAuth: true,
	}

	var response BudgetBalanceResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to increase budget balance for %s: %w", budgetID, err)
	}

	return &response, nil
}

// DecreaseBudgetBalanceRequest contains the parameters for decreasing budget balance
type DecreaseBudgetBalanceRequest struct {
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	MerchantTradeNo string  `json:"merchantTradeNo,omitempty"`
	Description     string  `json:"description,omitempty"`
}

// DecreaseBudgetBalance decreases the budget balance (withdraw)
// POST /open-api/v3/budgets/{id}/decrease
func (c *BudgetClient) DecreaseBudgetBalance(ctx context.Context, budgetID string, req *DecreaseBudgetBalanceRequest) (*BudgetBalanceResponse, error) {
	if budgetID == "" {
		return nil, fmt.Errorf("budget ID cannot be empty")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    fmt.Sprintf("/open-api/v3/budgets/%s/decrease", budgetID),
		Body:        req,
		RequireAuth: true,
	}

	var response BudgetBalanceResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to decrease budget balance for %s: %w", budgetID, err)
	}

	return &response, nil
}

// BudgetTransaction represents a budget transaction
type BudgetTransaction struct {
	ID              string  `json:"id"`
	BudgetID        string  `json:"budgetId"`
	Type            string  `json:"type"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Status          string  `json:"status"`
	Description     string  `json:"description"`
	MerchantTradeNo string  `json:"merchantTradeNo"`
	BalanceBefore   float64 `json:"balanceBefore"`
	BalanceAfter    float64 `json:"balanceAfter"`
	CardID          string  `json:"cardId"`
	CardholderID    string  `json:"cardholderId"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// GetBudgetTransaction retrieves details of a specific budget transaction
// GET /open-api/v3/budgets/{id}/transactions/{transactionId}
func (c *BudgetClient) GetBudgetTransaction(ctx context.Context, budgetID, transactionID string) (*BudgetTransaction, error) {
	if budgetID == "" {
		return nil, fmt.Errorf("budget ID cannot be empty")
	}
	if transactionID == "" {
		return nil, fmt.Errorf("transaction ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/budgets/%s/transactions/%s", budgetID, transactionID),
		RequireAuth: true,
	}

	var response BudgetTransaction
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget transaction %s: %w", transactionID, err)
	}

	return &response, nil
}

// ListBudgetTransactionsOptions contains the query parameters for listing budget transactions
type ListBudgetTransactionsOptions struct {
	Type      string `json:"type,omitempty"`
	Status    string `json:"status,omitempty"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Page      int    `json:"page,omitempty"`
}

// BudgetTransactionListResponse represents the response from list budget transactions API
type BudgetTransactionListResponse struct {
	List  []BudgetTransaction `json:"list"`
	Total int                 `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
}

// ListBudgetTransactions retrieves a list of budget transactions
// GET /open-api/v3/budgets/{id}/transactions
func (c *BudgetClient) ListBudgetTransactions(ctx context.Context, budgetID string, options *ListBudgetTransactionsOptions) (*BudgetTransactionListResponse, error) {
	if budgetID == "" {
		return nil, fmt.Errorf("budget ID cannot be empty")
	}

	queryParams := url.Values{}
	
	if options != nil {
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
		Endpoint:    fmt.Sprintf("/open-api/v3/budgets/%s/transactions", budgetID),
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var response BudgetTransactionListResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list budget transactions for %s: %w", budgetID, err)
	}

	return &response, nil
}
