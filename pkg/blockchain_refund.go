package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// BlockchainRefundClient handles blockchain refund operations
type BlockchainRefundClient struct {
	httpClient *HTTPClient
}

// NewBlockchainRefundClient creates a new blockchain refund client
func NewBlockchainRefundClient(httpClient *HTTPClient) *BlockchainRefundClient {
	return &BlockchainRefundClient{
		httpClient: httpClient,
	}
}

// BlockchainRefund represents a blockchain refund transaction
type BlockchainRefund struct {
	RefundID         string  `json:"refundId"`
	WalletID         string  `json:"walletId"`
	TransferID       string  `json:"transferId"`
	Chain            string  `json:"chain"`
	Currency         string  `json:"currency"`
	Amount           float64 `json:"amount"`
	FromAddress      string  `json:"fromAddress"`
	ToAddress        string  `json:"toAddress"`
	TxHash           string  `json:"txHash"`
	GasFee           float64 `json:"gasFee"`
	Status           string  `json:"status"` // PENDING, PROCESSING, COMPLETED, FAILED
	MerchantRefundNo string  `json:"merchantRefundNo,omitempty"`
	Reason           string  `json:"reason,omitempty"`
	CreatedAt        string  `json:"createdAt"`
	CompletedAt      string  `json:"completedAt,omitempty"`
	FailedReason     string  `json:"failedReason,omitempty"`
}

// CreateBlockchainRefundRequest represents blockchain refund creation request
type CreateBlockchainRefundRequest struct {
	WalletID         string  `json:"walletId"`
	TransferID       string  `json:"transferId"`
	Chain            string  `json:"chain"`
	Currency         string  `json:"currency"`
	Amount           float64 `json:"amount"`
	ToAddress        string  `json:"toAddress"`
	MerchantRefundNo string  `json:"merchantRefundNo,omitempty"`
	Reason           string  `json:"reason,omitempty"`
}

// RefundGasFee represents gas fee information for refund
type RefundGasFee struct {
	Chain          string  `json:"chain"`
	Currency       string  `json:"currency"`
	EstimatedGasFee float64 `json:"estimatedGasFee"`
	GasPrice       string  `json:"gasPrice"`
	GasLimit       int64   `json:"gasLimit"`
	ValidUntil     string  `json:"validUntil"`
}

// GetRefundGasFeeRequest represents gas fee query request
type GetRefundGasFeeRequest struct {
	Chain    string  `json:"chain"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

// ListBlockchainRefundsOptions represents options for listing refunds
type ListBlockchainRefundsOptions struct {
	WalletID   string `json:"walletId,omitempty"`
	TransferID string `json:"transferId,omitempty"`
	Chain      string `json:"chain,omitempty"`
	Currency   string `json:"currency,omitempty"`
	Status     string `json:"status,omitempty"`
	StartTime  string `json:"startTime,omitempty"`
	EndTime    string `json:"endTime,omitempty"`
	Page       int    `json:"page,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

// BlockchainRefundListResponse represents the response for listing refunds
type BlockchainRefundListResponse struct {
	Refunds    []BlockchainRefund `json:"refunds"`
	TotalCount int                `json:"totalCount"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
}

// CreateBlockchainRefund creates a blockchain refund transaction
// POST /open-api/v3/crypto/refund
func (c *BlockchainRefundClient) CreateBlockchainRefund(ctx context.Context, req *CreateBlockchainRefundRequest) (*BlockchainRefund, error) {
	if req == nil {
		return nil, fmt.Errorf("refund request is required")
	}
	if req.WalletID == "" {
		return nil, fmt.Errorf("wallet ID is required")
	}
	if req.TransferID == "" {
		return nil, fmt.Errorf("transfer ID is required")
	}
	if req.Chain == "" {
		return nil, fmt.Errorf("chain is required")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}
	if req.ToAddress == "" {
		return nil, fmt.Errorf("to address is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/crypto/refund",
		Body:        req,
		RequireAuth: true,
	}

	var refund BlockchainRefund
	err := c.httpClient.DoRequest(ctx, opts, &refund)
	if err != nil {
		return nil, fmt.Errorf("failed to create blockchain refund: %w", err)
	}

	return &refund, nil
}

// ListBlockchainRefunds retrieves all blockchain refunds with optional filtering
// GET /open-api/v3/crypto/refunds
func (c *BlockchainRefundClient) ListBlockchainRefunds(ctx context.Context, options *ListBlockchainRefundsOptions) (*BlockchainRefundListResponse, error) {
	var queryParams url.Values
	if options != nil {
		queryParams = url.Values{}
		if options.WalletID != "" {
			queryParams.Set("walletId", options.WalletID)
		}
		if options.TransferID != "" {
			queryParams.Set("transferId", options.TransferID)
		}
		if options.Chain != "" {
			queryParams.Set("chain", options.Chain)
		}
		if options.Currency != "" {
			queryParams.Set("currency", options.Currency)
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
		Endpoint:    "/open-api/v3/crypto/refunds",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var resp BlockchainRefundListResponse
	err := c.httpClient.DoRequest(ctx, opts, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to list blockchain refunds: %w", err)
	}

	return &resp, nil
}

// GetRefundGasFee retrieves estimated gas fee for a refund
// POST /open-api/v3/crypto/refund/gas-fee
func (c *BlockchainRefundClient) GetRefundGasFee(ctx context.Context, req *GetRefundGasFeeRequest) (*RefundGasFee, error) {
	if req == nil {
		return nil, fmt.Errorf("gas fee request is required")
	}
	if req.Chain == "" {
		return nil, fmt.Errorf("chain is required")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/crypto/refund/gas-fee",
		Body:        req,
		RequireAuth: true,
	}

	var gasFee RefundGasFee
	err := c.httpClient.DoRequest(ctx, opts, &gasFee)
	if err != nil {
		return nil, fmt.Errorf("failed to get refund gas fee: %w", err)
	}

	return &gasFee, nil
}

// GetBlockchainRefund retrieves a specific blockchain refund by ID
// GET /open-api/v3/crypto/refund/{refundId}
func (c *BlockchainRefundClient) GetBlockchainRefund(ctx context.Context, refundID string) (*BlockchainRefund, error) {
	if refundID == "" {
		return nil, fmt.Errorf("refund ID is required")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/crypto/refund/%s", refundID),
		RequireAuth: true,
	}

	var refund BlockchainRefund
	err := c.httpClient.DoRequest(ctx, opts, &refund)
	if err != nil {
		return nil, fmt.Errorf("failed to get blockchain refund: %w", err)
	}

	return &refund, nil
}
