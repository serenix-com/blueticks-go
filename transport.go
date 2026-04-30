package blueticks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Sentinel-ish constants for the retry policy. Kept unexported because the
// caller controls retry shape via WithMaxRetries and the HTTP client's
// timeout — these are internal tuning knobs.
const (
	defaultBaseURL    = "https://api.blueticks.co"
	defaultTimeout    = 30 * time.Second
	defaultMaxRetries = 3
	backoffBase       = 500 * time.Millisecond
	backoffCap        = 8 * time.Second
)

// retriableStatuses is the set of status codes that trigger a retry on
// idempotent methods. 429 retries even on POST when the server gives us a
// Retry-After hint, mirroring Stripe's behavior.
var retriableStatuses = map[int]struct{}{
	429: {},
	502: {},
	503: {},
	504: {},
}

// idempotentMethods is the set of HTTP methods we'll retry without an
// explicit Idempotency-Key header. PATCH is included to match
// Python/Node SDK behavior — most REST APIs treat partial updates as
// effectively idempotent in practice.
var idempotentMethods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodHead:    {},
	http.MethodOptions: {},
	http.MethodDelete:  {},
	http.MethodPatch:   {},
	http.MethodPut:     {},
}

// requestOptions are the per-call parameters threaded through Client.do.
type requestOptions struct {
	method         string
	path           string
	query          url.Values
	body           interface{}
	idempotencyKey string
}

// do executes a single HTTP request against the API with retries and
// envelope-based error mapping. The returned bytes are the raw response body
// (ready to json.Unmarshal into a typed model).
func (c *Client) do(ctx context.Context, opts requestOptions) ([]byte, error) {
	method := strings.ToUpper(opts.method)

	var bodyBytes []byte
	if opts.body != nil {
		var err error
		bodyBytes, err = json.Marshal(opts.body)
		if err != nil {
			return nil, &ValidationError{Message: "failed to marshal request body", Cause: err}
		}
	}

	endpoint := strings.TrimRight(c.baseURL, "/") + opts.path
	if len(opts.query) > 0 {
		endpoint += "?" + opts.query.Encode()
	}

	_, isIdempotent := idempotentMethods[method]
	if opts.idempotencyKey != "" {
		isIdempotent = true
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		// Build a fresh request each attempt — Go's http.Request body is
		// a single-shot io.ReadCloser that the transport may consume.
		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}
		req, err := http.NewRequestWithContext(ctx, method, endpoint, bodyReader)
		if err != nil {
			return nil, &APIConnectionError{Message: "failed to build request", Cause: err}
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.userAgent)
		if bodyBytes != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		if opts.idempotencyKey != "" {
			req.Header.Set("Idempotency-Key", opts.idempotencyKey)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			// Context cancellation should NOT be retried — the caller wants
			// to bail out, not back off into another attempt.
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			lastErr = &APIConnectionError{Message: "request failed", Cause: err}
			if attempt < c.maxRetries {
				if waitErr := sleepWithContext(ctx, backoffFor(attempt, 0)); waitErr != nil {
					return nil, waitErr
				}
				continue
			}
			return nil, lastErr
		}

		// Always drain + close, even on retry, to allow connection reuse.
		respBody, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			lastErr = &APIConnectionError{Message: "failed to read response body", Cause: readErr}
			if attempt < c.maxRetries {
				if waitErr := sleepWithContext(ctx, backoffFor(attempt, 0)); waitErr != nil {
					return nil, waitErr
				}
				continue
			}
			return nil, lastErr
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return respBody, nil
		}

		var retryAfter float64
		if resp.StatusCode == 429 {
			retryAfter = parseRetryAfter(resp.Header.Get("Retry-After"))
		}

		_, retriable := retriableStatuses[resp.StatusCode]
		if retriable && isIdempotent && attempt < c.maxRetries {
			wait := backoffFor(attempt, retryAfter)
			if waitErr := sleepWithContext(ctx, wait); waitErr != nil {
				return nil, waitErr
			}
			continue
		}

		return nil, fromEnvelope(resp.StatusCode, respBody, retryAfter)
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, &APIConnectionError{Message: fmt.Sprintf("exhausted %d retries", c.maxRetries)}
}

// decodeJSON unmarshals the raw bytes into out and returns a ValidationError
// on failure. Empty bodies leave out at its zero value (some endpoints —
// notably DELETE — return a non-meaningful body).
func decodeJSON(data []byte, out interface{}) error {
	if out == nil {
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return &ValidationError{Message: "failed to decode response body", Cause: err}
	}
	return nil
}

// parseRetryAfter parses the Retry-After header. Only the seconds form is
// supported — HTTP-date forms are rare in practice for this API. Returns 0
// when absent or unparseable; the caller falls back to exponential backoff.
func parseRetryAfter(raw string) float64 {
	if raw == "" {
		return 0
	}
	if v, err := strconv.ParseFloat(raw, 64); err == nil && v >= 0 {
		return v
	}
	return 0
}

// backoffFor returns the wait duration for a given attempt index. Honours an
// explicit Retry-After value when the server sent one; otherwise uses
// exponential backoff with full jitter, capped at backoffCap.
func backoffFor(attempt int, retryAfter float64) time.Duration {
	if retryAfter > 0 {
		return time.Duration(retryAfter * float64(time.Second))
	}
	exp := backoffBase * (1 << attempt)
	if exp > backoffCap {
		exp = backoffCap
	}
	jitter := time.Duration(rand.Int63n(int64(backoffBase)))
	return exp + jitter
}

// sleepWithContext sleeps for d, but returns early with ctx.Err() if the
// caller's context is cancelled — keeps cancellation responsive even mid-retry.
func sleepWithContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
