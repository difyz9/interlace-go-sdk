package interlace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// HTTPClient is a wrapper around http.Client that handles common operations
type HTTPClient struct {
	config      *Config
	httpClient  *http.Client
	accessToken string
}

// NewHTTPClient creates a new HTTP client wrapper
func NewHTTPClient(config *Config, accessToken string) *HTTPClient {
	if config == nil {
		config = DefaultConfig()
	}

	return &HTTPClient{
		config:      config,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// SetAccessToken updates the access token
func (c *HTTPClient) SetAccessToken(accessToken string) {
	c.accessToken = accessToken
}

// GetAccessToken returns the current access token
func (c *HTTPClient) GetAccessToken() string {
	return c.accessToken
}

// RequestOptions holds options for HTTP requests
type RequestOptions struct {
	Method       string
	Endpoint     string
	Body         interface{}
	QueryParams  url.Values
	Headers      map[string]string
	RequireAuth  bool
	ContentType  string
}

// DoRequest performs an HTTP request with common handling
func (c *HTTPClient) DoRequest(ctx context.Context, opts *RequestOptions, result interface{}) error {
	// Build URL
	fullURL := fmt.Sprintf("%s%s", c.config.BaseURL, opts.Endpoint)
	if opts.QueryParams != nil && len(opts.QueryParams) > 0 {
		fullURL = fmt.Sprintf("%s?%s", fullURL, opts.QueryParams.Encode())
	}

	// Prepare request body
	var bodyReader io.Reader
	if opts.Body != nil {
		if bodyBytes, ok := opts.Body.([]byte); ok {
			bodyReader = bytes.NewReader(bodyBytes)
		} else if reader, ok := opts.Body.(io.Reader); ok {
			// Body is already a reader (e.g., for file uploads)
			bodyReader = reader
		} else {
			// JSON marshal the body
			jsonData, err := json.Marshal(opts.Body)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %w", err)
			}
			bodyReader = bytes.NewBuffer(jsonData)
		}
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, opts.Method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.config.UserAgent)

	// Set content type
	if opts.ContentType != "" {
		req.Header.Set("Content-Type", opts.ContentType)
	} else if opts.Body != nil && opts.Method != "GET" {
		// Default to JSON for non-GET requests with body
		if _, ok := opts.Body.(io.Reader); !ok {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	// Set authentication header if required
	if opts.RequireAuth && c.accessToken != "" {
		req.Header.Set("x-access-token", c.accessToken)
	}

	// Set custom headers
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		apiError := ParseError(respBody)
		return apiError
	}

	// Parse JSON response if result is provided
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// DoJSONRequest is a convenience method for JSON requests
func (c *HTTPClient) DoJSONRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	opts := &RequestOptions{
		Method:      method,
		Endpoint:    endpoint,
		Body:        body,
		RequireAuth: true,
		ContentType: "application/json",
	}

	return c.DoRequest(ctx, opts, result)
}

// DoGetRequest is a convenience method for GET requests
func (c *HTTPClient) DoGetRequest(ctx context.Context, endpoint string, queryParams url.Values, result interface{}) error {
	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    endpoint,
		QueryParams: queryParams,
		RequireAuth: true,
	}

	return c.DoRequest(ctx, opts, result)
}

// DoPostRequest is a convenience method for POST requests
func (c *HTTPClient) DoPostRequest(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.DoJSONRequest(ctx, "POST", endpoint, body, result)
}

// DoGetRequestNoAuth is a convenience method for GET requests without authentication
func (c *HTTPClient) DoGetRequestNoAuth(ctx context.Context, endpoint string, queryParams url.Values, result interface{}) error {
	opts := &RequestOptions{
		Method:      "GET",
		Endpoint:    endpoint,
		QueryParams: queryParams,
		RequireAuth: false,
	}

	return c.DoRequest(ctx, opts, result)
}

// DoPostRequestNoAuth is a convenience method for POST requests without authentication
func (c *HTTPClient) DoPostRequestNoAuth(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    endpoint,
		Body:        body,
		RequireAuth: false,
		ContentType: "application/json",
	}

	return c.DoRequest(ctx, opts, result)
}