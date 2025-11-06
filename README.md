# Interlace Go SDK

A comprehensive Golang SDK for the Interlace Money API, providing easy-to-use methods for OAuth authentication, account management, and file upload functionality.

## Features

- **OAuth 2.0 Authentication**: Complete OAuth flow with authorization and token exchange
- **Account Management**: Register accounts, list accounts, and retrieve specific account details
- **Card Management**: List, filter, and manage cards with comprehensive querying options
- **KYC/CDD Operations**: Submit KYC information, check verification status, and retrieve CDD details
- **File Upload**: Upload single or multiple files with multipart form data
- **HTTP Client Abstraction**: Unified HTTP client that handles headers, authentication, and error processing
- **Security**: AES encryption for sensitive card data using client secret
- **Error Handling**: Comprehensive error handling with detailed error messages
- **Configurable**: Support for both sandbox and production environments
- **Type-Safe**: Full type definitions for all API requests and responses
- **Clean Architecture**: Minimal code duplication with centralized HTTP operations

## Installation

```bash
go get github.com/difyz9/interlace-go-sdk
```

## ⚙️ Configuration

Before using the SDK, you need to configure your API credentials. See the [Configuration Guide](./CONFIGURATION.md) for detailed instructions on:

- How to obtain API credentials
- Setting up environment variables
- Best practices for managing secrets
- Troubleshooting common configuration issues

**Quick Setup**:
1. Copy `.env.example` to `.env`
2. Add your Client ID: `INTERLACE_CLIENT_ID=your-client-id-here`
3. Load environment variables in your code

For security reasons, **never commit your `.env` file or hardcode credentials** in your source code.

## Architecture

The SDK uses a layered architecture with HTTP client abstraction:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Client API    │────│  Business Logic  │────│  HTTP Client    │
│  (OAuth, etc.)  │    │   (Types, etc.)  │    │   (Requests)    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

**Benefits of HTTP Client Abstraction:**
- **Unified Headers**: Automatic handling of Accept, Content-Type, User-Agent, and authentication headers
- **Centralized Error Handling**: Consistent error processing across all API calls  
- **Reduced Code Duplication**: No need to repeat HTTP setup in every method
- **Easy Maintenance**: Changes to HTTP logic only need to be made in one place
- **Type Safety**: Structured request options and response handling

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
    clientID := "your-client-id"
    
    // Quick setup - handles authentication automatically
    client, tokenData, err := interlace.QuickSetup(clientID, nil)
    if err != nil {
        log.Fatalf("Setup failed: %v", err)
    }
    
    fmt.Printf("Access Token: %s\n", tokenData.AccessToken)
    fmt.Printf("Client authenticated: %v\n", client.IsAuthenticated())
}
```

### Manual Setup

```go
package main

import (
    "context"
    
    interlace "github.com/difyz9/interlace-go-sdk/pkg"
)

func main() {
    // Create client with custom configuration
    config := interlace.DefaultConfig() // or interlace.ProductionConfig()
    config.ClientID = "your-client-id"
    
    client := interlace.NewClient(config)
    
    // Perform authentication
    ctx := context.Background()
    tokenData, err := client.Authenticate(ctx, "your-client-id")
    if err != nil {
        // handle error
    }
    
    // Now you can use the client for API calls
}
```

## API Reference

### OAuth Operations

#### Authorize and Get Token (Combined)

```go
tokenData, err := client.Authenticate(ctx, clientID)
```

#### Manual OAuth Flow

```go
// Step 1: Get authorization code
authData, err := client.OAuth.Authorize(ctx, clientID)

// Step 2: Exchange code for token
tokenData, err := client.OAuth.GetAccessToken(ctx, authData.Code, clientID)

// Step 3: Set token for future requests
client.SetAccessToken(tokenData.AccessToken)
```

### Account Management

#### Register Account

**Method 1: Using struct (recommended for complex scenarios)**
```go
registerReq := &interlace.AccountRegisterRequest{
    PhoneNumber:      "15900000042",
    Email:            "user@example.com",
    Name:             "John Doe",
    PhoneCountryCode: "86",
}

