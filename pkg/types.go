package interlace

import (
	"encoding/json"
	"fmt"
	"time"
)

// BaseResponse represents the common response structure from Interlace API
type BaseResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// OAuthAuthorizeResponse represents the OAuth authorization response
type OAuthAuthorizeResponse struct {
	BaseResponse
	Data OAuthAuthorizeData `json:"data"`
}

type OAuthAuthorizeData struct {
	Timestamp int64  `json:"timestamp"`
	Code      string `json:"code"`
}

// OAuthTokenRequest represents the OAuth token request
type OAuthTokenRequest struct {
	Code     string `json:"code"`
	ClientID string `json:"clientId"`
}

// OAuthTokenResponse represents the OAuth token response
type OAuthTokenResponse struct {
	BaseResponse
	Data OAuthTokenData `json:"data"`
}

type OAuthTokenData struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	Timestamp    int64  `json:"timestamp"`
}

// OAuthRefreshTokenRequest represents the OAuth refresh token request
type OAuthRefreshTokenRequest struct {
	ClientID     string `json:"clientId"`
	RefreshToken string `json:"refreshToken"`
}

// OAuthRefreshTokenResponse represents the OAuth refresh token response
type OAuthRefreshTokenResponse struct {
	BaseResponse
	Data OAuthRefreshTokenData `json:"data"`
}

type OAuthRefreshTokenData struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
	Timestamp   int64  `json:"timestamp"`
}

// AccountRegisterRequest represents the account registration request
type AccountRegisterRequest struct {
	PhoneNumber      string `json:"phoneNumber"`
	Email            string `json:"email"`
	Name             string `json:"name"`
	PhoneCountryCode string `json:"phoneCountryCode"`
}

// AccountRegisterResponse represents the account registration response
type AccountRegisterResponse struct {
	BaseResponse
	Data AccountData `json:"data"`
}

// AccountListResponse represents the account list response
type AccountListResponse struct {
	BaseResponse
	Data AccountListData `json:"data"`
}

type AccountListData struct {
	List  []AccountData `json:"list"`
	Total string        `json:"total"`
}

type AccountData struct {
	ID                string    `json:"id"`
	CreateTime        string    `json:"createTime"`
	Type              int       `json:"type"`
	Status            string    `json:"status"`
	VerifiedName      string    `json:"verifiedName"`
	VerifiedNameEn    *string   `json:"verifiedNameEn"`
	DisplayID         string    `json:"displayId"`
	ParentAccountID   *string   `json:"parentAccountId,omitempty"`
}

// Account status constants
const (
	AccountStatusActive   = "ACTIVE"
	AccountStatusInactive = "INACTIVE"
	AccountStatusPending  = "PENDING"
	AccountStatusSuspended = "SUSPENDED"
)

// Account type constants
const (
	AccountTypePersonal = 1
	AccountTypeBusiness = 2
	AccountTypeChild    = 3
)

// KYC ID type constants
const (
	IDTypeCNRIC            = "CN-RIC"           // China Resident Identity Card
	IDTypeHKHKID           = "HK-HKID"          // Hong Kong Identity Card
	IDTypePassport         = "PASSPORT"         // Passport
	IDTypeDLN              = "DLN"              // Driver's License Number
	IDTypeGovernmentID     = "Government-Issued ID Card"
)

// KYC status constants
const (
	KYCStatusPending   = "PENDING"
	KYCStatusApproved  = "APPROVED"
	KYCStatusRejected  = "REJECTED"
	KYCStatusExpired   = "EXPIRED"
)

