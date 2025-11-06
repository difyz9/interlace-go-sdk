package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// CardBinClient handles card BIN-related API calls
type CardBinClient struct {
	httpClient *HTTPClient
}

// NewCardBinClient creates a new card BIN client
func NewCardBinClient(httpClient *HTTPClient) *CardBinClient {
	return &CardBinClient{
		httpClient: httpClient,
	}
}

// ListCardBins retrieves a list of all available card BINs
// GET /open-api/v3/card/bins
func (c *CardBinClient) ListCardBins(ctx context.Context, accountID string) (*CardBinListResponse, error) {
	if accountID == "" {
		return nil, fmt.Errorf("accountId cannot be empty")
	}

	params := url.Values{}
	params.Set("accountId", accountID)

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    "/open-api/v3/card/bins?" + params.Encode(),
		RequireAuth: true,
	}

	var response CardBinListResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list card bins: %w", err)
	}

	return &response, nil
}

// ListCardBinsMaintain retrieves a list of card BINs under maintenance
// GET /open-api/v3/card/bins/maintain
func (c *CardBinClient) ListCardBinsMaintain(ctx context.Context, accountID string) (*CardBinListResponse, error) {
	if accountID == "" {
		return nil, fmt.Errorf("accountId cannot be empty")
	}

	params := url.Values{}
	params.Set("accountId", accountID)

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    "/open-api/v3/card/bins/maintain?" + params.Encode(),
		RequireAuth: true,
	}

	var response CardBinListResponse
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list card bins under maintenance: %w", err)
	}

	return &response, nil
}
