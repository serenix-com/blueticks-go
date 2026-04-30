package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestAccountRetrieve(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/account" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test_key" {
			t.Errorf("missing auth header: %q", got)
		}
		jsonResponse(w, `{"id":"acct_1","name":"Acme","timezone":"UTC","created_at":"2026-01-01T00:00:00Z"}`)
	})
	got, err := c.Account.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if got.ID != "acct_1" || got.Name != "Acme" {
		t.Fatalf("unexpected payload: %+v", got)
	}
}

func TestAccountRetrieveUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Account.Retrieve(context.Background())
	expectAuthError(t, err)
}
