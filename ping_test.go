package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestPingRetrieve(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/ping" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"account_id":"acct_1","key_prefix":"sk_live_1","scopes":["messages:write"]}`)
	})
	got, err := c.Ping.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if got.AccountID != "acct_1" {
		t.Fatalf("unexpected account_id: %q", got.AccountID)
	}
	if got.KeyPrefix != "sk_live_1" {
		t.Fatalf("unexpected key_prefix: %q", got.KeyPrefix)
	}
	if len(got.Scopes) != 1 || got.Scopes[0] != "messages:write" {
		t.Fatalf("unexpected scopes: %v", got.Scopes)
	}
}

func TestPingRetrieveUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Ping.Retrieve(context.Background())
	expectAuthError(t, err)
}

func TestPingDecodeEmpty(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &Ping{}); err != nil {
		t.Fatalf("Ping empty decode: %v", err)
	}
}
