package interlace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	assert.Equal(t, "https://api-sandbox.interlace.money", config.BaseURL)
	assert.Equal(t, "interlace-go-sdk/1.0.0", config.UserAgent)
	assert.NotZero(t, config.Timeout)
}

func TestProductionConfig(t *testing.T) {
	config := ProductionConfig()
	
	assert.Equal(t, "https://api.interlace.money", config.BaseURL)
}

func TestNewClient(t *testing.T) {
	client := NewClient(nil)
	
	assert.NotNil(t, client)
	assert.NotNil(t, client.OAuth)
	assert.NotNil(t, client.Account)
	assert.NotNil(t, client.File)
	assert.NotNil(t, client.KYC)
	assert.False(t, client.IsAuthenticated())
}

func TestClientSetAccessToken(t *testing.T) {
	client := NewClient(nil)
	token := "test-token"
	
	client.SetAccessToken(token)
	
	assert.Equal(t, token, client.GetAccessToken())
	assert.True(t, client.IsAuthenticated())
}

func TestAccountRegisterRequest(t *testing.T) {
	req := &AccountRegisterRequest{
		PhoneNumber:      "15900000042",
		Email:            "test@example.com",
		Name:             "Test User",
		PhoneCountryCode: "86",
	}
	
	assert.Equal(t, "15900000042", req.PhoneNumber)
	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "Test User", req.Name)
	assert.Equal(t, "86", req.PhoneCountryCode)
}

func TestError(t *testing.T) {
	err := &Error{
		Code:    "TEST_ERROR",
		Message: "This is a test error",
	}
	
	expectedMessage := "Interlace API Error - Code: TEST_ERROR, Message: This is a test error"
	assert.Equal(t, expectedMessage, err.Error())
}

// Example test for OAuth client
func TestOAuthClient(t *testing.T) {
	config := DefaultConfig()
	oauthClient := NewOAuthClient(config)
	
	assert.NotNil(t, oauthClient)
	assert.NotNil(t, oauthClient.httpClient)
}

// Example test for Account client
func TestAccountClient(t *testing.T) {
	config := DefaultConfig()
	token := "test-token"
	accountClient := NewAccountClient(config, token)
	
	assert.NotNil(t, accountClient)
	assert.Equal(t, token, accountClient.httpClient.GetAccessToken())
	
	// Test SetAccessToken
	newToken := "new-test-token"
	accountClient.SetAccessToken(newToken)
	assert.Equal(t, newToken, accountClient.httpClient.GetAccessToken())
}

// Example test for File client
func TestFileClient(t *testing.T) {
	config := DefaultConfig()
	token := "test-token"
	fileClient := NewFileClient(config, token)
	
	assert.NotNil(t, fileClient)
	assert.Equal(t, token, fileClient.httpClient.GetAccessToken())
}

func TestHTTPClient(t *testing.T) {
	config := DefaultConfig()
	token := "test-token"
	httpClient := NewHTTPClient(config, token)
	
	assert.NotNil(t, httpClient)
	assert.Equal(t, token, httpClient.GetAccessToken())
	
	// Test SetAccessToken
	newToken := "new-test-token"
	httpClient.SetAccessToken(newToken)
	assert.Equal(t, newToken, httpClient.GetAccessToken())
}

func TestAccountRegisterWithDetails(t *testing.T) {
	config := DefaultConfig()
	token := "test-token"
	accountClient := NewAccountClient(config, token)
	
	assert.NotNil(t, accountClient)
	
	// Test that RegisterWithDetails creates correct request structure
	// Note: This would require mocking the HTTP client for actual testing
}

func TestRegisterGolangTest(t *testing.T) {
	config := DefaultConfig()
	token := "test-token"
	accountClient := NewAccountClient(config, token)
	
	assert.NotNil(t, accountClient)
	
	// Test that RegisterGolangTest has correct preset values
	// Note: This would require mocking the HTTP client for actual testing
}

func TestAccountListOptions(t *testing.T) {
	opts := &AccountListOptions{
		AccountID: "test-account-id",
		Limit:     50,
		Page:      2,
		Status:    AccountStatusActive,
		Type:      AccountTypePersonal,
	}
	
	assert.Equal(t, "test-account-id", opts.AccountID)
	assert.Equal(t, 50, opts.Limit)
	assert.Equal(t, 2, opts.Page)
	assert.Equal(t, "ACTIVE", opts.Status)
	assert.Equal(t, 1, opts.Type)
}

