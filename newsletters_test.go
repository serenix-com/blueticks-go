package blueticks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestNewslettersCreate(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/newsletters" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var sent map[string]any
		_ = json.Unmarshal(body, &sent)
		if sent["name"] != "My Channel" {
			t.Errorf("unexpected body: %v", sent)
		}
		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, `{"id":"nl_1","name":"My Channel","description":null,"owner":null,"createdAt":"2026-01-01T00:00:00Z","subscribers":0,"invite":null,"verification":"UNVERIFIED"}`)
	})
	got, err := c.Newsletters.Create(context.Background(), CreateNewsletterParams{Name: "My Channel"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != "nl_1" {
		t.Fatalf("unexpected id: %q", got.ID)
	}
	if got.Name != "My Channel" {
		t.Fatalf("unexpected name: %q", got.Name)
	}
}

func TestNewslettersCreateUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Newsletters.Create(context.Background(), CreateNewsletterParams{Name: "Test"})
	expectAuthError(t, err)
}

func TestNewslettersList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/newsletters" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"data":[{"id":"nl_1","name":"My Channel","description":null,"owner":null,"createdAt":"2026-01-01T00:00:00Z","subscribers":100,"invite":"abc123","verification":"VERIFIED"}],"hasMore":false,"nextCursor":null}`)
	})
	got, err := c.Newsletters.List(context.Background(), ListNewslettersParams{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].ID != "nl_1" {
		t.Fatalf("unexpected data: %+v", got)
	}
	if got.HasMore {
		t.Fatalf("expected hasMore=false")
	}
}

func TestNewslettersListUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Newsletters.List(context.Background(), ListNewslettersParams{})
	expectAuthError(t, err)
}

func TestNewslettersRetrieve(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/newsletters/nl_1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"id":"nl_1","name":"My Channel","description":"A great channel","owner":"15551234567@s.whatsapp.net","createdAt":"2026-01-01T00:00:00Z","subscribers":42,"invite":"abc123","verification":"VERIFIED"}`)
	})
	got, err := c.Newsletters.Retrieve(context.Background(), "nl_1")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if got.ID != "nl_1" {
		t.Fatalf("unexpected id: %q", got.ID)
	}
}

func TestNewslettersRetrieveUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Newsletters.Retrieve(context.Background(), "nl_1")
	expectAuthError(t, err)
}

func TestNewslettersDecodeEmptyResponse(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &Newsletter{}); err != nil {
		t.Fatalf("Newsletter empty decode: %v", err)
	}
}
