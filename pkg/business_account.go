package interlace

import (
	"context"
	"fmt"
	"net/url"
)

// BusinessAccountClient handles business account operations
type BusinessAccountClient struct {
	httpClient *HTTPClient
}

// NewBusinessAccountClient creates a new business account client
func NewBusinessAccountClient(httpClient *HTTPClient) *BusinessAccountClient {
	return &BusinessAccountClient{
		httpClient: httpClient,
	}
}

// BusinessAccount represents a business account
type BusinessAccount struct {
	AccountID       string  `json:"accountId"`
	LegalEntityID   string  `json:"legalEntityId"`
	AccountType     string  `json:"accountType"` // CORPORATE, INDIVIDUAL
	Status          string  `json:"status"` // ACTIVE, SUSPENDED, CLOSED
	Currency        string  `json:"currency"`
	Balance         float64 `json:"balance"`
	AvailableBalance float64 `json:"availableBalance"`
	FrozenBalance   float64 `json:"frozenBalance"`
	AccountNumber   string  `json:"accountNumber"`
	IBAN            string  `json:"iban,omitempty"`
	SwiftCode       string  `json:"swiftCode,omitempty"`
	BankName        string  `json:"bankName,omitempty"`
	BankAddress     string  `json:"bankAddress,omitempty"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// BusinessAccountBalance represents account balance information
type BusinessAccountBalance struct {
	AccountID        string  `json:"accountId"`
	Currency         string  `json:"currency"`
	Balance          float64 `json:"balance"`
	AvailableBalance float64 `json:"availableBalance"`
	FrozenBalance    float64 `json:"frozenBalance"`
	PendingBalance   float64 `json:"pendingBalance"`
	LastUpdated      string  `json:"lastUpdated"`
}

// BusinessAccountTransaction represents a business account transaction
type BusinessAccountTransaction struct {
	TransactionID   string  `json:"transactionId"`
	AccountID       string  `json:"accountId"`
	Type            string  `json:"type"` // DEBIT, CREDIT
	Category        string  `json:"category"` // TRANSFER_IN, TRANSFER_OUT, FEE, REFUND, PAYOUT
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	BalanceBefore   float64 `json:"balanceBefore"`
	BalanceAfter    float64 `json:"balanceAfter"`
	Status          string  `json:"status"`
	Description     string  `json:"description,omitempty"`
	CounterpartyName string `json:"counterpartyName,omitempty"`
	CounterpartyAccount string `json:"counterpartyAccount,omitempty"`
	Reference       string  `json:"reference,omitempty"`
	CreatedAt       string  `json:"createdAt"`
}

// LegalEntity represents a legal entity
type LegalEntity struct {
	EntityID           string            `json:"entityId"`
	EntityType         string            `json:"entityType"` // COMPANY, INDIVIDUAL
	CompanyName        string            `json:"companyName,omitempty"`
	RegistrationNumber string            `json:"registrationNumber,omitempty"`
	TaxID              string            `json:"taxId,omitempty"`
	IncorporationDate  string            `json:"incorporationDate,omitempty"`
	Country            string            `json:"country"`
	Address            *Address          `json:"address"`
	ContactPerson      *ContactPerson    `json:"contactPerson,omitempty"`
	Directors          []Director        `json:"directors,omitempty"`
	UltimateBeneficiaryOwners []UBO     `json:"ultimateBeneficiaryOwners,omitempty"`
	Documents          []Document        `json:"documents,omitempty"`
	Status             string            `json:"status"` // PENDING, VERIFIED, REJECTED
	CreatedAt          string            `json:"createdAt"`
	UpdatedAt          string            `json:"updatedAt"`
}

// Address represents a physical address
type Address struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode"`
	Country      string `json:"country"`
}

// ContactPerson represents a contact person
type ContactPerson struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Position    string `json:"position,omitempty"`
}

// Director represents a company director
type Director struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DateOfBirth string `json:"dateOfBirth"`
	Nationality string `json:"nationality"`
	Position    string `json:"position"`
}

