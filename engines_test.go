package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestEnginesStatus(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/engines" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"connected":true,"state":"CONNECTED"}`)
	})
	st, err := c.Engines.Status(context.Background())
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if !st.Connected {
		t.Fatalf("expected connected")
	}
}

func TestEnginesDecodeEmptyResponse(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &EngineStatus{}); err != nil {
		t.Fatalf("EngineStatus empty decode: %v", err)
	}
}

func TestEnginesUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Engines.Status(context.Background())
	expectAuthError(t, err)
}
