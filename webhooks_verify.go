package blueticks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DefaultWebhookTolerance is the default replay-window passed to
// VerifyWebhook when no override is supplied. 5 minutes matches the
// other Blueticks SDKs.
const DefaultWebhookTolerance = 300 * time.Second

// VerifyWebhookOption configures VerifyWebhook.
type VerifyWebhookOption func(*verifyWebhookConfig)

type verifyWebhookConfig struct {
	tolerance time.Duration
	now       func() time.Time
}

// WithVerifyTolerance overrides the replay-window used to reject stale
// signed timestamps. Pass 0 to disable the check (NOT recommended in
// production — the only valid use is back-filling old deliveries from a
// log replay tool).
func WithVerifyTolerance(d time.Duration) VerifyWebhookOption {
	return func(c *verifyWebhookConfig) { c.tolerance = d }
}

// withNowFunc is an internal hook used by tests so we can pin a fake
// "current time" without exposing a public knob nobody needs.
func withNowFunc(f func() time.Time) VerifyWebhookOption {
	return func(c *verifyWebhookConfig) { c.now = f }
}

// VerifyWebhook validates a Blueticks webhook delivery and returns the
// parsed event.
//
// Verification:
//  1. The Blueticks-Webhook-Timestamp and Blueticks-Webhook-Signature
//     headers must both be present.
//  2. The timestamp must be within the configured tolerance of the
//     server clock (default 5 minutes).
//  3. HMAC-SHA256(secret, "{timestamp}.{rawBody}") must match the v1
//     scheme value carried in the signature header.
//
// Pass the *raw* request body — the bytes the server signed. Unmarshalling
// and re-marshalling JSON will change byte ordering / whitespace and
// invalidate the signature.
func VerifyWebhook(payload []byte, headers http.Header, secret string, opts ...VerifyWebhookOption) (*WebhookEvent, error) {
	cfg := &verifyWebhookConfig{
		tolerance: DefaultWebhookTolerance,
		now:       time.Now,
	}
	for _, o := range opts {
		o(cfg)
	}

	timestamp := headerLookup(headers, "Blueticks-Webhook-Timestamp")
	signature := headerLookup(headers, "Blueticks-Webhook-Signature")
	if timestamp == "" || signature == "" {
		return nil, &WebhookVerificationError{Reason: "missing required headers"}
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return nil, &WebhookVerificationError{Reason: "invalid timestamp"}
	}
	if cfg.tolerance > 0 {
		now := cfg.now().Unix()
		drift := now - ts
		if drift < 0 {
			drift = -drift
		}
		if time.Duration(drift)*time.Second > cfg.tolerance {
			return nil, &WebhookVerificationError{Reason: "expired timestamp"}
		}
	}

	supplied := ""
	for _, part := range strings.Split(signature, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "v1=") {
			supplied = strings.TrimPrefix(part, "v1=")
			break
		}
	}
	if supplied == "" {
		return nil, &WebhookVerificationError{Reason: "invalid_signature: missing v1 scheme"}
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "."))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))

	suppliedBytes, decodeErr := hex.DecodeString(supplied)
	expectedBytes, _ := hex.DecodeString(expected)
	if decodeErr != nil || !hmac.Equal(expectedBytes, suppliedBytes) {
		return nil, &WebhookVerificationError{Reason: "invalid_signature: mismatch"}
	}

	out := &WebhookEvent{}
	if err := json.Unmarshal(payload, out); err != nil {
		return nil, &WebhookVerificationError{Reason: "invalid_signature: payload not JSON object"}
	}
	return out, nil
}

// headerLookup performs a case-insensitive header lookup. http.Header.Get
// already canonicalizes keys, but we accept either the canonical form or
// any cased variant (matching how proxies sometimes deliver headers).
func headerLookup(h http.Header, name string) string {
	if v := h.Get(name); v != "" {
		return v
	}
	lower := strings.ToLower(name)
	for k, v := range h {
		if strings.ToLower(k) == lower && len(v) > 0 {
			return v[0]
		}
	}
	return ""
}
