package interlace

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// AccountClient handles account operations
type AccountClient struct {
	httpClient *HTTPClient
}

// NewAccountClient creates a new account client
func NewAccountClient(httpClient *HTTPClient) *AccountClient {
	return &AccountClient{
		httpClient: httpClient,
	}
}



// Register creates a new account
func (c *AccountClient) Register(ctx context.Context, req *AccountRegisterRequest) (*AccountData, error) {
	var registerResp AccountRegisterResponse
	err := c.httpClient.DoPostRequest(ctx, "/open-api/v3/accounts/register", req, &registerResp)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if registerResp.Code != "000000" {
		return nil, &Error{
			Code:    registerResp.Code,
			Message: registerResp.Message,
		}
	}

	return &registerResp.Data, nil
}

// RegisterWithDetails creates a new account with specific field ordering as expected by the API
// This function matches the exact curl command structure you provided
func (c *AccountClient) RegisterWithDetails(ctx context.Context, phoneCountryCode, phoneNumber, email, name string) (*AccountData, error) {
	// Create request with exact field order from the curl command
	req := &AccountRegisterRequest{
		PhoneCountryCode: phoneCountryCode,
		PhoneNumber:      phoneNumber,
		Email:           email,
		Name:            name,
	}

	return c.Register(ctx, req)
}

// RegisterGolangTest creates a test account with the exact data from your curl command
func (c *AccountClient) RegisterGolangTest(ctx context.Context) (*AccountData, error) {
	req := &AccountRegisterRequest{
		PhoneCountryCode: "86",
		PhoneNumber:      "15900000031",
		Email:            "15900000031@qq.com",
		Name:             "golang_test",
	}

	return c.Register(ctx, req)
}

// AccountListOptions represents options for listing accounts
type AccountListOptions struct {
	AccountID string // Filter by specific account ID
	Limit     int    // Number of results per page (default: 10, max: 100)
	Page      int    // Page number (default: 1)
	Status    string // Filter by account status (ACTIVE, INACTIVE, etc.)
	Type      int    // Filter by account type
}

// List retrieves a list of accounts with optional filtering
func (c *AccountClient) List(ctx context.Context, opts *AccountListOptions) (*AccountListData, error) {
	// Build query parameters
	params := url.Values{}
	
	if opts != nil {
		if opts.AccountID != "" {
			params.Add("accountId", opts.AccountID)
		}
		if opts.Limit > 0 {
			// Ensure limit doesn't exceed maximum
			if opts.Limit > 100 {
				params.Add("limit", "100")
			} else {
				params.Add("limit", strconv.Itoa(opts.Limit))
			}
		} else {
			params.Add("limit", "10") // Default limit
		}
		if opts.Page > 0 {
			params.Add("page", strconv.Itoa(opts.Page))
		} else {
			params.Add("page", "1") // Default page
		}
		if opts.Status != "" {
			params.Add("status", opts.Status)
		}
		if opts.Type > 0 {
			params.Add("type", strconv.Itoa(opts.Type))
		}
	} else {
		// Default values when no options provided
		params.Add("limit", "10")
		params.Add("page", "1")
	}

	var listResp AccountListResponse
	err := c.httpClient.DoGetRequest(ctx, "/open-api/v3/accounts", params, &listResp)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if listResp.Code != "000000" {
		return nil, &Error{
			Code:    listResp.Code,
			Message: listResp.Message,
		}
	}

	return &listResp.Data, nil
}

// Get retrieves a specific account by ID
func (c *AccountClient) Get(ctx context.Context, accountID string) (*AccountData, error) {
	opts := &AccountListOptions{
		AccountID: accountID,
		Limit:     1,
		Page:      1,
	}

	listData, err := c.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	if len(listData.List) == 0 {
		return nil, &Error{
			Code:    "ACCOUNT_NOT_FOUND",
			Message: fmt.Sprintf("Account with ID %s not found", accountID),
		}
	}

	return &listData.List[0], nil
}

// ListAll retrieves all accounts without pagination (automatically handles pagination)
func (c *AccountClient) ListAll(ctx context.Context) ([]AccountData, error) {
	var allAccounts []AccountData
	page := 1
	limit := 100 // Use maximum limit for efficiency

	for {
		opts := &AccountListOptions{
			Limit: limit,
			Page:  page,
		}

		accountList, err := c.List(ctx, opts)
		if err != nil {
			return nil, err
		}

		allAccounts = append(allAccounts, accountList.List...)

		// If we got less than the limit, we've reached the end
		if len(accountList.List) < limit {
			break
		}

		page++
	}

	return allAccounts, nil
}

// ListByStatus retrieves accounts filtered by status
func (c *AccountClient) ListByStatus(ctx context.Context, status string) ([]AccountData, error) {
	opts := &AccountListOptions{
		Status: status,
		Limit:  100,
		Page:   1,
	}

	accountList, err := c.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return accountList.List, nil
}

// ListActiveAccounts retrieves all active accounts
func (c *AccountClient) ListActiveAccounts(ctx context.Context) ([]AccountData, error) {
	return c.ListByStatus(ctx, "ACTIVE")
}

// ListInactiveAccounts retrieves all inactive accounts
func (c *AccountClient) ListInactiveAccounts(ctx context.Context) ([]AccountData, error) {
	return c.ListByStatus(ctx, "INACTIVE")
}

// ListByType retrieves accounts filtered by type
func (c *AccountClient) ListByType(ctx context.Context, accountType int) ([]AccountData, error) {
	opts := &AccountListOptions{
		Type:  accountType,
		Limit: 100,
		Page:  1,
	}

	accountList, err := c.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return accountList.List, nil
}

// Count returns the total number of accounts
func (c *AccountClient) Count(ctx context.Context) (int, error) {
	opts := &AccountListOptions{
		Limit: 1,
		Page:  1,
	}

	accountList, err := c.List(ctx, opts)
	if err != nil {
		return 0, err
	}

	// Convert total string to int
	total, err := strconv.Atoi(accountList.Total)
	if err != nil {
		return 0, fmt.Errorf("failed to parse total count: %w", err)
	}

	return total, nil
}

// GetAccountsByPage retrieves accounts with specific pagination settings
func (c *AccountClient) GetAccountsByPage(ctx context.Context, page, limit int) (*AccountListData, error) {
	opts := &AccountListOptions{
		Page:  page,
		Limit: limit,
	}

	return c.List(ctx, opts)
}