// KYCSubmitRequest represents the KYC submission request
type KYCSubmitRequest struct {
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	MiddleName          string `json:"middleName,omitempty"`
	DateOfBirth         string `json:"dateOfBirth"`         // Format: YYYY-MM-DD
	Gender              string `json:"gender"`              // M or F
	Nationality         string `json:"nationality"`         // Country code
	CountryOfResidence  string `json:"countryOfResidence"`  // Country code
	Address             string `json:"address"`
	City                string `json:"city"`
	State               string `json:"state,omitempty"`
	PostalCode          string `json:"postalCode"`
	Country             string `json:"country"`             // Country code
	IDType              string `json:"idType"`              // See ID type constants
	IDNumber            string `json:"idNumber"`
	IDExpiryDate        string `json:"idExpiryDate,omitempty"` // Format: YYYY-MM-DD
	Occupation          string `json:"occupation"`          // Occupation code
	SourceOfIncome      string `json:"sourceOfIncome"`
	AnnualIncome        string `json:"annualIncome,omitempty"`
	PurposeOfAccount    string `json:"purposeOfAccount"`
	ExpectedTxnVolume   string `json:"expectedTxnVolume,omitempty"`
	IDFrontImageFileID  string `json:"idFrontImageFileId"`  // File ID from upload API
	IDBackImageFileID   string `json:"idBackImageFileId,omitempty"`
	SelfieImageFileID   string `json:"selfieImageFileId"`   // File ID from upload API
}

// KYCSubmitResponse represents the KYC submission response
type KYCSubmitResponse struct {
	BaseResponse
	Data KYCSubmitData `json:"data"`
}

type KYCSubmitData struct {
	KYCApplicationID string `json:"kycApplicationId"`
	Status           string `json:"status"`
	SubmittedTime    string `json:"submittedTime"`
}

// KYCStatusResponse represents the KYC status query response
type KYCStatusResponse struct {
	BaseResponse
	Data KYCStatusData `json:"data"`
}

type KYCStatusData struct {
	AccountID        string `json:"accountId"`
	KYCApplicationID string `json:"kycApplicationId"`
	Status           string `json:"status"`
	SubmittedTime    string `json:"submittedTime,omitempty"`
	ReviewedTime     string `json:"reviewedTime,omitempty"`
	RejectionReason  string `json:"rejectionReason,omitempty"`
	ExpiryTime       string `json:"expiryTime,omitempty"`
}

// CDDDetailResponse represents the CDD (Customer Due Diligence) detail response
type CDDDetailResponse struct {
	BaseResponse
	Data CDDDetailData `json:"data"`
}

type CDDDetailData struct {
	AccountID         string              `json:"accountId"`
	KYCVerification   *KYCVerificationDetail `json:"kycVerification,omitempty"`
	KYBVerification   *KYBVerificationDetail `json:"kybVerification,omitempty"`
	OverallStatus     string              `json:"overallStatus"`
	LastUpdated       string              `json:"lastUpdated"`
}

// KYCVerificationDetail represents individual KYC verification details
type KYCVerificationDetail struct {
	ApplicationID       string                 `json:"applicationId"`
	Status             string                 `json:"status"`
	SubmittedTime      string                 `json:"submittedTime"`
	ReviewedTime       string                 `json:"reviewedTime,omitempty"`
	ExpiryTime         string                 `json:"expiryTime,omitempty"`
	RejectionReason    string                 `json:"rejectionReason,omitempty"`
	PersonalInfo       *PersonalInfo          `json:"personalInfo,omitempty"`
	DocumentInfo       *DocumentInfo          `json:"documentInfo,omitempty"`
	VerificationChecks *VerificationChecks    `json:"verificationChecks,omitempty"`
	RiskAssessment     *RiskAssessment        `json:"riskAssessment,omitempty"`
}

// KYBVerificationDetail represents business KYB verification details
type KYBVerificationDetail struct {
	ApplicationID       string                 `json:"applicationId"`
	Status             string                 `json:"status"`
	SubmittedTime      string                 `json:"submittedTime"`
	ReviewedTime       string                 `json:"reviewedTime,omitempty"`
	ExpiryTime         string                 `json:"expiryTime,omitempty"`
	RejectionReason    string                 `json:"rejectionReason,omitempty"`
	BusinessInfo       *BusinessInfo          `json:"businessInfo,omitempty"`
	ComplianceChecks   *ComplianceChecks      `json:"complianceChecks,omitempty"`
	RiskAssessment     *RiskAssessment        `json:"riskAssessment,omitempty"`
}

// PersonalInfo represents personal information from KYC
type PersonalInfo struct {
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	MiddleName         string `json:"middleName,omitempty"`
	DateOfBirth        string `json:"dateOfBirth"`
	Gender             string `json:"gender"`
	Nationality        string `json:"nationality"`
	CountryOfResidence string `json:"countryOfResidence"`
	Address            string `json:"address"`
	City               string `json:"city"`
	State              string `json:"state,omitempty"`
	PostalCode         string `json:"postalCode"`
	Country            string `json:"country"`
}

