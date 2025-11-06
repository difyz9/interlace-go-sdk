package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// CommonClient handles common/system API calls
type CommonClient struct {
	httpClient *HTTPClient
}

// NewCommonClient creates a new common API client
func NewCommonClient(httpClient *HTTPClient) *CommonClient {
	return &CommonClient{
		httpClient: httpClient,
	}
}

// WalletBalance represents wallet balance information
type WalletBalance struct {
	WalletID        string  `json:"walletId"`
	Currency        string  `json:"currency"`
	AvailableBalance float64 `json:"availableBalance"`
	FrozenBalance   float64 `json:"frozenBalance"`
	TotalBalance    float64 `json:"totalBalance"`
}

// CardBinRecommendation represents recommended card BIN for high success rate
type CardBinRecommendation struct {
	BinID           string  `json:"binId"`
	CardBrand       string  `json:"cardBrand"`
	Currency        string  `json:"currency"`
	SuccessRate     float64 `json:"successRate"`
	RecommendationScore float64 `json:"recommendationScore"`
	Region          string  `json:"region"`
}

// SetConsumptionScenarioRequest represents request to set transaction scenarios
type SetConsumptionScenarioRequest struct {
	CardID     string   `json:"cardId"`
	ScenarioIDs []string `json:"scenarioIds"`
}

// SetConsumptionScenarioResponse represents response for setting transaction scenarios
type SetConsumptionScenarioResponse struct {
	CardID      string   `json:"cardId"`
	ScenarioIDs []string `json:"scenarioIds"`
	UpdatedAt   string   `json:"updatedAt"`
	Message     string   `json:"message"`
}

// ListConsumptionScenarios retrieves a list of all Infinity Card transaction scenarios
// GET /open-api/v3/card/sys/consumption-scenarios
func (c *CommonClient) ListConsumptionScenarios(ctx context.Context, accountID string) (*ConsumptionScenarioListResponse, error) {
	if accountID == "" {
		return nil, fmt.Errorf("accountId cannot be empty")
	}

	params := url.Values{}
	params.Set("accountId", accountID)

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    "/open-api/v3/card/sys/consumption-scenarios?" + params.Encode(),
		RequireAuth: true,
	}

	var response ConsumptionScenarioListResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list consumption scenarios: %w", err)
	}

	return &response, nil
}

// ListWallets retrieves wallet balances for an account
// GET /open-api/v3/wallets
func (c *CommonClient) ListWallets(ctx context.Context, accountID string) ([]WalletBalance, error) {
	if accountID == "" {
		return nil, fmt.Errorf("accountId cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/account/%s/wallets", accountID),
		RequireAuth: true,
	}

	var wallets []WalletBalance
	err := c.httpClient.DoRequest(ctx, opts, &wallets)
	if err != nil {
		return nil, fmt.Errorf("failed to list wallets: %w", err)
	}

	return wallets, nil
}

// GetCardBinRecommendation queries card BIN for high success rate
// GET /open-api/v3/card/bin/recommendation
func (c *CommonClient) GetCardBinRecommendation(ctx context.Context, currency, region string) ([]CardBinRecommendation, error) {
	if currency == "" {
		return nil, fmt.Errorf("currency cannot be empty")
	}

	params := url.Values{}
	params.Set("currency", currency)
	if region != "" {
		params.Set("region", region)
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    "/open-api/v3/card/bin/recommendation",
		QueryParams: params,
		RequireAuth: true,
	}

	var recommendations []CardBinRecommendation
	err := c.httpClient.DoRequest(ctx, opts, &recommendations)
	if err != nil {
		return nil, fmt.Errorf("failed to get card BIN recommendation: %w", err)
	}

	return recommendations, nil
}

// SetConsumptionScenario sets transaction scenarios for a card
// POST /open-api/v3/card/consumption-scenario
func (c *CommonClient) SetConsumptionScenario(ctx context.Context, req *SetConsumptionScenarioRequest) (*SetConsumptionScenarioResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	if req.CardID == "" {
		return nil, fmt.Errorf("card ID is required")
	}
	if len(req.ScenarioIDs) == 0 {
		return nil, fmt.Errorf("at least one scenario ID is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/card/consumption-scenario",
		Body:        req,
		RequireAuth: true,
	}

	var resp SetConsumptionScenarioResponse
	err := c.httpClient.DoRequest(ctx, opts, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to set consumption scenario: %w", err)
	}

	return &resp, nil
}