func TestAccountConstants(t *testing.T) {
	// Test account status constants
	assert.Equal(t, "ACTIVE", AccountStatusActive)
	assert.Equal(t, "INACTIVE", AccountStatusInactive)
	assert.Equal(t, "PENDING", AccountStatusPending)
	assert.Equal(t, "SUSPENDED", AccountStatusSuspended)
	
	// Test account type constants
	assert.Equal(t, 1, AccountTypePersonal)
	assert.Equal(t, 2, AccountTypeBusiness)
	assert.Equal(t, 3, AccountTypeChild)
}

func TestKYCConstants(t *testing.T) {
	// Test KYC ID type constants
	assert.Equal(t, "CN-RIC", IDTypeCNRIC)
	assert.Equal(t, "HK-HKID", IDTypeHKHKID)
	assert.Equal(t, "PASSPORT", IDTypePassport)
	assert.Equal(t, "DLN", IDTypeDLN)
	
	// Test KYC status constants
	assert.Equal(t, "PENDING", KYCStatusPending)
	assert.Equal(t, "APPROVED", KYCStatusApproved)
	assert.Equal(t, "REJECTED", KYCStatusRejected)
	assert.Equal(t, "EXPIRED", KYCStatusExpired)
}

func TestKYCBuilder(t *testing.T) {
	builder := NewKYCBuilder()
	
	kycRequest := builder.
		SetPersonalInfo("John", "Doe", "1990-01-15", "M").
		SetNationality("US", "US").
		SetAddress("123 Main St", "New York", "10001", "US").
		SetIDInfo(IDTypePassport, "A12345678").
		SetOccupationInfo("01", "Employment").
		SetAccountPurpose("Personal Banking").
		SetDocumentFiles("front-id", "selfie-id").
		Build()
	
	assert.Equal(t, "John", kycRequest.FirstName)
	assert.Equal(t, "Doe", kycRequest.LastName)
	assert.Equal(t, "1990-01-15", kycRequest.DateOfBirth)
	assert.Equal(t, "M", kycRequest.Gender)
	assert.Equal(t, "US", kycRequest.Nationality)
	assert.Equal(t, "US", kycRequest.CountryOfResidence)
	assert.Equal(t, "123 Main St", kycRequest.Address)
	assert.Equal(t, "New York", kycRequest.City)
	assert.Equal(t, "10001", kycRequest.PostalCode)
	assert.Equal(t, "US", kycRequest.Country)
	assert.Equal(t, IDTypePassport, kycRequest.IDType)
	assert.Equal(t, "A12345678", kycRequest.IDNumber)
	assert.Equal(t, "01", kycRequest.Occupation)
	assert.Equal(t, "Employment", kycRequest.SourceOfIncome)
	assert.Equal(t, "Personal Banking", kycRequest.PurposeOfAccount)
	assert.Equal(t, "front-id", kycRequest.IDFrontImageFileID)
	assert.Equal(t, "selfie-id", kycRequest.SelfieImageFileID)
}

func TestKYCBuilderValidation(t *testing.T) {
	// Test incomplete request validation
	builder := NewKYCBuilder()
	err := builder.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "firstName is required")
	
	// Test complete request validation
	completeBuilder := NewKYCBuilder().
		SetPersonalInfo("John", "Doe", "1990-01-15", "M").
		SetNationality("US", "US").
		SetAddress("123 Main St", "New York", "10001", "US").
		SetIDInfo(IDTypePassport, "A12345678").
		SetOccupationInfo("01", "Employment").
		SetAccountPurpose("Personal Banking").
		SetDocumentFiles("front-id", "selfie-id")
	
	err = completeBuilder.Validate()
	assert.NoError(t, err)
}

func TestKYCClient(t *testing.T) {
	config := DefaultConfig()
	token := "test-token"
	kycClient := NewKYCClient(config, token)
	
	assert.NotNil(t, kycClient)
	assert.Equal(t, token, kycClient.httpClient.GetAccessToken())
	
	// Test SetAccessToken
	newToken := "new-test-token"
	kycClient.SetAccessToken(newToken)
	assert.Equal(t, newToken, kycClient.httpClient.GetAccessToken())
}