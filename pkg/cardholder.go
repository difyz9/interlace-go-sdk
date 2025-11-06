package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// CardholderClient handles cardholder-related API calls
type CardholderClient struct {
	httpClient *HTTPClient
}

// NewCardholderClient creates a new cardholder client
func NewCardholderClient(httpClient *HTTPClient) *CardholderClient {
	return &CardholderClient{
		httpClient: httpClient,
	}
}

// CreateCardholderRequest represents the request to create a cardholder
type CreateCardholderRequest struct {
	AccountID          string                  `json:"accountId"`
	BinID              string                  `json:"binId"` // Card BIN ID
	FirstName          string                  `json:"firstName"`
	LastName           string                  `json:"lastName"`
	Email              string                  `json:"email"`
	PhoneNumber        string                  `json:"phoneNumber"`
	PhoneCountryCode   string                  `json:"phoneCountryCode"`
	DateOfBirth        string                  `json:"dateOfBirth"` // Format: YYYY-MM-DD
	Nationality        string                  `json:"nationality"` // ISO 3166-1 alpha-2 country code
	Gender             string                  `json:"gender,omitempty"`
	Occupation         string                  `json:"occupation,omitempty"`
	Address            *CardholderAddress      `json:"address,omitempty"`
	IdentityDocument   *IdentityDocument       `json:"identityDocument,omitempty"`
	IdempotencyKey     string                  `json:"idempotencyKey,omitempty"`
}

// CardholderAddress represents an address
type CardholderAddress struct {
	Country    string `json:"country"` // ISO 3166-1 alpha-2 country code
	State      string `json:"state,omitempty"`
	City       string `json:"city"`
	PostalCode string `json:"postalCode,omitempty"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
}

// IdentityDocument represents identity document information
type IdentityDocument struct {
	Type           string `json:"type"`           // CN-RIC, HK-HKID, PASSPORT, DLN, etc.
	Number         string `json:"number"`
	IssuingCountry string `json:"issuingCountry"` // ISO 3166-1 alpha-2 country code
	ExpiryDate     string `json:"expiryDate,omitempty"` // Format: YYYY-MM-DD
}

// UpdateCardholderRequest represents the request to update a cardholder
type UpdateCardholderRequest struct {
	Email            string             `json:"email,omitempty"`
	PhoneNumber      string             `json:"phoneNumber,omitempty"`
	PhoneCountryCode string             `json:"phoneCountryCode,omitempty"`
	Address          *CardholderAddress `json:"address,omitempty"`
	Occupation       string             `json:"occupation,omitempty"`
}

// CardholderListOptions represents options for listing cardholders
type CardholderListOptions struct {
	AccountID string `json:"accountId,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// CreateCardholder creates a new cardholder
// POST /open-api/v3/cardholders
func (c *CardholderClient) CreateCardholder(ctx context.Context, req *CreateCardholderRequest) (*Cardholder, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.AccountID == "" || req.BinID == "" || req.FirstName == "" || req.LastName == "" || req.Email == "" {
		return nil, fmt.Errorf("accountId, binId, firstName, lastName and email are required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/cardholders",
		Body:        req,
		RequireAuth: true,
	}

	var response Cardholder
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create cardholder: %w", err)
	}

	return &response, nil
}

// ListCardholders retrieves a list of cardholders
// GET /open-api/v3/cardholders
func (c *CardholderClient) ListCardholders(ctx context.Context, opts *CardholderListOptions) (*CardholderListResponse, error) {
	params := url.Values{}
	if opts != nil {
		if opts.AccountID != "" {
			params.Set("accountId", opts.AccountID)
		}
		if opts.Page > 0 {
			params.Set("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.Limit > 0 {
			params.Set("limit", fmt.Sprintf("%d", opts.Limit))
		}
	}

	endpoint := "/open-api/v3/cardholders"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	reqOpts := &RequestOptions{
		Method:      "GET",
		Endpoint:    endpoint,
		RequireAuth: true,
	}

	var response CardholderListResponse
	err := c.httpClient.DoRequest(ctx, reqOpts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list cardholders: %w", err)
	}

	return &response, nil
}

// GetCardholder retrieves a specific cardholder by ID
// GET /open-api/v3/cardholders/{id}
func (c *CardholderClient) GetCardholder(ctx context.Context, cardholderID string) (*Cardholder, error) {
	if cardholderID == "" {
		return nil, fmt.Errorf("cardholder ID cannot be empty")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/cardholders/%s", cardholderID),
		RequireAuth: true,
	}

	var response Cardholder
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get cardholder: %w", err)
	}

	return &response, nil
}

// UpdateCardholder updates a cardholder's information
// PATCH /open-api/v3/cardholders/{id}
func (c *CardholderClient) UpdateCardholder(ctx context.Context, cardholderID string, req *UpdateCardholderRequest) (*Cardholder, error) {
	if cardholderID == "" {
		return nil, fmt.Errorf("cardholder ID cannot be empty")
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	opts := &RequestOptions{
		Method:      "PATCH",
		Endpoint:    fmt.Sprintf("/open-api/v3/cardholders/%s", cardholderID),
		Body:        req,
		RequireAuth: true,
	}

	var response Cardholder
	err := c.httpClient.DoRequest(ctx, opts, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update cardholder: %w", err)
	}

	return &response, nil
}