account, err := client.Account.Register(ctx, registerReq)
```

**Method 2: Using direct parameters (convenient)**
```go
account, err := client.Account.RegisterWithDetails(ctx, "86", "15900000031", "15900000031@qq.com", "golang_test")
```

**Method 3: Quick test account (for development)**
```go
// Creates account with preset test data: phone=15900000031, email=15900000031@qq.com, name=golang_test
account, err := client.Account.RegisterGolangTest(ctx)
```

#### List Accounts

**Basic pagination**
```go
opts := &interlace.AccountListOptions{
    Limit: 10,
    Page:  1,
}
accountList, err := client.Account.List(ctx, opts)
```

**Filter by account ID**
```go
opts := &interlace.AccountListOptions{
    AccountID: "specific-account-id",
    Limit:     10,
    Page:      1,
}
accountList, err := client.Account.List(ctx, opts)
```

**Filter by status and type**
```go
opts := &interlace.AccountListOptions{
    Status: interlace.AccountStatusActive, // ACTIVE, INACTIVE, PENDING, SUSPENDED
    Type:   interlace.AccountTypePersonal, // 1=Personal, 2=Business, 3=Child
    Limit:  50,
    Page:   1,
}
accountList, err := client.Account.List(ctx, opts)
```

**Convenience methods**
```go
// Get all accounts (auto-pagination)
allAccounts, err := client.Account.ListAll(ctx)

// Get active accounts only
activeAccounts, err := client.Account.ListActiveAccounts(ctx)

// Get accounts by status
inactiveAccounts, err := client.Account.ListByStatus(ctx, "INACTIVE")

// Get accounts by type
personalAccounts, err := client.Account.ListByType(ctx, interlace.AccountTypePersonal)

// Get total account count
count, err := client.Account.Count(ctx)

// Get specific page
pageData, err := client.Account.GetAccountsByPage(ctx, 2, 20) // page 2, 20 per page
```

#### Get Specific Account

```go
account, err := client.Account.Get(ctx, "account-id")
```

### Card Management Operations

#### List Cards

**Basic listing with pagination**
```go
opts := &interlace.CardListOptions{
    Limit: 10,
    Page:  1,
}
cardList, err := client.Card.ListCards(ctx, opts)
```

**Filter by account ID**
```go
opts := &interlace.CardListOptions{
    AccountID: "specific-account-id",
    Limit:     10,
    Page:      1,
}
cardList, err := client.Card.ListCards(ctx, opts)
```

**Filter by status and type**
```go
opts := &interlace.CardListOptions{
    CardStatus: interlace.CardStatusActive,   // ACTIVE, INACTIVE, BLOCKED, etc.
    CardType:   interlace.CardTypeVirtual,    // VIRTUAL, PHYSICAL, PREPAID, CREDIT
    IsActive:   &[]bool{true}[0],             // Filter by active status
    Limit:      50,
    Page:       1,
}
cardList, err := client.Card.ListCards(ctx, opts)
```

**Convenience methods**
```go
// Get all cards (auto-pagination)
allCards, err := client.Card.ListAllCards(ctx)

// Get cards for specific account
accountCards, err := client.Card.ListCardsByAccount(ctx, "account-id")

// Get active cards only
activeCards, err := client.Card.ListActiveCards(ctx)

// Get cards by status
blockedCards, err := client.Card.ListCardsByStatus(ctx, interlace.CardStatusBlocked)

// Get cards by type
virtualCards, err := client.Card.ListCardsByType(ctx, interlace.CardTypeVirtual)

// Get total card count
count, err := client.Card.CountCards(ctx)

// Get specific page
pageData, err := client.Card.GetCardsByPage(ctx, 2, 20) // page 2, 20 per page
```

#### Get Card Private Information

```go
// Get sensitive card information (encrypted)
cardInfo, err := client.Card.GetCardPrivateInfo(ctx, "card-id")

// Note: Sensitive data like card number and CVV will be encrypted using AES
// with clientSecret. Non-sensitive data like BIN and last 4 digits remain in plaintext
```

#### Remove Card

```go
// Delete a card (WARNING: Cannot be recovered)
removeResp, err := client.Card.RemoveCard(ctx, "card-id")

