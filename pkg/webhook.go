package interlace

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// WebhookClient handles webhook event processing
type WebhookClient struct {
	webhookSecret string
}

// NewWebhookClient creates a new webhook client
func NewWebhookClient(webhookSecret string) *WebhookClient {
	return &WebhookClient{
		webhookSecret: webhookSecret,
	}
}

// WebhookEvent represents a webhook event notification
type WebhookEvent struct {
	EventID   string                 `json:"eventId"`
	EventType string                 `json:"eventType"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// WebhookEventType defines the types of webhook events
const (
	// Card events
	EventCardCreated    = "card.created"
	EventCardActivated  = "card.activated"
	EventCardSuspended  = "card.suspended"
	EventCardDeleted    = "card.deleted"
	
	// Transaction events
	EventTransactionAuthorized = "transaction.authorized"
	EventTransactionDeclined   = "transaction.declined"
	EventTransactionCleared    = "transaction.cleared"
	
	// Transfer events
	EventTransferCreated   = "transfer.created"
	EventTransferCompleted = "transfer.completed"
	EventTransferFailed    = "transfer.failed"
	
	// Refund events
	EventRefundCreated   = "refund.created"
	EventRefundCompleted = "refund.completed"
	EventRefundFailed    = "refund.failed"
	
	// Account events
	EventAccountCreated  = "account.created"
	EventAccountUpdated  = "account.updated"
	EventAccountSuspended = "account.suspended"
	
	// Budget events
	EventBudgetCreated  = "budget.created"
	EventBudgetUpdated  = "budget.updated"
	EventBudgetExceeded = "budget.exceeded"
	
	// Payout events
	EventPayoutCreated   = "payout.created"
	EventPayoutCompleted = "payout.completed"
	EventPayoutFailed    = "payout.failed"
)

// VerifyWebhookSignature verifies the webhook signature
func (c *WebhookClient) VerifyWebhookSignature(payload []byte, signature string) bool {
	if c.webhookSecret == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(c.webhookSecret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// ParseWebhookEvent parses a webhook event from HTTP request
func (c *WebhookClient) ParseWebhookEvent(r *http.Request) (*WebhookEvent, error) {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	defer r.Body.Close()

	// Verify signature if webhook secret is set
	if c.webhookSecret != "" {
		signature := r.Header.Get("X-Interlace-Signature")
		if signature == "" {
			return nil, fmt.Errorf("missing webhook signature")
		}

		if !c.VerifyWebhookSignature(body, signature) {
			return nil, fmt.Errorf("invalid webhook signature")
		}
	}

	// Parse event
	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook event: %w", err)
	}

	return &event, nil
}

// GenerateWebhookSignature generates a signature for webhook payload (for testing)
func (c *WebhookClient) GenerateWebhookSignature(payload []byte) string {
	if c.webhookSecret == "" {
		return ""
	}

	mac := hmac.New(sha256.New, []byte(c.webhookSecret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// HandleWebhookEvent is a helper function to handle webhook events
// Returns true if event was handled successfully
type WebhookHandler func(event *WebhookEvent) error

// WebhookServer represents a webhook server configuration
type WebhookServer struct {
	client   *WebhookClient
	handlers map[string]WebhookHandler
}

// NewWebhookServer creates a new webhook server
func NewWebhookServer(webhookSecret string) *WebhookServer {
	return &WebhookServer{
		client:   NewWebhookClient(webhookSecret),
		handlers: make(map[string]WebhookHandler),
	}
}

// RegisterHandler registers a handler for a specific event type
func (s *WebhookServer) RegisterHandler(eventType string, handler WebhookHandler) {
	s.handlers[eventType] = handler
}

// HandleWebhook handles incoming webhook HTTP requests
func (s *WebhookServer) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	event, err := s.client.ParseWebhookEvent(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse webhook: %v", err), http.StatusBadRequest)
		return
	}

	// Find and execute handler
	handler, exists := s.handlers[event.EventType]
	if !exists {
		// No handler registered, but still return success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ignored","message":"No handler registered for event type"}`))
		return
	}

	// Execute handler
	if err := handler(event); err != nil {
		http.Error(w, fmt.Sprintf("Handler error: %v", err), http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