// DocumentInfo represents document verification information
type DocumentInfo struct {
	IDType              string `json:"idType"`
	IDNumber            string `json:"idNumber"`
	IDExpiryDate        string `json:"idExpiryDate,omitempty"`
	IDFrontImageStatus  string `json:"idFrontImageStatus"`
	IDBackImageStatus   string `json:"idBackImageStatus,omitempty"`
	SelfieImageStatus   string `json:"selfieImageStatus"`
	DocumentMatch       bool   `json:"documentMatch"`
	FaceMatch          bool   `json:"faceMatch"`
}

// BusinessInfo represents business information from KYB
type BusinessInfo struct {
	CompanyName         string `json:"companyName"`
	BusinessType        string `json:"businessType"`
	RegistrationNumber  string `json:"registrationNumber"`
	RegistrationCountry string `json:"registrationCountry"`
	BusinessAddress     string `json:"businessAddress"`
	BusinessCity        string `json:"businessCity"`
	BusinessState       string `json:"businessState,omitempty"`
	BusinessPostalCode  string `json:"businessPostalCode"`
	BusinessCountry     string `json:"businessCountry"`
	Industry            string `json:"industry"`
	Website             string `json:"website,omitempty"`
}

// VerificationChecks represents various verification checks performed
type VerificationChecks struct {
	IdentityVerification   *CheckResult `json:"identityVerification,omitempty"`
	DocumentVerification   *CheckResult `json:"documentVerification,omitempty"`
	BiometricVerification  *CheckResult `json:"biometricVerification,omitempty"`
	AddressVerification    *CheckResult `json:"addressVerification,omitempty"`
	WatchlistScreening     *CheckResult `json:"watchlistScreening,omitempty"`
	SanctionsScreening     *CheckResult `json:"sanctionsScreening,omitempty"`
	PEPScreening          *CheckResult `json:"pepScreening,omitempty"`
}

// ComplianceChecks represents business compliance checks
type ComplianceChecks struct {
	BusinessRegistration   *CheckResult `json:"businessRegistration,omitempty"`
	DirectorsScreening     *CheckResult `json:"directorsScreening,omitempty"`
	ShareholdersScreening  *CheckResult `json:"shareholdersScreening,omitempty"`
	UBOVerification       *CheckResult `json:"uboVerification,omitempty"`
	LicenseVerification   *CheckResult `json:"licenseVerification,omitempty"`
	WatchlistScreening    *CheckResult `json:"watchlistScreening,omitempty"`
	SanctionsScreening    *CheckResult `json:"sanctionsScreening,omitempty"`
}

// CheckResult represents the result of a verification check
type CheckResult struct {
	Status    string `json:"status"`    // PASS, FAIL, WARNING, PENDING
	Details   string `json:"details,omitempty"`
	Timestamp string `json:"timestamp"`
	Score     *int   `json:"score,omitempty"` // Risk score if applicable
}

// RiskAssessment represents overall risk assessment
type RiskAssessment struct {
	RiskLevel   string `json:"riskLevel"`   // LOW, MEDIUM, HIGH
	RiskScore   int    `json:"riskScore"`   // Numerical risk score
	Factors     []string `json:"factors,omitempty"` // Risk factors identified
	LastUpdated string `json:"lastUpdated"`
}

// CDD status constants
const (
	CDDStatusPending   = "PENDING"
	CDDStatusApproved  = "APPROVED"
	CDDStatusRejected  = "REJECTED"
	CDDStatusExpired   = "EXPIRED"
	CDDStatusIncomplete = "INCOMPLETE"
)

// Verification check status constants
const (
	CheckStatusPass    = "PASS"
	CheckStatusFail    = "FAIL"
	CheckStatusWarning = "WARNING"
	CheckStatusPending = "PENDING"
)

// Risk level constants
const (
	RiskLevelLow    = "LOW"
	RiskLevelMedium = "MEDIUM"
	RiskLevelHigh   = "HIGH"
)

// FileUploadResponse represents the file upload response
type FileUploadResponse struct {
	BaseResponse
	Data interface{} `json:"data"` // The structure may vary based on API response
}

