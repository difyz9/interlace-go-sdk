package interlace

import (
	"context"
	"fmt"
)

// PhysicalCardClient handles physical card operations
type PhysicalCardClient struct {
	client *HTTPClient
}

// NewPhysicalCardClient creates a new physical card client
func NewPhysicalCardClient(client *HTTPClient) *PhysicalCardClient {
	return &PhysicalCardClient{client: client}
}

// PhysicalCardFee represents physical card fee information
type PhysicalCardFee struct {
	BinID              string  `json:"binId"`
	CardBrand          string  `json:"cardBrand"`
	Currency           string  `json:"currency"`
	ShippingFee        float64 `json:"shippingFee"`
	CardProductionFee  float64 `json:"cardProductionFee"`
	TotalFee           float64 `json:"totalFee"`
	EstimatedDelivery  string  `json:"estimatedDelivery"`
}

// ShippingAddress represents shipping address for physical card
type ShippingAddress struct {
	RecipientName string `json:"recipientName"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2,omitempty"`
	City          string `json:"city"`
	State         string `json:"state,omitempty"`
	PostalCode    string `json:"postalCode"`
	Country       string `json:"country"`
	PhoneNumber   string `json:"phoneNumber"`
}

// BulkShipRequest represents bulk shipping request for physical cards
type BulkShipRequest struct {
	CardIDs         []string         `json:"cardIds"`
	ShippingAddress *ShippingAddress `json:"shippingAddress"`
	ShippingMethod  string           `json:"shippingMethod,omitempty"` // STANDARD, EXPRESS
	Notes           string           `json:"notes,omitempty"`
}

// BulkShipResponse represents bulk shipping response
type BulkShipResponse struct {
	ShipmentID      string   `json:"shipmentId"`
	CardIDs         []string `json:"cardIds"`
	TrackingNumber  string   `json:"trackingNumber"`
	Status          string   `json:"status"`
	EstimatedDelivery string `json:"estimatedDelivery"`
	CreatedAt       string   `json:"createdAt"`
}

// ConfirmCardholderIdentityRequest represents cardholder identity confirmation request
type ConfirmCardholderIdentityRequest struct {
	CardholderID   string `json:"cardholderId"`
	VerificationID string `json:"verificationId"`
	Verified       bool   `json:"verified"`
	Notes          string `json:"notes,omitempty"`
}

// ConfirmCardholderIdentityResponse represents cardholder identity confirmation response
type ConfirmCardholderIdentityResponse struct {
	CardholderID string `json:"cardholderId"`
	Status       string `json:"status"`
	VerifiedAt   string `json:"verifiedAt"`
	Message      string `json:"message"`
}

// CardholderIdentityURLResponse represents cardholder identity URL response
type CardholderIdentityURLResponse struct {
	CardholderID   string `json:"cardholderId"`
	VerificationID string `json:"verificationId"`
	IdentityURL    string `json:"identityUrl"`
	ExpiresAt      string `json:"expiresAt"`
	Status         string `json:"status"`
}

// ActivatePhysicalCardRequest represents physical card activation request
type ActivatePhysicalCardRequest struct {
	CardID      string `json:"cardId"`
	LastFourDigits string `json:"lastFourDigits"`
	CVV         string `json:"cvv"`
	ExpiryMonth string `json:"expiryMonth"`
	ExpiryYear  string `json:"expiryYear"`
}

// ActivatePhysicalCardResponse represents physical card activation response
type ActivatePhysicalCardResponse struct {
	CardID      string `json:"cardId"`
	Status      string `json:"status"`
	ActivatedAt string `json:"activatedAt"`
	Message     string `json:"message"`
}

