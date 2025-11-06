package interlace

import (
	"context"
	"fmt"
)

// SecurityClient handles card security operations
type SecurityClient struct {
	client *HTTPClient
}

// NewSecurityClient creates a new security client
func NewSecurityClient(client *HTTPClient) *SecurityClient {
	return &SecurityClient{client: client}
}

// UpdatePINRequest represents card PIN update request
type UpdatePINRequest struct {
	CardID     string `json:"cardId"`
	NewPIN     string `json:"newPin"`
	ConfirmPIN string `json:"confirmPin"`
}

// UpdatePINResponse represents card PIN update response
type UpdatePINResponse struct {
	CardID    string `json:"cardId"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updatedAt"`
	Message   string `json:"message"`
}

// UpdateCardPIN updates the PIN for a card
func (c *SecurityClient) UpdateCardPIN(ctx context.Context, req *UpdatePINRequest) (*UpdatePINResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("update PIN request is required")
	}
	if req.CardID == "" {
		return nil, fmt.Errorf("card ID is required")
	}
	if req.NewPIN == "" {
		return nil, fmt.Errorf("new PIN is required")
	}
	if req.ConfirmPIN == "" {
		return nil, fmt.Errorf("confirm PIN is required")
	}
	if req.NewPIN != req.ConfirmPIN {
		return nil, fmt.Errorf("new PIN and confirm PIN must match")
	}
	if len(req.NewPIN) < 4 || len(req.NewPIN) > 6 {
		return nil, fmt.Errorf("PIN must be 4-6 digits")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/card/update-pin",
		Body:        req,
		RequireAuth: true,
	}

	var resp UpdatePINResponse
	err := c.client.DoRequest(ctx, opts, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to update card PIN: %w", err)
	}
	return &resp, nil
}
