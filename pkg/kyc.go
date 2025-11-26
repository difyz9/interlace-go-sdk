package interlace

import (
	"context"
	"fmt"
)

// KYCClient handles KYC (Know Your Customer) operations
type KYCClient struct {
	httpClient *HTTPClient
}

// NewKYCClient creates a new KYC client
func NewKYCClient(httpClient *HTTPClient) *KYCClient {
	return &KYCClient{
		httpClient: httpClient,
	}
}



// SubmitKYC submits KYC information for the specified account
func (c *KYCClient) SubmitKYC(ctx context.Context, accountID string, req *KYCSubmitRequest) (*KYCSubmitData, error) {
	endpoint := fmt.Sprintf("/open-api/v3/accounts/%s/kyc", accountID)
	
	var kycResp KYCSubmitResponse
	err := c.httpClient.DoPostRequest(ctx, endpoint, req, &kycResp)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if kycResp.GetCode() != "000000" {
		return nil, &Error{
			Code:    kycResp.GetCode(),
			Message: kycResp.Message,
		}
	}

	return &kycResp.Data, nil
}

// GetKYCStatus retrieves the KYC status for the specified account
func (c *KYCClient) GetKYCStatus(ctx context.Context, accountID string) (*KYCStatusData, error) {
	endpoint := fmt.Sprintf("/open-api/v3/accounts/%s/kyc", accountID)
	
	var statusResp KYCStatusResponse
	err := c.httpClient.DoGetRequest(ctx, endpoint, nil, &statusResp)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if statusResp.GetCode() != "000000" {
		return nil, &Error{
			Code:    statusResp.GetCode(),
			Message: statusResp.Message,
		}
	}

	return &statusResp.Data, nil
}

// IsKYCApproved checks if the account's KYC is approved
func (c *KYCClient) IsKYCApproved(ctx context.Context, accountID string) (bool, error) {
	status, err := c.GetKYCStatus(ctx, accountID)
	if err != nil {
		return false, err
	}

	return status.Status == KYCStatusApproved, nil
}

// IsKYCPending checks if the account's KYC is pending
func (c *KYCClient) IsKYCPending(ctx context.Context, accountID string) (bool, error) {
	status, err := c.GetKYCStatus(ctx, accountID)
	if err != nil {
		return false, err
	}

	return status.Status == KYCStatusPending, nil
}

// IsKYCRejected checks if the account's KYC is rejected
func (c *KYCClient) IsKYCRejected(ctx context.Context, accountID string) (bool, error) {
	status, err := c.GetKYCStatus(ctx, accountID)
	if err != nil {
		return false, err
	}

	return status.Status == KYCStatusRejected, nil
}

// WaitForKYCApproval waits for KYC approval with polling
// This is a convenience method for testing/development
func (c *KYCClient) WaitForKYCApproval(ctx context.Context, accountID string, maxAttempts int) (*KYCStatusData, error) {
	for i := 0; i < maxAttempts; i++ {
		status, err := c.GetKYCStatus(ctx, accountID)
		if err != nil {
			return nil, err
		}

		switch status.Status {
		case KYCStatusApproved:
			return status, nil
		case KYCStatusRejected:
			return status, fmt.Errorf("KYC was rejected: %s", status.RejectionReason)
		case KYCStatusExpired:
			return status, fmt.Errorf("KYC has expired")
		case KYCStatusPending:
			// Continue polling
			continue
		default:
			return status, fmt.Errorf("unknown KYC status: %s", status.Status)
		}
	}

	return nil, fmt.Errorf("KYC approval timeout after %d attempts", maxAttempts)
}

// GetCDDDetail retrieves comprehensive CDD (Customer Due Diligence) details including KYC and KYB verification results
func (c *KYCClient) GetCDDDetail(ctx context.Context, accountID string) (*CDDDetailData, error) {
	endpoint := fmt.Sprintf("/open-api/v3/accounts/cdd/detail/%s", accountID)
	
	var cddResp CDDDetailResponse
	err := c.httpClient.DoGetRequest(ctx, endpoint, nil, &cddResp)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if cddResp.GetCode() != "000000" {
		return nil, &Error{
			Code:    cddResp.GetCode(),
			Message: cddResp.Message,
		}
	}

	return &cddResp.Data, nil
}

