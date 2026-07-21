package blueticks

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestServer spins up an httptest.Server and returns a *Client wired to
// it. The handler's r.URL.Path / r.Method drive routing in the per-resource
// tests.
func newTestServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c, err := NewClient(
		WithAPIKey("test_key"),
		WithBaseURL(srv.URL),
		WithMaxRetries(0),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c, srv
}

// unauthorizedHandler returns a 401 with the standard envelope.
func unauthorizedHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":{"code":"unauthorized","message":"bad key","request_id":"req_1"}}`))
}

// expectAuthError fails if err doesn't unwrap to ErrAuthentication.
func expectAuthError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected auth error, got nil")
	}
	if !errors.Is(err, ErrAuthentication) {
		t.Fatalf("expected ErrAuthentication, got %v", err)
	}
}

// jsonResponse writes a JSON body with status 200.
func jsonResponse(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(body))
}
