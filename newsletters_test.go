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
		jsonResponse(w, `{"success":true,"data":{"newsletterId":"123@newsletter","name":"My Channel","description":null,"createdAt":"2026-01-01T00:00:00Z","subscribers":0,"invite":null,"verification":"UNVERIFIED"}}`)
	})
	got, err := c.Newsletters.Create(context.Background(), CreateNewsletterParams{Name: "My Channel"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.NewsletterID != "123@newsletter" || got.Name != "My Channel" {
		t.Fatalf("unexpected: %+v", got)
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
		jsonResponse(w, `{"success":true,"data":[{"chatId":"123@newsletter","name":"My Channel","description":null,"createdAt":"2026-01-01T00:00:00Z","subscribers":100,"invite":"abc123","verification":"VERIFIED"}],"limit":50,"skip":0,"total":1}`)
	})
	got, err := c.Newsletters.List(context.Background(), ListNewslettersParams{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].ChatID != "123@newsletter" {
		t.Fatalf("unexpected data: %+v", got)
	}
}

func TestNewslettersListUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Newsletters.List(context.Background(), ListNewslettersParams{})
	expectAuthError(t, err)
}

func TestNewslettersRetrieve(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/newsletters/123@newsletter" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":{"newsletterId":"123@newsletter","name":"My Channel","description":"A great channel","createdAt":"2026-01-01T00:00:00Z","subscribers":42,"invite":"abc123","verification":"VERIFIED"}}`)
	})
	got, err := c.Newsletters.Retrieve(context.Background(), "123@newsletter")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if got.NewsletterID != "123@newsletter" {
		t.Fatalf("unexpected id: %q", got.NewsletterID)
	}
}

func TestNewslettersRetrieveUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Newsletters.Retrieve(context.Background(), "123@newsletter")
	expectAuthError(t, err)
}

func TestNewslettersDecodeEmptyResponse(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &Newsletter{}); err != nil {
		t.Fatalf("Newsletter empty decode: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &NewsletterListItem{}); err != nil {
		t.Fatalf("NewsletterListItem empty decode: %v", err)
	}
}
