package interlace

import (
	"context"
	"fmt"
)

// SweepingClient handles cryptocurrency sweeping operations
type SweepingClient struct {
	httpClient *HTTPClient
}

// NewSweepingClient creates a new sweeping client
func NewSweepingClient(httpClient *HTTPClient) *SweepingClient {
	return &SweepingClient{
		httpClient: httpClient,
	}
}

// SweepingRequest represents sweeping request
type SweepingRequest struct {
	WalletID      string   `json:"walletId"`
	FromAddresses []string `json:"fromAddresses"`
	ToAddress     string   `json:"toAddress"`
	Chain         string   `json:"chain"`
	Currency      string   `json:"currency"`
	MinAmount     float64  `json:"minAmount,omitempty"` // Minimum amount to sweep
}

// SweepingResponse represents sweeping response
type SweepingResponse struct {
	SweepingID    string              `json:"sweepingId"`
	WalletID      string              `json:"walletId"`
	Chain         string              `json:"chain"`
	Currency      string              `json:"currency"`
	TotalAmount   float64             `json:"totalAmount"`
	TotalGasFee   float64             `json:"totalGasFee"`
	Transactions  []SweepTransaction  `json:"transactions"`
	Status        string              `json:"status"` // PENDING, PROCESSING, COMPLETED, FAILED
	CreatedAt     string              `json:"createdAt"`
	CompletedAt   string              `json:"completedAt,omitempty"`
}

// SweepTransaction represents individual sweep transaction
type SweepTransaction struct {
	FromAddress string  `json:"fromAddress"`
	ToAddress   string  `json:"toAddress"`
	Amount      float64 `json:"amount"`
	GasFee      float64 `json:"gasFee"`
	TxHash      string  `json:"txHash"`
	Status      string  `json:"status"`
}

// Sweeping performs cryptocurrency sweeping from multiple addresses to a single address
// POST /open-api/v3/crypto/sweeping
func (c *SweepingClient) Sweeping(ctx context.Context, req *SweepingRequest) (*SweepingResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("sweeping request is required")
	}
	if req.WalletID == "" {
		return nil, fmt.Errorf("wallet ID is required")
	}
	if len(req.FromAddresses) == 0 {
		return nil, fmt.Errorf("at least one from address is required")
	}
	if req.ToAddress == "" {
		return nil, fmt.Errorf("to address is required")
	}
	if req.Chain == "" {
		return nil, fmt.Errorf("chain is required")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	// Validate that to address is not in from addresses
	for _, fromAddr := range req.FromAddresses {
		if fromAddr == req.ToAddress {
			return nil, fmt.Errorf("to address cannot be in from addresses list")
		}
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/crypto/sweeping",
		Body:        req,
		RequireAuth: true,
	}

	var response SweepingResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to perform sweeping: %w", err)
	}

	return &response, nil
}
