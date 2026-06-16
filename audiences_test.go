package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestAudiencesCreate(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/audiences" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"id":"aud_1","name":"VIPs","contactCount":0,"createdAt":"2026-01-01T00:00:00Z"}`)
	})
	got, err := c.Audiences.Create(context.Background(), CreateAudienceParams{Name: "VIPs"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != "aud_1" || got.Name != "VIPs" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestAudiencesListGetUpdateDelete(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/audiences":
			jsonResponse(w, `{"data":[{"id":"aud_1","name":"V","contactCount":0,"createdAt":"2026-01-01T00:00:00Z"}],"hasMore":false}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/audiences/aud_1":
			jsonResponse(w, `{"id":"aud_1","name":"V","contactCount":0,"createdAt":"2026-01-01T00:00:00Z"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/v1/audiences/aud_1":
			jsonResponse(w, `{"id":"aud_1","name":"V2","contactCount":0,"createdAt":"2026-01-01T00:00:00Z"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/audiences/aud_1":
			jsonResponse(w, `{"id":"aud_1","deleted":true}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	if _, err := c.Audiences.List(ctx, ListAudiencesParams{Limit: 10}); err != nil {
		t.Fatalf("List: %v", err)
	}
	if _, err := c.Audiences.Get(ctx, "aud_1", GetAudienceParams{}); err != nil {
		t.Fatalf("Get: %v", err)
	}
	got, err := c.Audiences.Update(ctx, "aud_1", UpdateAudienceParams{Name: "V2"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got.Name != "V2" {
		t.Fatalf("unexpected name: %q", got.Name)
	}
	del, err := c.Audiences.Delete(ctx, "aud_1")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if del.ID != "aud_1" || !del.Deleted {
		t.Fatalf("unexpected delete ref: %+v", del)
	}
}

func TestAudiencesAppendAndContactOps(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/audiences/aud_1/contacts":
			jsonResponse(w, `{"added":2,"contactCount":2}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/v1/audiences/aud_1/contacts/c_1":
			jsonResponse(w, `{"id":"c_1","to":"+1","addedAt":"2026-01-01T00:00:00Z"}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/audiences/aud_1/contacts/c_1":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	res, err := c.Audiences.AppendContacts(ctx, "aud_1", []AudienceContactInput{{To: "+1"}, {To: "+2"}})
	if err != nil {
		t.Fatalf("AppendContacts: %v", err)
	}
	if res.Added != 2 {
		t.Fatalf("expected added=2, got %d", res.Added)
	}
	if _, err := c.Audiences.UpdateContact(ctx, "aud_1", "c_1", UpdateAudienceContactParams{To: "+1"}); err != nil {
		t.Fatalf("UpdateContact: %v", err)
	}
	if err := c.Audiences.DeleteContact(ctx, "aud_1", "c_1"); err != nil {
		t.Fatalf("DeleteContact: %v", err)
	}
}

func TestAudiencesUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Audiences.List(context.Background(), ListAudiencesParams{})
	expectAuthError(t, err)
}
