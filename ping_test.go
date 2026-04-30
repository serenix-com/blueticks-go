package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestPingRetrieve(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/ping" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"account_id":"acct_1","key_prefix":"sk_live_xxx","scopes":["messages:write"]}`)
	})
	got, err := c.Ping.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if got.AccountID != "acct_1" || len(got.Scopes) != 1 {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestPingRetrieveUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Ping.Retrieve(context.Background())
	expectAuthError(t, err)
}