// UBO represents an Ultimate Beneficiary Owner
type UBO struct {
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	DateOfBirth    string  `json:"dateOfBirth"`
	Nationality    string  `json:"nationality"`
	OwnershipPercentage float64 `json:"ownershipPercentage"`
}

// Document represents a legal document
type Document struct {
	DocumentID   string `json:"documentId"`
	DocumentType string `json:"documentType"`
	FileID       string `json:"fileId"`
	FileName     string `json:"fileName"`
	UploadedAt   string `json:"uploadedAt"`
}

// CreateLegalEntityRequest represents legal entity creation request
type CreateLegalEntityRequest struct {
	EntityType         string            `json:"entityType"`
	CompanyName        string            `json:"companyName,omitempty"`
	RegistrationNumber string            `json:"registrationNumber,omitempty"`
	TaxID              string            `json:"taxId,omitempty"`
	IncorporationDate  string            `json:"incorporationDate,omitempty"`
	Country            string            `json:"country"`
	Address            *Address          `json:"address"`
	ContactPerson      *ContactPerson    `json:"contactPerson"`
	Directors          []Director        `json:"directors,omitempty"`
	UltimateBeneficiaryOwners []UBO     `json:"ultimateBeneficiaryOwners,omitempty"`
}

// UpdateLegalEntityRequest represents legal entity update request
type UpdateLegalEntityRequest struct {
	CompanyName   string         `json:"companyName,omitempty"`
	Address       *Address       `json:"address,omitempty"`
	ContactPerson *ContactPerson `json:"contactPerson,omitempty"`
	TaxID         string         `json:"taxId,omitempty"`
}

// CreateVirtualAccountRequest represents virtual account creation request
type CreateVirtualAccountRequest struct {
	LegalEntityID string `json:"legalEntityId"`
	Currency      string `json:"currency"`
	AccountType   string `json:"accountType,omitempty"`
	Reference     string `json:"reference,omitempty"`
}

