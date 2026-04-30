package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestContactsList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/contacts" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"data":[{"chat_id":"c1@c.us","name":"Alice","is_business":false}],"has_more":false}`)
	})
	got, err := c.Contacts.List(context.Background(), ListContactsParams{Query: "a"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].ChatID != "c1@c.us" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestContactsGetProfilePicture(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/contacts/c1@c.us/profile_picture" {
			t.Errorf("unexpected: %s", r.URL.Path)
		}
		jsonResponse(w, `{"url":"https://cdn/x.jpg"}`)
	})
	got, err := c.Contacts.GetProfilePicture(context.Background(), "c1@c.us")
	if err != nil {
		t.Fatalf("GetProfilePicture: %v", err)
	}
	if got.URL == "" {
		t.Fatalf("expected URL, got empty")
	}
}

func TestContactsUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Contacts.List(context.Background(), ListContactsParams{})
	expectAuthError(t, err)
}
