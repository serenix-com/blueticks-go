package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestRealtimeToken(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/realtime/token" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, `{"success":true,"data":{"brokerUrl":"wss://mqtt.example/mqtt","username":"u","password":"p","expiresAt":"2030-01-01T00:00:00Z","expiresInSeconds":1000,"topics":{"messages":"acct/messages"}}}`)
	})
	got, err := c.Realtime.Token(context.Background())
	if err != nil {
		t.Fatalf("Token: %v", err)
	}
	if got.BrokerURL != "wss://mqtt.example/mqtt" || got.Username != "u" || got.Password != "p" || got.ExpiresInSeconds != 1000 || got.Topics.Messages != "acct/messages" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestRealtimeTokenUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Realtime.Token(context.Background())
	expectAuthError(t, err)
}

func TestRealtimeTokenDecodeEmptyResponse(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) { jsonResponse(w, `{}`) })
	got, err := c.Realtime.Token(context.Background())
	if err != nil {
		t.Fatalf("Token: %v", err)
	}
	if got.BrokerURL != "" || got.Topics.Messages != "" {
		t.Fatalf("expected zero value, got %+v", got)
	}
}
