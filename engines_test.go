package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestEnginesList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/engines" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":[{"id":"eng_1","type":"gateway","connected":true,"state":"CONNECTED","stream":"CONNECTED","hasSynced":true}]}`)
	})
	engines, err := c.Engines.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(engines) != 1 || engines[0].ID != "eng_1" || !engines[0].Connected {
		t.Fatalf("unexpected engines: %+v", engines)
	}
	if engines[0].State == nil || *engines[0].State != "CONNECTED" {
		t.Fatalf("unexpected state: %+v", engines[0].State)
	}
}

// TestEnginesOptionalFieldsNullAndAbsent locks in that optional response fields
// tolerate BOTH an explicit JSON null (disconnected engine) AND being omitted
// (connected engine) — both decode to a nil pointer, no error. Confirmed live:
// /v1/engines omits state/stream/hasSynced when CONNECTED, returns them as
// explicit null when DISCONNECTED.
func TestEnginesOptionalFieldsNullAndAbsent(t *testing.T) {
	// Explicit null.
	var withNull EngineStatus
	if err := decodeJSON([]byte(`{"id":"e","type":"gateway","connected":false,"state":null,"stream":null,"hasSynced":null}`), &withNull); err != nil {
		t.Fatalf("null decode: %v", err)
	}
	if withNull.State != nil || withNull.Stream != nil || withNull.HasSynced != nil {
		t.Fatalf("expected nil pointers for explicit null, got %+v", withNull)
	}
	// Omitted entirely.
	var withAbsent EngineStatus
	if err := decodeJSON([]byte(`{"id":"e","type":"gateway","connected":true}`), &withAbsent); err != nil {
		t.Fatalf("absent decode: %v", err)
	}
	if withAbsent.State != nil || withAbsent.Stream != nil || withAbsent.HasSynced != nil {
		t.Fatalf("expected nil pointers for absent fields, got %+v", withAbsent)
	}
}

func TestEnginesDecodeEmptyResponse(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &EngineStatus{}); err != nil {
		t.Fatalf("EngineStatus empty decode: %v", err)
	}
}

func TestEnginesUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Engines.List(context.Background())
	expectAuthError(t, err)
}
