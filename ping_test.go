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
		jsonResponse(w, `{"success":true,"data":{"api":"ok","accountId":"acct_1","whatsappConnections":[{"id":"eng_1","type":"gateway","connected":true}]}}`)
	})
	got, err := c.Ping.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if got.API != "ok" || got.AccountID != "acct_1" {
		t.Fatalf("unexpected: %+v", got)
	}
	if len(got.WhatsappConnections) != 1 || got.WhatsappConnections[0].ID != "eng_1" {
		t.Fatalf("unexpected connections: %+v", got.WhatsappConnections)
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
