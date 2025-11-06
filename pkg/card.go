package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// CardListOptions contains the query parameters for listing cards
type CardListOptions struct {
	AccountID  string `json:"accountId,omitempty"`
	CardStatus string `json:"cardStatus,omitempty"`
	CardType   string `json:"cardType,omitempty"`
	IsActive   *bool  `json:"isActive,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Page       int    `json:"page,omitempty"`
}

// CardClient handles card-related operations
type CardClient struct {
	httpClient *HTTPClient
}

// NewCardClient creates a new card client
func NewCardClient(httpClient *HTTPClient) *CardClient {
	return &CardClient{
		httpClient: httpClient,
	}
}

// ListCards retrieves a list of cards with optional filtering
// GET /open-api/v3/card-list
func (c *CardClient) ListCards(ctx context.Context, options *CardListOptions) (*CardListResponse, error) {
	queryParams := url.Values{}
	
	if options != nil {
		if options.AccountID != "" {
			queryParams.Set("accountId", options.AccountID)
		}
		if options.CardStatus != "" {
			queryParams.Set("cardStatus", options.CardStatus)
		}
		if options.CardType != "" {
			queryParams.Set("cardType", options.CardType)
		}
		if options.IsActive != nil {
			if *options.IsActive {
				queryParams.Set("isActive", "true")
			} else {
				queryParams.Set("isActive", "false")
			}
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
		Endpoint:    "/open-api/v3/card-list",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var response CardListResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list cards: %w", err)
	}

	return &response, nil
}

// GetCardPrivateInfo retrieves sensitive card information (encrypted)
// GET /open-api/v3/cards/{id}
func (c *CardClient) GetCardPrivateInfo(ctx context.Context, cardID string) (*CardPrivateInfo, error) {
	if cardID == "" {
		return nil, fmt.Errorf("card ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/cards/%s", cardID),
		RequireAuth: true,
	}

	var response CardPrivateInfo
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get card private info for card %s: %w", cardID, err)
	}

	return &response, nil
}

// RemoveCard deletes a card
// DELETE /open-api/v3/cards/{id}
func (c *CardClient) RemoveCard(ctx context.Context, cardID string) (*CardRemoveResponse, error) {
	if cardID == "" {
		return nil, fmt.Errorf("card ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "DELETE",
		Endpoint:    fmt.Sprintf("/open-api/v3/cards/%s", cardID),
		RequireAuth: true,
	}

	var response CardRemoveResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to remove card %s: %w", cardID, err)
	}

	return &response, nil
}

// FreezeCard freezes a card to prevent transactions
// POST /open-api/v3/cards/{id}/freeze
func (c *CardClient) FreezeCard(ctx context.Context, cardID string) (*Card, error) {
	if cardID == "" {
		return nil, fmt.Errorf("card ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    fmt.Sprintf("/open-api/v3/cards/%s/freeze", cardID),
		RequireAuth: true,
	}

	var response Card
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to freeze card %s: %w", cardID, err)
	}

	return &response, nil
}

// UnfreezeCard unfreezes a card to allow transactions
// POST /open-api/v3/cards/{id}/unfreeze
func (c *CardClient) UnfreezeCard(ctx context.Context, cardID string) (*Card, error) {
	if cardID == "" {
		return nil, fmt.Errorf("card ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    fmt.Sprintf("/open-api/v3/cards/%s/unfreeze", cardID),
		RequireAuth: true,
	}

	var response Card
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unfreeze card %s: %w", cardID, err)
	}

	return &response, nil
}

// VelocityControlRequest contains the velocity control parameters
type VelocityControlRequest struct {
	DailySpendingLimit *float64 `json:"dailySpendingLimit,omitempty"`
	SingleTransLimit   *float64 `json:"singleTransLimit,omitempty"`
}

// SetCardVelocityControl sets transaction limits for a card
// PUT /open-api/v3/cards/{id}/velocity-control
func (c *CardClient) SetCardVelocityControl(ctx context.Context, cardID string, req *VelocityControlRequest) (*Card, error) {
	if cardID == "" {
		return nil, fmt.Errorf("card ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "PUT",
		Endpoint:    fmt.Sprintf("/open-api/v3/cards/%s/velocity-control", cardID),
		Body:        req,
		RequireAuth: true,
	}

	var response Card
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to set velocity control for card %s: %w", cardID, err)
	}

	return &response, nil
}

// CreatePrepaidCardRequest contains the parameters for creating a prepaid card
type CreatePrepaidCardRequest struct {
	BinID                    string  `json:"binId"`
	CardholderID             string  `json:"cardholderId"`
	ShippingAddressID        string  `json:"shippingAddressId,omitempty"`
	Label                    string  `json:"label,omitempty"`
	DailySpendingLimit       float64 `json:"dailySpendingLimit,omitempty"`
	SingleTransLimit         float64 `json:"singleTransLimit,omitempty"`
	MonthlySpendingLimit     float64 `json:"monthlySpendingLimit,omitempty"`
	ThreeDSecureAuthRequired bool    `json:"threeDSecureAuthRequired,omitempty"`
}

// CreatePrepaidCard creates a single prepaid card synchronously
// POST /open-api/v3/prepaid-card
func (c *CardClient) CreatePrepaidCard(ctx context.Context, req *CreatePrepaidCardRequest) (*Card, error) {
	if req.BinID == "" {
		return nil, fmt.Errorf("binId is required")
	}
	if req.CardholderID == "" {
		return nil, fmt.Errorf("cardholderId is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/prepaid-card",
		Body:        req,
		RequireAuth: true,
	}

	var response Card
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create prepaid card: %w", err)
	}

	return &response, nil
}

// BatchCreatePrepaidCardsRequest contains the parameters for batch creating prepaid cards
type BatchCreatePrepaidCardsRequest struct {
	Cards []CreatePrepaidCardRequest `json:"cards"`
}

// BatchCreatePrepaidCardsResponse represents the response from batch creating prepaid cards
type BatchCreatePrepaidCardsResponse struct {
	List  []Card `json:"list"`
	Total int    `json:"total"`
}

// BatchCreatePrepaidCards creates multiple prepaid cards at once (max 100)
// POST /open-api/v3/prepaid-cards
func (c *CardClient) BatchCreatePrepaidCards(ctx context.Context, cards []CreatePrepaidCardRequest) (*BatchCreatePrepaidCardsResponse, error) {
	if len(cards) == 0 {
		return nil, fmt.Errorf("cards list cannot be empty")
	}
	if len(cards) > 100 {
		return nil, fmt.Errorf("cannot create more than 100 cards at once")
	}

	req := &BatchCreatePrepaidCardsRequest{
		Cards: cards,
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/prepaid-cards",
		Body:        req,
		RequireAuth: true,
	}

	var response BatchCreatePrepaidCardsResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to batch create prepaid cards: %w", err)
	}

	return &response, nil
}

// CreateBudgetCardRequest contains the parameters for creating a budget card
type CreateBudgetCardRequest struct {
	BinID                    string  `json:"binId"`
	CardholderID             string  `json:"cardholderId"`
	BudgetID                 string  `json:"budgetId"`
	ShippingAddressID        string  `json:"shippingAddressId,omitempty"`
	Label                    string  `json:"label,omitempty"`
	DailySpendingLimit       float64 `json:"dailySpendingLimit,omitempty"`
	SingleTransLimit         float64 `json:"singleTransLimit,omitempty"`
	MonthlySpendingLimit     float64 `json:"monthlySpendingLimit,omitempty"`
	ThreeDSecureAuthRequired bool    `json:"threeDSecureAuthRequired,omitempty"`
}

// CreateBudgetCard creates a single budget card synchronously
// POST /open-api/v3/budget-card
func (c *CardClient) CreateBudgetCard(ctx context.Context, req *CreateBudgetCardRequest) (*Card, error) {
	if req.BinID == "" {
		return nil, fmt.Errorf("binId is required")
	}
	if req.CardholderID == "" {
		return nil, fmt.Errorf("cardholderId is required")
	}
	if req.BudgetID == "" {
		return nil, fmt.Errorf("budgetId is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/budget-card",
		Body:        req,
		RequireAuth: true,
	}

	var response Card
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create budget card: %w", err)
	}

	return &response, nil
}

// BatchCreateBudgetCardsRequest contains the parameters for batch creating budget cards
type BatchCreateBudgetCardsRequest struct {
	Cards []CreateBudgetCardRequest `json:"cards"`
}

// BatchCreateBudgetCardsResponse represents the response from batch creating budget cards
type BatchCreateBudgetCardsResponse struct {
	List  []Card `json:"list"`
	Total int    `json:"total"`
}

// BatchCreateBudgetCards creates multiple budget cards at once (max 100)
// POST /open-api/v3/budget-cards
func (c *CardClient) BatchCreateBudgetCards(ctx context.Context, cards []CreateBudgetCardRequest) (*BatchCreateBudgetCardsResponse, error) {
	if len(cards) == 0 {
		return nil, fmt.Errorf("cards list cannot be empty")
	}
	if len(cards) > 100 {
		return nil, fmt.Errorf("cannot create more than 100 cards at once")
	}

	req := &BatchCreateBudgetCardsRequest{
		Cards: cards,
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/budget-cards",
		Body:        req,
		RequireAuth: true,
	}

	var response BatchCreateBudgetCardsResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to batch create budget cards: %w", err)
	}

	return &response, nil
}

// CardSummary represents the summary information for a card
type CardSummary struct {
	CardID               string  `json:"cardId"`
	AvailableBalance     float64 `json:"availableBalance"`
	CurrentBalance       float64 `json:"currentBalance"`
	PendingTransactions  float64 `json:"pendingTransactions"`
	DailySpendingLimit   float64 `json:"dailySpendingLimit"`
	MonthlySpendingLimit float64 `json:"monthlySpendingLimit"`
	SingleTransLimit     float64 `json:"singleTransLimit"`
	SpentToday           float64 `json:"spentToday"`
	SpentThisMonth       float64 `json:"spentThisMonth"`
}

// GetCardSummary retrieves the summary information for a card
// GET /open-api/v3/cards/{id}/card-summary
func (c *CardClient) GetCardSummary(ctx context.Context, cardID string) (*CardSummary, error) {
	if cardID == "" {
		return nil, fmt.Errorf("card ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/cards/%s/card-summary", cardID),
		RequireAuth: true,
	}

	var response CardSummary
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get card summary for card %s: %w", cardID, err)
	}

	return &response, nil
}

// UpdateCardRequest contains the parameters for updating a card
type UpdateCardRequest struct {
	CardID                   string   `json:"cardId"`
	Label                    string   `json:"label,omitempty"`
	DailySpendingLimit       *float64 `json:"dailySpendingLimit,omitempty"`
	SingleTransLimit         *float64 `json:"singleTransLimit,omitempty"`
	MonthlySpendingLimit     *float64 `json:"monthlySpendingLimit,omitempty"`
	ThreeDSecureAuthRequired *bool    `json:"threeDSecureAuthRequired,omitempty"`
}

// UpdateCard updates card information
// PUT /open-api/v3/card
func (c *CardClient) UpdateCard(ctx context.Context, req *UpdateCardRequest) (*Card, error) {
	if req.CardID == "" {
		return nil, fmt.Errorf("cardId is required")
	}

	opts := &RequestOptions{
		Method:      "PUT",
		Endpoint:    "/open-api/v3/card",
		Body:        req,
		RequireAuth: true,
	}

	var response Card
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update card %s: %w", req.CardID, err)
	}

	return &response, nil
}

// BindWalletRequest contains the parameters for binding a wallet to a card
type BindWalletRequest struct {
	WalletID string `json:"walletId"`
}

// BindWallet binds a wallet to a card for spending (supports USDC and USDT)
// POST /open-api/v3/cards/{id}/bind-wallet
func (c *CardClient) BindWallet(ctx context.Context, cardID string, req *BindWalletRequest) (*Card, error) {
	if cardID == "" {
		return nil, fmt.Errorf("card ID cannot be empty")
	}
	if req.WalletID == "" {
		return nil, fmt.Errorf("walletId is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    fmt.Sprintf("/open-api/v3/cards/%s/bind-wallet", cardID),
		Body:        req,
		RequireAuth: true,
	}

	var response Card
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to bind wallet to card %s: %w", cardID, err)
	}

	return &response, nil
}