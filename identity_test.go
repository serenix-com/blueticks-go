package blueticks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestIdentityGet(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/identity/123@lid" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":{"input":"123@lid","resolved":true,"canonical":"15550100@c.us","phone":"+15550100","lid":"123@lid","aliases":["15550100@c.us"],"type":"user","source":"engine"}}`)
	})
	got, err := c.Identity.Get(context.Background(), "123@lid")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !got.Resolved || got.Phone == nil || *got.Phone != "+15550100" || got.Type != IdentityTypeUser || got.Source != IdentitySourceEngine {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestIdentityGetUnresolved(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		jsonResponse(w, `{"success":true,"data":{"input":"9@lid","resolved":false,"aliases":[],"type":"unknown","source":"unknown"}}`)
	})
	got, err := c.Identity.Get(context.Background(), "9@lid")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Resolved || got.Phone != nil || got.Canonical != nil {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestIdentityResolve(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/identity/resolve" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var body struct {
			IDs []string `json:"ids"`
		}
		_ = json.Unmarshal(raw, &body)
		if len(body.IDs) != 2 || body.IDs[0] != "a@lid" {
			t.Errorf("unexpected body: %s", raw)
		}
		jsonResponse(w, `{"success":true,"data":[{"input":"a@lid","resolved":true,"aliases":[],"type":"user","source":"cache"},{"input":"b@g.us","resolved":true,"aliases":[],"type":"group","source":"engine"}]}`)
	})
	got, err := c.Identity.Resolve(context.Background(), []string{"a@lid", "b@g.us"})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if len(got.Data) != 2 || got.Data[0].Input != "a@lid" || got.Data[1].Type != IdentityTypeGroup {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestIdentityUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Identity.Get(context.Background(), "a@lid")
	expectAuthError(t, err)
	_, err = c.Identity.Resolve(context.Background(), []string{"a@lid"})
	expectAuthError(t, err)
}

func TestIdentityDecodeEmptyResponse(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) { jsonResponse(w, `{}`) })
	got, err := c.Identity.Get(context.Background(), "a@lid")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Input != "" || got.Resolved {
		t.Fatalf("expected zero value, got %+v", got)
	}
	res, err := c.Identity.Resolve(context.Background(), []string{"a@lid"})
	if err != nil || len(res.Data) != 0 {
		t.Fatalf("Resolve empty: %v %+v", err, res)
	}
}