// ListPhysicalCardFees lists all physical card fees
func (c *PhysicalCardClient) ListPhysicalCardFees(ctx context.Context) ([]PhysicalCardFee, error) {
	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    "/open-api/v3/physical-card/fees",
		RequireAuth: true,
	}

	var fees []PhysicalCardFee
	err := c.client.DoRequest(ctx, opts, &fees)
	if err != nil {
		return nil, fmt.Errorf("failed to list physical card fees: %w", err)
	}
	return fees, nil
}

// BulkShipPhysicalCards ships multiple physical cards in bulk
func (c *PhysicalCardClient) BulkShipPhysicalCards(ctx context.Context, req *BulkShipRequest) (*BulkShipResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("bulk ship request is required")
	}
	if len(req.CardIDs) == 0 {
		return nil, fmt.Errorf("at least one card ID is required")
	}
	if req.ShippingAddress == nil {
		return nil, fmt.Errorf("shipping address is required")
	}
	if req.ShippingAddress.RecipientName == "" {
		return nil, fmt.Errorf("recipient name is required")
	}
	if req.ShippingAddress.AddressLine1 == "" {
		return nil, fmt.Errorf("address line 1 is required")
	}
	if req.ShippingAddress.City == "" {
		return nil, fmt.Errorf("city is required")
	}
	if req.ShippingAddress.PostalCode == "" {
		return nil, fmt.Errorf("postal code is required")
	}
	if req.ShippingAddress.Country == "" {
		return nil, fmt.Errorf("country is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/physical-card/bulk-ship",
		Body:        req,
		RequireAuth: true,
	}

	var resp BulkShipResponse
	err := c.client.DoRequest(ctx, opts, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk ship physical cards: %w", err)
	}
	return &resp, nil
}

// ConfirmCardholderIdentity confirms cardholder identity for physical card
func (c *PhysicalCardClient) ConfirmCardholderIdentity(ctx context.Context, req *ConfirmCardholderIdentityRequest) (*ConfirmCardholderIdentityResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("confirm identity request is required")
	}
	if req.CardholderID == "" {
		return nil, fmt.Errorf("cardholder ID is required")
	}
	if req.VerificationID == "" {
		return nil, fmt.Errorf("verification ID is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/cardholder/confirm-identity",
		Body:        req,
		RequireAuth: true,
	}

	var resp ConfirmCardholderIdentityResponse
	err := c.client.DoRequest(ctx, opts, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm cardholder identity: %w", err)
	}
	return &resp, nil
}

// GenerateCardholderIdentityURL generates a URL for cardholder identity verification
func (c *PhysicalCardClient) GenerateCardholderIdentityURL(ctx context.Context, cardholderID string) (*CardholderIdentityURLResponse, error) {
	if cardholderID == "" {
		return nil, fmt.Errorf("cardholder ID is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    fmt.Sprintf("/open-api/v3/cardholder/%s/identity-url", cardholderID),
		RequireAuth: true,
	}

	var resp CardholderIdentityURLResponse
	err := c.client.DoRequest(ctx, opts, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate cardholder identity URL: %w", err)
	}
	return &resp, nil
}

// ActivatePhysicalCard activates a physical card
func (c *PhysicalCardClient) ActivatePhysicalCard(ctx context.Context, req *ActivatePhysicalCardRequest) (*ActivatePhysicalCardResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("activate card request is required")
	}
	if req.CardID == "" {
		return nil, fmt.Errorf("card ID is required")
	}
	if req.LastFourDigits == "" {
		return nil, fmt.Errorf("last four digits are required")
	}
	if req.CVV == "" {
		return nil, fmt.Errorf("CVV is required")
	}
	if req.ExpiryMonth == "" {
		return nil, fmt.Errorf("expiry month is required")
	}
	if req.ExpiryYear == "" {
		return nil, fmt.Errorf("expiry year is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/physical-card/activate",
		Body:        req,
		RequireAuth: true,
	}

	var resp ActivatePhysicalCardResponse
	err := c.client.DoRequest(ctx, opts, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to activate physical card: %w", err)
	}
	return &resp, nil
}
