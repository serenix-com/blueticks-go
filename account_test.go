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
		jsonResponse(w, `{"success":true,"data":{"id":"acct_1","name":"Acme","userEmail":"a@b.co","timezone":"UTC","createdAt":"2026-01-01T00:00:00Z"}}`)
	})
	got, err := c.Account.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if got.ID != "acct_1" || got.Name != "Acme" {
		t.Fatalf("unexpected payload: %+v", got)
	}
	if got.UserEmail == nil || *got.UserEmail != "a@b.co" {
		t.Fatalf("unexpected userEmail: %+v", got.UserEmail)
	}
}

func TestAccountRetrieveUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Account.Retrieve(context.Background())
	expectAuthError(t, err)
}

func TestAccountDecodeEmptyResponse(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &Account{}); err != nil {
		t.Fatalf("Account empty decode: %v", err)
	}
}