// Config represents the SDK configuration
type Config struct {
	BaseURL   string
	ClientID  string
	UserAgent string
	Timeout   time.Duration
}

// DefaultConfig returns the default configuration for sandbox environment
func DefaultConfig() *Config {
	return &Config{
		BaseURL:   "https://api-sandbox.interlace.money",
		UserAgent: "interlace-go-sdk/1.0.0",
		Timeout:   30 * time.Second,
	}
}

// Error represents an API error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Interlace API Error - Code: %s, Message: %s", e.Code, e.Message)
}

// ParseError parses an error response from the API
func ParseError(body []byte) *Error {
	var apiError BaseResponse
	if err := json.Unmarshal(body, &apiError); err != nil {
		return &Error{
			Code:    "PARSE_ERROR",
			Message: "Failed to parse error response",
		}
	}

	return &Error{
		Code:    apiError.Code,
		Message: apiError.Message,
	}
}

// Card Management Types

// Card represents a card entity from the list cards API
type Card struct {
	ID                string  `json:"id"`
	AccountID         string  `json:"accountId"`
	CardType          string  `json:"cardType"`
	CardStatus        string  `json:"cardStatus"`
	PersonCardStatus  string  `json:"personCardStatus,omitempty"`
	CardBIN           string  `json:"cardBin"`
	Last4Digits       string  `json:"last4Digits"`
	ExpiryMonth       string  `json:"expiryMonth"`
	ExpiryYear        string  `json:"expiryYear"`
	Currency          string  `json:"currency"`
	Balance           *float64 `json:"balance,omitempty"`
	CreditLimit       *float64 `json:"creditLimit,omitempty"`
	AvailableLimit    *float64 `json:"availableLimit,omitempty"`
	IsPhysical        bool    `json:"isPhysical"`
	IsActive          bool    `json:"isActive"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
	CardholderName    string  `json:"cardholderName,omitempty"`
}

// CardListResponse represents the response from list cards API
type CardListResponse struct {
	Cards      []Card `json:"cards"`
	TotalCount int    `json:"totalCount"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	HasMore    bool   `json:"hasMore"`
}

// CardPrivateInfo represents sensitive card information from get card details API
// Sensitive fields (cardNumber, cvv) are encrypted using AES with clientSecret
type CardPrivateInfo struct {
	ID               string `json:"id"`
	CardNumber       string `json:"cardNumber"`       // Encrypted with AES
	CVV              string `json:"cvv"`              // Encrypted with AES
	ExpiryMonth      string `json:"expiryMonth"`
	ExpiryYear       string `json:"expiryYear"`
	CardholderName   string `json:"cardholderName"`
	CardBIN          string `json:"cardBin"`          // Plain text
	Last4Digits      string `json:"last4Digits"`     // Plain text
	CardStatus       string `json:"cardStatus"`
	IsActive         bool   `json:"isActive"`
}

// CardRemoveResponse represents the response from remove card API
type CardRemoveResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	CardID    string `json:"cardId"`
	RemovedAt string `json:"removedAt"`
}

// Crypto Wallet Management Types

