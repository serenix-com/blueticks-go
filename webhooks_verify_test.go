package blueticks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func sign(secret, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "."))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyWebhookHappyPath(t *testing.T) {
	secret := "whsec_test"
	body := []byte(`{"id":"evt_1","type":"message.delivered","created_at":"2026-01-01T00:00:00Z","data":{"message_id":"msg_1"}}`)
	now := time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC)
	ts := strconv.FormatInt(now.Unix(), 10)
	hex := sign(secret, ts, body)

	headers := http.Header{}
	headers.Set("Blueticks-Webhook-Timestamp", ts)
	headers.Set("Blueticks-Webhook-Signature", "v1="+hex)

	evt, err := VerifyWebhook(body, headers, secret, withNowFunc(func() time.Time { return now }))
	if err != nil {
		t.Fatalf("VerifyWebhook: %v", err)
	}
	if evt.ID != "evt_1" || evt.Type != "message.delivered" {
		t.Fatalf("unexpected event: %+v", evt)
	}
}

func TestVerifyWebhookMismatchedSignature(t *testing.T) {
	secret := "whsec_test"
	body := []byte(`{"id":"evt_1","type":"x","created_at":"2026-01-01T00:00:00Z","data":{}}`)
	now := time.Now()
	ts := strconv.FormatInt(now.Unix(), 10)

	headers := http.Header{}
	headers.Set("Blueticks-Webhook-Timestamp", ts)
	headers.Set("Blueticks-Webhook-Signature", "v1=deadbeef")

	_, err := VerifyWebhook(body, headers, secret, withNowFunc(func() time.Time { return now }))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrWebhookVerification) {
		t.Fatalf("expected ErrWebhookVerification, got %v", err)
	}
}

func TestVerifyWebhookExpiredTimestamp(t *testing.T) {
	secret := "whsec_test"
	body := []byte(`{"id":"evt_1","type":"x","created_at":"2026-01-01T00:00:00Z","data":{}}`)
	now := time.Now()
	staleTS := strconv.FormatInt(now.Add(-1*time.Hour).Unix(), 10)
	hex := sign(secret, staleTS, body)

	headers := http.Header{}
	headers.Set("Blueticks-Webhook-Timestamp", staleTS)
	headers.Set("Blueticks-Webhook-Signature", "v1="+hex)

	_, err := VerifyWebhook(body, headers, secret, withNowFunc(func() time.Time { return now }))
	if !errors.Is(err, ErrWebhookVerification) {
		t.Fatalf("expected verification error, got %v", err)
	}
}

func TestVerifyWebhookMissingHeaders(t *testing.T) {
	_, err := VerifyWebhook([]byte("{}"), http.Header{}, "x")
	if !errors.Is(err, ErrWebhookVerification) {
		t.Fatalf("expected verification error, got %v", err)
	}
}

func TestVerifyWebhookSignatureNotIncludingID(t *testing.T) {
	// Regression: prior reference SDKs occasionally drifted to "{id}.{ts}.{body}".
	// Confirm Go sticks to "{ts}.{body}" by signing with that exact formula and
	// ignoring any Blueticks-Webhook-Id header in the signed string.
	secret := "whsec_test"
	body := []byte(`{"id":"evt_1","type":"x","created_at":"2026-01-01T00:00:00Z","data":{}}`)
	now := time.Now()
	ts := strconv.FormatInt(now.Unix(), 10)
	hex := sign(secret, ts, body)

	headers := http.Header{}
	headers.Set("Blueticks-Webhook-Id", "wh_some_id_should_not_affect_signing")
	headers.Set("Blueticks-Webhook-Timestamp", ts)
	headers.Set("Blueticks-Webhook-Signature", "v1="+hex)

	if _, err := VerifyWebhook(body, headers, secret, withNowFunc(func() time.Time { return now })); err != nil {
		t.Fatalf("VerifyWebhook: %v", err)
	}
}