if err == nil {
    fmt.Printf("Card removed at: %s\n", removeResp.RemovedAt)
    
    // Check if balance was transferred
    if removeResp.BalanceTransfer != nil {
        fmt.Printf("Balance %.2f %s transferred to %s\n",
            removeResp.BalanceTransfer.Amount,
            removeResp.BalanceTransfer.Currency,
            removeResp.BalanceTransfer.TransferredTo)
    }
}
```

#### Card Status and Type Constants

**Card Statuses**
- `interlace.CardStatusActive` - Active and usable
- `interlace.CardStatusInactive` - Inactive
- `interlace.CardStatusBlocked` - Blocked/frozen  
- `interlace.CardStatusExpired` - Expired
- `interlace.CardStatusPending` - Pending activation
- `interlace.CardStatusCancelled` - Cancelled

**Card Types**
- `interlace.CardTypeVirtual` - Virtual card
- `interlace.CardTypePhysical` - Physical card
- `interlace.CardTypePrepaid` - Prepaid card
- `interlace.CardTypeCredit` - Credit card

#### Convenience Methods

```go
// Check if user has any cards
hasCards, err := client.Card.HasCards(ctx)

// Check if specific card is active
isActive, err := client.Card.IsCardActive(ctx, "card-id")

// Get card status
status, err := client.Card.GetCardStatus(ctx, "card-id")

// Count cards by account
accountCardCount, err := client.Card.CountCardsByAccount(ctx, "account-id")

// Count active cards
activeCount, err := client.Card.CountActiveCards(ctx)
```

### KYC (Know Your Customer) Operations

#### Submit KYC Information

**Using builder pattern (recommended)**
```go
// Build KYC request using fluent interface
kycRequest := interlace.NewKYCBuilder().
    SetPersonalInfo("John", "Doe", "1990-01-15", "M").
    SetNationality("US", "US").
    SetAddress("123 Main St", "New York", "10001", "US").
    SetIDInfo(interlace.IDTypePassport, "A12345678").
    SetIDExpiryDate("2030-01-15").
    SetOccupationInfo("01", "Employment").
    SetAccountPurpose("Personal Banking").
    SetDocumentFiles("id-front-file-id", "selfie-file-id").
    Build()

// Validate the request
if err := kycBuilder.Validate(); err != nil {
    // handle validation error
}

// Submit KYC
kycResult, err := client.KYC.SubmitKYC(ctx, accountID, kycRequest)
```

**Using struct directly**
```go
kycRequest := &interlace.KYCSubmitRequest{
    FirstName:          "John",
    LastName:           "Doe",
    DateOfBirth:        "1990-01-15",
    Gender:             "M",
    Nationality:        "US",
    CountryOfResidence: "US",
    Address:            "123 Main St",
    City:               "New York",
    PostalCode:         "10001",
    Country:            "US",
    IDType:             interlace.IDTypePassport,
    IDNumber:           "A12345678",
    Occupation:         "01", // See occupation codes documentation
    SourceOfIncome:     "Employment",
    PurposeOfAccount:   "Personal Banking",
    IDFrontImageFileID: "id-front-file-id", // From file upload API
    SelfieImageFileID:  "selfie-file-id",   // From file upload API
}

kycResult, err := client.KYC.SubmitKYC(ctx, accountID, kycRequest)
```

#### Check KYC Status

**Get detailed status**
```go
kycStatus, err := client.KYC.GetKYCStatus(ctx, accountID)
fmt.Printf("Status: %s\n", kycStatus.Status)
fmt.Printf("Application ID: %s\n", kycStatus.KYCApplicationID)
```

**Quick status checks**
```go
// Check if KYC is approved
isApproved, err := client.KYC.IsKYCApproved(ctx, accountID)

// Check if KYC is pending
isPending, err := client.KYC.IsKYCPending(ctx, accountID)

// Check if KYC is rejected
isRejected, err := client.KYC.IsKYCRejected(ctx, accountID)
```

**Wait for approval (polling)**
```go
// Wait for KYC approval with maximum 10 attempts
finalStatus, err := client.KYC.WaitForKYCApproval(ctx, accountID, 10)
```

#### Supported ID Types
- `interlace.IDTypeCNRIC` - China Resident Identity Card
- `interlace.IDTypeHKHKID` - Hong Kong Identity Card  
- `interlace.IDTypePassport` - Passport
- `interlace.IDTypeDLN` - Driver's License Number
- `interlace.IDTypeGovernmentID` - Government-Issued ID Card

#### KYC Status Values
- `interlace.KYCStatusPending` - Under review
- `interlace.KYCStatusApproved` - Approved
- `interlace.KYCStatusRejected` - Rejected
- `interlace.KYCStatusExpired` - Expired

### File Upload

#### Upload Single File

```go
uploadResp, err := client.File.UploadFile(ctx, "/path/to/file.jpg", accountID)
```

#### Upload Multiple Files

```go
filePaths := []string{
    "/path/to/file1.jpg",
    "/path/to/file2.pdf",
}