// ListBusinessAccountTransactionsOptions represents options for listing transactions
type ListBusinessAccountTransactionsOptions struct {
	AccountID string `json:"accountId,omitempty"`
	Type      string `json:"type,omitempty"`
	Category  string `json:"category,omitempty"`
	Status    string `json:"status,omitempty"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	MinAmount float64 `json:"minAmount,omitempty"`
	MaxAmount float64 `json:"maxAmount,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// BusinessAccountTransactionListResponse represents the response for listing transactions
type BusinessAccountTransactionListResponse struct {
	Transactions []BusinessAccountTransaction `json:"transactions"`
	TotalCount   int                          `json:"totalCount"`
	Page         int                          `json:"page"`
	Limit        int                          `json:"limit"`
}

// GetBusinessAccounts retrieves all business accounts
// GET /open-api/v3/business/accounts
func (c *BusinessAccountClient) GetBusinessAccounts(ctx context.Context, legalEntityID string) ([]BusinessAccount, error) {
	queryParams := url.Values{}
	if legalEntityID != "" {
		queryParams.Set("legalEntityId", legalEntityID)
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    "/open-api/v3/business/accounts",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var accounts []BusinessAccount
	err := c.httpClient.DoRequest(ctx, opts, &accounts)
	if err != nil {
		return nil, fmt.Errorf("failed to get business accounts: %w", err)
	}

	return accounts, nil
}

// GetAccountBalance retrieves the balance of a business account
// GET /open-api/v3/business/account/{accountId}/balance
func (c *BusinessAccountClient) GetAccountBalance(ctx context.Context, accountID string) (*BusinessAccountBalance, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/business/account/%s/balance", accountID),
		RequireAuth: true,
	}

	var balance BusinessAccountBalance
	err := c.httpClient.DoRequest(ctx, opts, &balance)
	if err != nil {
		return nil, fmt.Errorf("failed to get account balance: %w", err)
	}

	return &balance, nil
}

// GetAccountTransactions retrieves transactions for a business account
// GET /open-api/v3/business/account/transactions
func (c *BusinessAccountClient) GetAccountTransactions(ctx context.Context, options *ListBusinessAccountTransactionsOptions) (*BusinessAccountTransactionListResponse, error) {
	var queryParams url.Values
	if options != nil {
		queryParams = url.Values{}
		if options.AccountID != "" {
			queryParams.Set("accountId", options.AccountID)
		}
		if options.Type != "" {
			queryParams.Set("type", options.Type)
		}
		if options.Category != "" {
			queryParams.Set("category", options.Category)
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
		if options.MinAmount > 0 {
			queryParams.Set("minAmount", fmt.Sprintf("%.2f", options.MinAmount))
		}
		if options.MaxAmount > 0 {
			queryParams.Set("maxAmount", fmt.Sprintf("%.2f", options.MaxAmount))
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
		Endpoint:    "/open-api/v3/business/account/transactions",
		QueryParams: queryParams,
		RequireAuth: true,
	}

	var resp BusinessAccountTransactionListResponse
	err := c.httpClient.DoRequest(ctx, opts, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get account transactions: %w", err)
	}

	return &resp, nil
}

// CreateLegalEntity creates a new legal entity
// POST /open-api/v3/business/legal-entity
func (c *BusinessAccountClient) CreateLegalEntity(ctx context.Context, req *CreateLegalEntityRequest) (*LegalEntity, error) {
	if req == nil {
		return nil, fmt.Errorf("legal entity request is required")
	}
	if req.EntityType == "" {
		return nil, fmt.Errorf("entity type is required")
	}
	if req.Country == "" {
		return nil, fmt.Errorf("country is required")
	}
	if req.Address == nil {
		return nil, fmt.Errorf("address is required")
	}
	if req.ContactPerson == nil {
		return nil, fmt.Errorf("contact person is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/business/legal-entity",
		Body:        req,
		RequireAuth: true,
	}

	var entity LegalEntity
	err := c.httpClient.DoRequest(ctx, opts, &entity)
	if err != nil {
		return nil, fmt.Errorf("failed to create legal entity: %w", err)
	}

	return &entity, nil
}

// GetLegalEntity retrieves a legal entity by ID
// GET /open-api/v3/business/legal-entity/{entityId}
func (c *BusinessAccountClient) GetLegalEntity(ctx context.Context, entityID string) (*LegalEntity, error) {
	if entityID == "" {
		return nil, fmt.Errorf("entity ID is required")
	}

	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    fmt.Sprintf("/open-api/v3/business/legal-entity/%s", entityID),
		RequireAuth: true,
	}

	var entity LegalEntity
	err := c.httpClient.DoRequest(ctx, opts, &entity)
	if err != nil {
		return nil, fmt.Errorf("failed to get legal entity: %w", err)
	}

	return &entity, nil
}

// UpdateLegalEntity updates an existing legal entity
// PUT /open-api/v3/business/legal-entity/{entityId}
func (c *BusinessAccountClient) UpdateLegalEntity(ctx context.Context, entityID string, req *UpdateLegalEntityRequest) (*LegalEntity, error) {
	if entityID == "" {
		return nil, fmt.Errorf("entity ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("update request is required")
	}

	opts := &RequestOptions{
		Method:      "PUT",
		Endpoint:    fmt.Sprintf("/open-api/v3/business/legal-entity/%s", entityID),
		Body:        req,
		RequireAuth: true,
	}

	var entity LegalEntity
	err := c.httpClient.DoRequest(ctx, opts, &entity)
	if err != nil {
		return nil, fmt.Errorf("failed to update legal entity: %w", err)
	}

	return &entity, nil
}

// CreateVirtualAccount creates a virtual bank account for a legal entity
// POST /open-api/v3/business/virtual-account
func (c *BusinessAccountClient) CreateVirtualAccount(ctx context.Context, req *CreateVirtualAccountRequest) (*BusinessAccount, error) {
	if req == nil {
		return nil, fmt.Errorf("virtual account request is required")
	}
	if req.LegalEntityID == "" {
		return nil, fmt.Errorf("legal entity ID is required")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/business/virtual-account",
		Body:        req,
		RequireAuth: true,
	}

	var account BusinessAccount
	err := c.httpClient.DoRequest(ctx, opts, &account)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual account: %w", err)
	}

	return &account, nil
}