// GetKYCVerificationDetail extracts KYC verification details from CDD data
func (c *KYCClient) GetKYCVerificationDetail(ctx context.Context, accountID string) (*KYCVerificationDetail, error) {
	cddDetail, err := c.GetCDDDetail(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if cddDetail.KYCVerification == nil {
		return nil, fmt.Errorf("no KYC verification data found for account %s", accountID)
	}

	return cddDetail.KYCVerification, nil
}

// GetKYBVerificationDetail extracts KYB verification details from CDD data
func (c *KYCClient) GetKYBVerificationDetail(ctx context.Context, accountID string) (*KYBVerificationDetail, error) {
	cddDetail, err := c.GetCDDDetail(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if cddDetail.KYBVerification == nil {
		return nil, fmt.Errorf("no KYB verification data found for account %s", accountID)
	}

	return cddDetail.KYBVerification, nil
}

// GetRiskAssessment extracts risk assessment from CDD data
func (c *KYCClient) GetRiskAssessment(ctx context.Context, accountID string) (*RiskAssessment, error) {
	cddDetail, err := c.GetCDDDetail(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Try KYC risk assessment first
	if cddDetail.KYCVerification != nil && cddDetail.KYCVerification.RiskAssessment != nil {
		return cddDetail.KYCVerification.RiskAssessment, nil
	}

	// Try KYB risk assessment
	if cddDetail.KYBVerification != nil && cddDetail.KYBVerification.RiskAssessment != nil {
		return cddDetail.KYBVerification.RiskAssessment, nil
	}

	return nil, fmt.Errorf("no risk assessment found for account %s", accountID)
}

// IsHighRisk checks if the account is classified as high risk
func (c *KYCClient) IsHighRisk(ctx context.Context, accountID string) (bool, error) {
	riskAssessment, err := c.GetRiskAssessment(ctx, accountID)
	if err != nil {
		return false, err
	}

	return riskAssessment.RiskLevel == RiskLevelHigh, nil
}

// GetVerificationChecks extracts verification checks from KYC data
func (c *KYCClient) GetVerificationChecks(ctx context.Context, accountID string) (*VerificationChecks, error) {
	kycDetail, err := c.GetKYCVerificationDetail(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if kycDetail.VerificationChecks == nil {
		return nil, fmt.Errorf("no verification checks found for account %s", accountID)
	}

	return kycDetail.VerificationChecks, nil
}

// GetComplianceChecks extracts compliance checks from KYB data
func (c *KYCClient) GetComplianceChecks(ctx context.Context, accountID string) (*ComplianceChecks, error) {
	kybDetail, err := c.GetKYBVerificationDetail(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if kybDetail.ComplianceChecks == nil {
		return nil, fmt.Errorf("no compliance checks found for account %s", accountID)
	}

	return kybDetail.ComplianceChecks, nil
}

// HasPassedAllChecks verifies if all required verification checks have passed
func (c *KYCClient) HasPassedAllChecks(ctx context.Context, accountID string) (bool, []string, error) {
	var failedChecks []string

	// Check KYC verification checks
	kycDetail, err := c.GetKYCVerificationDetail(ctx, accountID)
	if err == nil && kycDetail.VerificationChecks != nil {
		checks := kycDetail.VerificationChecks
		
		if checks.IdentityVerification != nil && checks.IdentityVerification.Status != CheckStatusPass {
			failedChecks = append(failedChecks, "Identity Verification")
		}
		if checks.DocumentVerification != nil && checks.DocumentVerification.Status != CheckStatusPass {
			failedChecks = append(failedChecks, "Document Verification")
		}
		if checks.BiometricVerification != nil && checks.BiometricVerification.Status != CheckStatusPass {
			failedChecks = append(failedChecks, "Biometric Verification")
		}
		if checks.WatchlistScreening != nil && checks.WatchlistScreening.Status == CheckStatusFail {
			failedChecks = append(failedChecks, "Watchlist Screening")
		}
		if checks.SanctionsScreening != nil && checks.SanctionsScreening.Status == CheckStatusFail {
			failedChecks = append(failedChecks, "Sanctions Screening")
		}
		if checks.PEPScreening != nil && checks.PEPScreening.Status == CheckStatusFail {
			failedChecks = append(failedChecks, "PEP Screening")
		}
	}

	// Check KYB compliance checks
	kybDetail, err := c.GetKYBVerificationDetail(ctx, accountID)
	if err == nil && kybDetail.ComplianceChecks != nil {
		checks := kybDetail.ComplianceChecks
		
		if checks.BusinessRegistration != nil && checks.BusinessRegistration.Status != CheckStatusPass {
			failedChecks = append(failedChecks, "Business Registration")
		}
		if checks.DirectorsScreening != nil && checks.DirectorsScreening.Status == CheckStatusFail {
			failedChecks = append(failedChecks, "Directors Screening")
		}
		if checks.UBOVerification != nil && checks.UBOVerification.Status != CheckStatusPass {
			failedChecks = append(failedChecks, "UBO Verification")
		}
	}

	return len(failedChecks) == 0, failedChecks, nil
}

// KYCBuilder helps build KYC requests with a fluent interface
type KYCBuilder struct {
	request *KYCSubmitRequest
}

// NewKYCBuilder creates a new KYC request builder
func NewKYCBuilder() *KYCBuilder {
	return &KYCBuilder{
		request: &KYCSubmitRequest{},
	}
}

// SetPersonalInfo sets basic personal information
func (b *KYCBuilder) SetPersonalInfo(firstName, lastName, dateOfBirth, gender string) *KYCBuilder {
	b.request.FirstName = firstName
	b.request.LastName = lastName
	b.request.DateOfBirth = dateOfBirth
	b.request.Gender = gender
	return b
}

// SetMiddleName sets middle name (optional)
func (b *KYCBuilder) SetMiddleName(middleName string) *KYCBuilder {
	b.request.MiddleName = middleName
	return b
}

// SetNationality sets nationality and country of residence
func (b *KYCBuilder) SetNationality(nationality, countryOfResidence string) *KYCBuilder {
	b.request.Nationality = nationality
	b.request.CountryOfResidence = countryOfResidence
	return b
}

// SetAddress sets address information
func (b *KYCBuilder) SetAddress(address, city, postalCode, country string) *KYCBuilder {
	b.request.Address = address
	b.request.City = city
	b.request.PostalCode = postalCode
	b.request.Country = country
	return b
}

// SetState sets state (optional, for countries that have states)
func (b *KYCBuilder) SetState(state string) *KYCBuilder {
	b.request.State = state
	return b
}

// SetIDInfo sets ID document information
func (b *KYCBuilder) SetIDInfo(idType, idNumber string) *KYCBuilder {
	b.request.IDType = idType
	b.request.IDNumber = idNumber
	return b
}

// SetIDExpiryDate sets ID expiry date (optional)
func (b *KYCBuilder) SetIDExpiryDate(expiryDate string) *KYCBuilder {
	b.request.IDExpiryDate = expiryDate
	return b
}

// SetOccupationInfo sets occupation and income information
func (b *KYCBuilder) SetOccupationInfo(occupation, sourceOfIncome string) *KYCBuilder {
	b.request.Occupation = occupation
	b.request.SourceOfIncome = sourceOfIncome
	return b
}

// SetAnnualIncome sets annual income (optional)
func (b *KYCBuilder) SetAnnualIncome(annualIncome string) *KYCBuilder {
	b.request.AnnualIncome = annualIncome
	return b
}

// SetAccountPurpose sets purpose of account and expected transaction volume
func (b *KYCBuilder) SetAccountPurpose(purposeOfAccount string) *KYCBuilder {
	b.request.PurposeOfAccount = purposeOfAccount
	return b
}

// SetExpectedTxnVolume sets expected transaction volume (optional)
func (b *KYCBuilder) SetExpectedTxnVolume(expectedVolume string) *KYCBuilder {
	b.request.ExpectedTxnVolume = expectedVolume
	return b
}

// SetDocumentFiles sets the file IDs for ID documents and selfie
func (b *KYCBuilder) SetDocumentFiles(idFrontFileID, selfieFileID string) *KYCBuilder {
	b.request.IDFrontImageFileID = idFrontFileID
	b.request.SelfieImageFileID = selfieFileID
	return b
}

// SetIDBackFile sets the back side of ID document (optional)
func (b *KYCBuilder) SetIDBackFile(idBackFileID string) *KYCBuilder {
	b.request.IDBackImageFileID = idBackFileID
	return b
}

// Build returns the completed KYC request
func (b *KYCBuilder) Build() *KYCSubmitRequest {
	return b.request
}

// Validate checks if the KYC request has all required fields
func (b *KYCBuilder) Validate() error {
	req := b.request

	if req.FirstName == "" {
		return fmt.Errorf("firstName is required")
	}
	if req.LastName == "" {
		return fmt.Errorf("lastName is required")
	}
	if req.DateOfBirth == "" {
		return fmt.Errorf("dateOfBirth is required")
	}
	if req.Gender == "" {
		return fmt.Errorf("gender is required")
	}
	if req.Nationality == "" {
		return fmt.Errorf("nationality is required")
	}
	if req.CountryOfResidence == "" {
		return fmt.Errorf("countryOfResidence is required")
	}
	if req.Address == "" {
		return fmt.Errorf("address is required")
	}
	if req.City == "" {
		return fmt.Errorf("city is required")
	}
	if req.PostalCode == "" {
		return fmt.Errorf("postalCode is required")
	}
	if req.Country == "" {
		return fmt.Errorf("country is required")
	}
	if req.IDType == "" {
		return fmt.Errorf("idType is required")
	}
	if req.IDNumber == "" {
		return fmt.Errorf("idNumber is required")
	}
	if req.Occupation == "" {
		return fmt.Errorf("occupation is required")
	}
	if req.SourceOfIncome == "" {
		return fmt.Errorf("sourceOfIncome is required")
	}
	if req.PurposeOfAccount == "" {
		return fmt.Errorf("purposeOfAccount is required")
	}
	if req.IDFrontImageFileID == "" {
		return fmt.Errorf("idFrontImageFileId is required")
	}
	if req.SelfieImageFileID == "" {
		return fmt.Errorf("selfieImageFileId is required")
	}

	return nil
}