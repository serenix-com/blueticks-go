package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestEnginesAll(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/engines":
			jsonResponse(w, `{"connected":true,"state":"CONNECTED"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/engines/me":
			jsonResponse(w, `{"phone":"+1","name":"Bot","platform":"web"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/engines/logout":
			jsonResponse(w, `{"ok":true}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/engines/reload":
			jsonResponse(w, `{"ok":true}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	st, err := c.Engines.Status(ctx)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if !st.Connected {
		t.Fatalf("expected connected")
	}
	if _, err := c.Engines.Me(ctx); err != nil {
		t.Fatalf("Me: %v", err)
	}
	if _, err := c.Engines.Logout(ctx); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	if _, err := c.Engines.Reload(ctx); err != nil {
		t.Fatalf("Reload: %v", err)
	}
}

func TestEnginesUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Engines.Status(context.Background())
	expectAuthError(t, err)
}
