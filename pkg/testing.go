package interlace

import (
	"context"
	"fmt"
)

// TestingClient handles testing and simulation operations
type TestingClient struct {
	httpClient *HTTPClient
}

// NewTestingClient creates a new testing client
func NewTestingClient(httpClient *HTTPClient) *TestingClient {
	return &TestingClient{
		httpClient: httpClient,
	}
}

// SimulateAuthorizationRequest represents card authorization simulation request
type SimulateAuthorizationRequest struct {
	CardID          string  `json:"cardId"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	MerchantName    string  `json:"merchantName"`
	MerchantCountry string  `json:"merchantCountry,omitempty"`
	MerchantCity    string  `json:"merchantCity,omitempty"`
	MCC             string  `json:"mcc,omitempty"` // Merchant Category Code
	AuthType        string  `json:"authType,omitempty"` // PURCHASE, ATM_WITHDRAWAL, REFUND
	IsOnline        bool    `json:"isOnline"`
	IsFallback      bool    `json:"isFallback,omitempty"` // Chip fallback to magnetic stripe
}

// SimulateAuthorizationResponse represents card authorization simulation response
type SimulateAuthorizationResponse struct {
	SimulationID     string  `json:"simulationId"`
	CardID           string  `json:"cardId"`
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
	AuthorizationCode string `json:"authorizationCode"`
	Status           string  `json:"status"` // APPROVED, DECLINED
	DeclineReason    string  `json:"declineReason,omitempty"`
	AvailableBalance float64 `json:"availableBalance"`
	MerchantName     string  `json:"merchantName"`
	TransactionID    string  `json:"transactionId,omitempty"`
	CreatedAt        string  `json:"createdAt"`
}

// SimulateCardAuthorization simulates a card authorization request for testing
// POST /open-api/v3/testing/simulate-authorization
func (c *TestingClient) SimulateCardAuthorization(ctx context.Context, req *SimulateAuthorizationRequest) (*SimulateAuthorizationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("simulation request is required")
	}
	if req.CardID == "" {
		return nil, fmt.Errorf("card ID is required")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}
	if req.MerchantName == "" {
		return nil, fmt.Errorf("merchant name is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/testing/simulate-authorization",
		Body:        req,
		RequireAuth: true,
	}

	var response SimulateAuthorizationResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to simulate card authorization: %w", err)
	}

	return &response, nil
}
