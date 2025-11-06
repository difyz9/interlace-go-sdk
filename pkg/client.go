package interlace

import (
	"context"
	"fmt"
)

// Client represents the main Interlace SDK client
type Client struct {
	config          *Config
	httpClient      *HTTPClient
	OAuth           *OAuthClient
	Account         *AccountClient
	File            *FileClient
	KYC             *KYCClient
	Card            *CardClient
	CardTransaction *CardTransactionClient
	Budget          *BudgetClient
	Payout          *PayoutClient
	Wallet          *WalletClient
	Transfer        *TransferClient
	Payment         *PaymentClient
	Cardholder      *CardholderClient
	CardBin         *CardBinClient
	Common          *CommonClient
	PhysicalCard      *PhysicalCardClient
	Security          *SecurityClient
	Convert           *ConvertClient
	Iframe            *IframeClient
	BlockchainRefund  *BlockchainRefundClient
	BusinessTransfer  *BusinessTransferClient
	InfinityAccount   *InfinityAccountClient
	Sweeping          *SweepingClient
	Testing           *TestingClient
	BusinessAccount   *BusinessAccountClient
}

// NewClient creates a new Interlace SDK client
func NewClient(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	httpClient := NewHTTPClient(config, "")

	client := &Client{
		config:     config,
		httpClient: httpClient,
	}

	// Initialize sub-clients
	client.OAuth = NewOAuthClient(httpClient)
	client.Account = NewAccountClient(httpClient)
	client.File = NewFileClient(httpClient)
	client.KYC = NewKYCClient(httpClient)
	client.Card = NewCardClient(httpClient)
	client.CardTransaction = NewCardTransactionClient(httpClient)
	client.Budget = NewBudgetClient(httpClient)
	client.Payout = NewPayoutClient(httpClient)
	client.Wallet = NewWalletClient(httpClient)
	client.Transfer = NewTransferClient(httpClient)
	client.Payment = NewPaymentClient(httpClient)
	client.Cardholder = NewCardholderClient(httpClient)
	client.CardBin = NewCardBinClient(httpClient)
	client.Common = NewCommonClient(httpClient)
	client.PhysicalCard = NewPhysicalCardClient(httpClient)
	client.Security = NewSecurityClient(httpClient)
	client.Convert = NewConvertClient(httpClient)
	client.Iframe = NewIframeClient(httpClient)
	client.BlockchainRefund = NewBlockchainRefundClient(httpClient)
	client.BusinessTransfer = NewBusinessTransferClient(httpClient)
	client.InfinityAccount = NewInfinityAccountClient(httpClient)
	client.Sweeping = NewSweepingClient(httpClient)
	client.Testing = NewTestingClient(httpClient)
	client.BusinessAccount = NewBusinessAccountClient(httpClient)

	return client
}

// NewClientWithToken creates a new Interlace SDK client with an existing access token
func NewClientWithToken(config *Config, accessToken string) *Client {
	client := NewClient(config)
	client.SetAccessToken(accessToken)
	return client
}

// SetAccessToken sets the access token for authenticated requests
func (c *Client) SetAccessToken(accessToken string) {
	c.httpClient.SetAccessToken(accessToken)
}

// GetAccessToken returns the current access token
func (c *Client) GetAccessToken() string {
	return c.httpClient.GetAccessToken()
}

// Authenticate performs the full OAuth flow and sets the access token
func (c *Client) Authenticate(ctx context.Context, clientID string) (*OAuthTokenData, error) {
	tokenData, err := c.OAuth.AuthorizeAndGetToken(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	c.SetAccessToken(tokenData.AccessToken)
	return tokenData, nil
}

// IsAuthenticated checks if the client has an access token
func (c *Client) IsAuthenticated() bool {
	return c.GetAccessToken() != ""
}

// Config returns the client configuration
func (c *Client) Config() *Config {
	return c.config
}

// SetConfig updates the client configuration
func (c *Client) SetConfig(config *Config) {
	c.config = config
	accessToken := c.GetAccessToken()
	
	// Update HTTP client and sub-clients with new config
	c.httpClient = NewHTTPClient(config, accessToken)
	c.OAuth = NewOAuthClient(c.httpClient)
	c.Account = NewAccountClient(c.httpClient)
	c.File = NewFileClient(c.httpClient)
	c.KYC = NewKYCClient(c.httpClient)
	c.Card = NewCardClient(c.httpClient)
	c.CardTransaction = NewCardTransactionClient(c.httpClient)
	c.Budget = NewBudgetClient(c.httpClient)
	c.Payout = NewPayoutClient(c.httpClient)
	c.Wallet = NewWalletClient(c.httpClient)
	c.Transfer = NewTransferClient(c.httpClient)
	c.Payment = NewPaymentClient(c.httpClient)
	c.Cardholder = NewCardholderClient(c.httpClient)
	c.CardBin = NewCardBinClient(c.httpClient)
	c.Common = NewCommonClient(c.httpClient)
	c.PhysicalCard = NewPhysicalCardClient(c.httpClient)
	c.Security = NewSecurityClient(c.httpClient)
	c.Convert = NewConvertClient(c.httpClient)
	c.Iframe = NewIframeClient(c.httpClient)
	c.BlockchainRefund = NewBlockchainRefundClient(c.httpClient)
	c.BusinessTransfer = NewBusinessTransferClient(c.httpClient)
	c.InfinityAccount = NewInfinityAccountClient(c.httpClient)
	c.Sweeping = NewSweepingClient(c.httpClient)
	c.Testing = NewTestingClient(c.httpClient)
	c.BusinessAccount = NewBusinessAccountClient(c.httpClient)
}

// SetBaseURL is a convenience method to update just the base URL
func (c *Client) SetBaseURL(baseURL string) {
	c.config.BaseURL = baseURL
	c.SetConfig(c.config)
}

// SetClientID is a convenience method to update the client ID in config
func (c *Client) SetClientID(clientID string) {
	c.config.ClientID = clientID
}

// GetClientID returns the client ID from config
func (c *Client) GetClientID() string {
	return c.config.ClientID
}

// QuickSetup is a convenience method for common initialization pattern
// It performs authentication and returns a ready-to-use client
func QuickSetup(clientID string, config *Config) (*Client, *OAuthTokenData, error) {
	if config == nil {
		config = DefaultConfig()
	}
	config.ClientID = clientID

	client := NewClient(config)
	
	ctx := context.Background()
	tokenData, err := client.Authenticate(ctx, clientID)
	if err != nil {
		return nil, nil, fmt.Errorf("quick setup failed: %w", err)
	}

	return client, tokenData, nil
}

// ProductionConfig returns a configuration for the production environment
func ProductionConfig() *Config {
	config := DefaultConfig()
	config.BaseURL = "https://api.interlace.money" // Production URL
	return config
}

// SandboxConfig returns a configuration for the sandbox environment (alias for DefaultConfig)
func SandboxConfig() *Config {
	return DefaultConfig()
}