// Wallet represents a crypto wallet
type Wallet struct {
	ID          string `json:"id"`
	AccountID   string `json:"accountId"`
	Nickname    string `json:"nickname,omitempty"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// WalletListResponse represents the response from list wallets API
type WalletListResponse struct {
	Wallets    []Wallet `json:"wallets"`
	TotalCount int      `json:"totalCount"`
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	HasMore    bool     `json:"hasMore"`
}

// WalletAddress represents a blockchain address for a wallet
type WalletAddress struct {
	ID          string `json:"id"`
	WalletID    string `json:"walletId"`
	Currency    string `json:"currency"`
	Chain       string `json:"chain"`
	Address     string `json:"address"`
	Tag         string `json:"tag,omitempty"`
	CreatedAt   string `json:"createdAt"`
}

// Blockchain Transfer Types

// BlockchainTransfer represents a blockchain transfer
type BlockchainTransfer struct {
	ID              string  `json:"id"`
	WalletID        string  `json:"walletId"`
	Currency        string  `json:"currency"`
	Chain           string  `json:"chain"`
	Amount          string  `json:"amount"`
	Fee             string  `json:"fee,omitempty"`
	ToAddress       string  `json:"toAddress"`
	Tag             string  `json:"tag,omitempty"`
	Status          string  `json:"status"`
	TxHash          string  `json:"txHash,omitempty"`
	Confirmations   int     `json:"confirmations,omitempty"`
	RequiredConfirms int    `json:"requiredConfirms,omitempty"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// TransferListResponse represents the response from list transfers API
type TransferListResponse struct {
	Transfers  []BlockchainTransfer `json:"transfers"`
	TotalCount int                  `json:"totalCount"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	HasMore    bool                 `json:"hasMore"`
}

// TransferKYT represents KYT (Know Your Transaction) information
type TransferKYT struct {
	TransferID   string `json:"transferId"`
	RiskScore    int    `json:"riskScore"`
	RiskLevel    string `json:"riskLevel"`
	Alerts       []KYTAlert `json:"alerts,omitempty"`
	CheckedAt    string `json:"checkedAt"`
}

// KYTAlert represents a KYT alert
type KYTAlert struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

// FeeAndQuota represents transfer fee and cross-chain quota information
type FeeAndQuota struct {
	Fee              string `json:"fee"`
	FeeCurrency      string `json:"feeCurrency"`
	MinAmount        string `json:"minAmount"`
	MaxAmount        string `json:"maxAmount"`
	AvailableQuota   string `json:"availableQuota"`
	EstimatedArrival string `json:"estimatedArrival,omitempty"`
}

// Payment and Refund Types

// Payment represents a payment order
type Payment struct {
	ID              string `json:"id"`
	MerchantTradeNo string `json:"merchantTradeNo"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	Status          string `json:"status"`
	PaymentMethod   string `json:"paymentMethod,omitempty"`
	Description     string `json:"description,omitempty"`
	CheckoutURL     string `json:"checkoutUrl,omitempty"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
	PaidAt          string `json:"paidAt,omitempty"`
}

// Refund represents a refund order
type Refund struct {
	ID              string `json:"id"`
	PaymentID       string `json:"paymentId"`
	MerchantTradeNo string `json:"merchantTradeNo"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	Status          string `json:"status"`
	Reason          string `json:"reason,omitempty"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
	ProcessedAt     string `json:"processedAt,omitempty"`
}

// SearchResult represents search results for payments and refunds
type SearchResult struct {
	Payments []Payment `json:"payments,omitempty"`
	Refunds  []Refund  `json:"refunds,omitempty"`
}

// Cardholder represents a cardholder
type Cardholder struct {
	ID               string  `json:"id"`
	AccountID        string  `json:"accountId"`
	FirstName        string  `json:"firstName"`
	LastName         string  `json:"lastName"`
	Email            string  `json:"email"`
	PhoneNumber      string  `json:"phoneNumber"`
	PhoneCountryCode string  `json:"phoneCountryCode"`
	DateOfBirth      string  `json:"dateOfBirth"`
	Nationality      string  `json:"nationality"`
	Gender           string  `json:"gender,omitempty"`
	Occupation       string  `json:"occupation,omitempty"`
	Status           string  `json:"status"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

// CardholderListResponse represents the response for listing cardholders
type CardholderListResponse struct {
	List  []Cardholder `json:"list"`
	Total string       `json:"total"`
}

// CardBin represents a card BIN (Bank Identification Number)
type CardBin struct {
	ID          string `json:"id"`
	Bin         string `json:"bin"`
	Brand       string `json:"brand"`       // VISA, MASTERCARD, etc.
	Type        string `json:"type"`        // VIRTUAL, PHYSICAL
	Currency    string `json:"currency"`
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CardBinListResponse represents the response for listing card BINs
type CardBinListResponse struct {
	List  []CardBin `json:"list"`
	Total string    `json:"total,omitempty"`
}

// ConsumptionScenario represents a card transaction scenario
type ConsumptionScenario struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
	Status      string `json:"status,omitempty"`
}

// ConsumptionScenarioListResponse represents the response for listing consumption scenarios
type ConsumptionScenarioListResponse struct {
	List  []ConsumptionScenario `json:"list"`
	Total string                `json:"total,omitempty"`
}
