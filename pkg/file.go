package interlace

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// FileClient handles file operations
type FileClient struct {
	httpClient *HTTPClient
}

// NewFileClient creates a new file client
func NewFileClient(httpClient *HTTPClient) *FileClient {
	return &FileClient{
		httpClient: httpClient,
	}
}



// UploadFile uploads a file to the Interlace API
func (c *FileClient) UploadFile(ctx context.Context, filePath, accountID string) (*FileUploadResponse, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileName := filepath.Base(filePath)
	return c.UploadFileFromReader(ctx, file, fileName, accountID)
}

// UploadFileFromReader uploads a file from an io.Reader
func (c *FileClient) UploadFileFromReader(ctx context.Context, reader io.Reader, fileName, accountID string) (*FileUploadResponse, error) {
	// Create a buffer to store the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the file field
	fileWriter, err := writer.CreateFormFile("files", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy reader content to the form
	_, err = io.Copy(fileWriter, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Add the accountId field
	err = writer.WriteField("accountId", accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to write accountId field: %w", err)
	}

	// Close the writer to finalize the form
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Use the HTTP client with multipart form data
	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/files/upload",
		Body:        &requestBody,
		RequireAuth: true,
		ContentType: writer.FormDataContentType(),
	}

	var uploadResp FileUploadResponse
	err = c.httpClient.DoRequest(ctx, opts, &uploadResp)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if uploadResp.GetCode() != "000000" {
		return nil, &Error{
			Code:    uploadResp.GetCode(),
			Message: uploadResp.Message,
		}
	}

	return &uploadResp, nil
}

// UploadMultipleFiles uploads multiple files to the Interlace API
func (c *FileClient) UploadMultipleFiles(ctx context.Context, filePaths []string, accountID string) (*FileUploadResponse, error) {
	// Create a buffer to store the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add each file
	for _, filePath := range filePaths {
		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		
		// Create form file field
		fileName := filepath.Base(filePath)
		fileWriter, err := writer.CreateFormFile("files", fileName)
		if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to create form file for %s: %w", fileName, err)
		}

		// Copy file content
		_, err = io.Copy(fileWriter, file)
		file.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to copy file content for %s: %w", fileName, err)
		}
	}

	// Add the accountId field
	err := writer.WriteField("accountId", accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to write accountId field: %w", err)
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Use the HTTP client with multipart form data
	opts := &RequestOptions{
		Method:      "POST",
		Endpoint:    "/open-api/v3/files/upload",
		Body:        &requestBody,
		RequireAuth: true,
		ContentType: writer.FormDataContentType(),
	}

	var uploadResp FileUploadResponse
	err = c.httpClient.DoRequest(ctx, opts, &uploadResp)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if uploadResp.GetCode() != "000000" {
		return nil, &Error{
			Code:    uploadResp.GetCode(),
			Message: uploadResp.Message,
		}
	}

	return &uploadResp, nil
}