uploadResp, err := client.File.UploadMultipleFiles(ctx, filePaths, accountID)
```

#### Upload from Reader

```go
file, err := os.Open("file.jpg")
if err != nil {
    // handle error
}
defer file.Close()

uploadResp, err := client.File.UploadFileFromReader(ctx, file, "file.jpg", accountID)
```

## Configuration

### Default Configuration (Sandbox)

```go
config := interlace.DefaultConfig()
// or
config := interlace.SandboxConfig()
```

### Production Configuration

```go
config := interlace.ProductionConfig()
```

### Custom Configuration

```go
config := &interlace.Config{
    BaseURL:   "https://api-sandbox.interlace.money",
    ClientID:  "your-client-id",
    UserAgent: "my-app/1.0.0",
    Timeout:   30 * time.Second,
}
```

## Error Handling

The SDK provides detailed error information:

```go
account, err := client.Account.Register(ctx, registerReq)
if err != nil {
    if apiErr, ok := err.(*interlace.Error); ok {
        fmt.Printf("API Error - Code: %s, Message: %s\n", apiErr.Code, apiErr.Message)
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

## Examples

### Complete Workflow (Replicating curl commands)

```go
func main() {
    // Replace with your actual client ID
    clientID := "your-client-id-here"
    
    // Step 1: OAuth Authorization
    client := interlace.NewClient(nil)
    ctx := context.Background()
    
    authData, err := client.OAuth.Authorize(ctx, clientID)
    if err != nil {
        log.Fatalf("Authorization failed: %v", err)
    }
    
    // Step 2: Get Access Token
    tokenData, err := client.OAuth.GetAccessToken(ctx, authData.Code, clientID)
    if err != nil {
        log.Fatalf("Token exchange failed: %v", err)
    }
    
    client.SetAccessToken(tokenData.AccessToken)
    
    // Step 3: Register Account
    registerReq := &interlace.AccountRegisterRequest{
        PhoneNumber:      "15900000042",
        Email:            "15900000042@qq.com",
        Name:             "saf",
        PhoneCountryCode: "86",
    }
    
    account, err := client.Account.Register(ctx, registerReq)
    if err != nil {
        log.Fatalf("Account registration failed: %v", err)
    }
    
    // Step 4: List Accounts
    listOpts := &interlace.AccountListOptions{
        AccountID: account.ID,
        Limit:     10,
        Page:      1,
    }
    
    accountList, err := client.Account.List(ctx, listOpts)
    if err != nil {
        log.Fatalf("Account listing failed: %v", err)
    }
    
    // Step 5: Upload File
    uploadResp, err := client.File.UploadFile(ctx, "file.jpg", account.ID)
    if err != nil {
        log.Fatalf("File upload failed: %v", err)
    }
}
```

## API Endpoints

The SDK covers the following Interlace API endpoints:

**Authentication**
- `GET /open-api/v3/oauth/authorize` - OAuth authorization
- `POST /open-api/v3/oauth/access-token` - Token exchange

**Account Management**  
- `POST /open-api/v3/accounts/register` - Account registration
- `GET /open-api/v3/accounts` - List accounts

**KYC (Know Your Customer)**
- `POST /open-api/v3/accounts/{accountId}/kyc` - Submit KYC information
- `GET /open-api/v3/accounts/{accountId}/kyc` - Get KYC status

**File Operations**
- `POST /open-api/v3/files/upload` - File upload

**Card Management**
- `GET /open-api/v3/card-list` - List all cards with filtering and pagination
- `GET /open-api/v3/cards/{id}` - Get card private information (encrypted sensitive data)
- `DELETE /open-api/v3/cards/{id}` - Remove/delete card (irreversible)

## Environment Support

- **Sandbox**: `https://api-sandbox.interlace.money` (default)
- **Production**: `https://api.interlace.money`

## License

This project is licensed under the MIT License.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

For issues and questions:
- Create an issue on GitHub
- Check the [Interlace API Documentation](https://developer.interlace.money/)

---

Built with ❤️ for the Interlace Money ecosystem.