package blueticks

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Sentinel errors. Use errors.Is to match against these — every concrete
// error type returned by the client wraps the appropriate sentinel so that
// callers can switch on category without importing concrete struct types.
var (
	ErrAuthentication      = errors.New("blueticks: authentication failed")
	ErrPermissionDenied    = errors.New("blueticks: permission denied")
	ErrNotFound            = errors.New("blueticks: not found")
	ErrBadRequest          = errors.New("blueticks: bad request")
	ErrRateLimit           = errors.New("blueticks: rate limited")
	ErrAPI                 = errors.New("blueticks: api error")
	ErrAPIConnection       = errors.New("blueticks: api connection error")
	ErrValidation          = errors.New("blueticks: validation error")
	ErrWebhookVerification = errors.New("blueticks: webhook verification failed")
)

// APIError is the canonical typed error for non-2xx responses from the Blueticks
// API. It carries the HTTP status, the parsed envelope (`error.code`,
// `error.message`, `error.request_id`), and the underlying sentinel for
// errors.Is matching.
type APIError struct {
	// StatusCode is the HTTP status code, or 0 for non-HTTP failures.
	StatusCode int
	// Code is `error.code` from the response envelope, when present.
	Code string
	// Message is `error.message` from the response envelope, or a fallback.
	Message string
	// RequestID is `error.request_id` from the response envelope, when present.
	RequestID string
	// RetryAfter is the parsed value of the `Retry-After` header on 429s,
	// in seconds. Zero when absent or unparseable.
	RetryAfter float64
	// sentinel is the category sentinel (ErrAuthentication, ErrNotFound, …).
	sentinel error
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Code != "" {
		return fmt.Sprintf("blueticks: %d %s: %s (request_id=%s)", e.StatusCode, e.Code, e.Message, e.RequestID)
	}
	return fmt.Sprintf("blueticks: %d: %s (request_id=%s)", e.StatusCode, e.Message, e.RequestID)
}

// Unwrap returns the category sentinel so errors.Is(err, ErrAuthentication)
// (and friends) work without the caller needing to type-assert.
func (e *APIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.sentinel
}

// APIConnectionError is a transport-level failure (DNS, TLS, connection reset,
// timeout) that didn't yield an HTTP response.
type APIConnectionError struct {
	Message string
	Cause   error
}

func (e *APIConnectionError) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil {
		return fmt.Sprintf("blueticks: api connection error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("blueticks: api connection error: %s", e.Message)
}

func (e *APIConnectionError) Unwrap() error { return ErrAPIConnection }

// ValidationError is raised when a server response can't be decoded into the
// expected typed model (missing required fields, wrong scalar type, etc.).
type ValidationError struct {
	Message string
	Cause   error
}

func (e *ValidationError) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil {
		return fmt.Sprintf("blueticks: validation error: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("blueticks: validation error: %s", e.Message)
}

func (e *ValidationError) Unwrap() error { return ErrValidation }

// WebhookVerificationError indicates a webhook signature or timestamp didn't
// validate. Always wraps ErrWebhookVerification for errors.Is matching.
type WebhookVerificationError struct {
	Reason string
}

func (e *WebhookVerificationError) Error() string {
	if e == nil {
		return ""
	}
	return "blueticks: webhook verification failed: " + e.Reason
}

func (e *WebhookVerificationError) Unwrap() error { return ErrWebhookVerification }

// fromEnvelope maps a non-2xx response body to an APIError carrying the
// appropriate category sentinel.
func fromEnvelope(statusCode int, body []byte, retryAfter float64) *APIError {
	sentinel := sentinelForStatus(statusCode)
	apiErr := &APIError{
		StatusCode: statusCode,
		RetryAfter: retryAfter,
		sentinel:   sentinel,
	}

	// Try to parse the standard envelope: {"error": {"code", "message", "request_id"}}.
	var envelope struct {
		Error struct {
			Code      string `json:"code"`
			Message   string `json:"message"`
			RequestID string `json:"requestId"`
		} `json:"error"`
	}
	if len(body) > 0 && json.Unmarshal(body, &envelope) == nil && envelope.Error.Code != "" {
		apiErr.Code = envelope.Error.Code
		apiErr.Message = envelope.Error.Message
		apiErr.RequestID = envelope.Error.RequestID
		return apiErr
	}

	// Fallback: render the body as a truncated string.
	raw := string(body)
	if len(raw) > 200 {
		raw = raw[:200] + "..."
	}
	apiErr.Message = raw
	return apiErr
}

func sentinelForStatus(status int) error {
	switch status {
	case 400, 422:
		return ErrBadRequest
	case 401:
		return ErrAuthentication
	case 403:
		return ErrPermissionDenied
	case 404:
		return ErrNotFound
	case 429:
		return ErrRateLimit
	default:
		return ErrAPI
	}
}
