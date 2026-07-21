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
		jsonResponse(w, `{"success":true,"data":{"id":"aud_1","name":"VIPs","contactCount":0,"createdAt":"2026-01-01T00:00:00Z"}}`)
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
			jsonResponse(w, `{"success":true,"data":[{"id":"aud_1","name":"V","contactCount":0,"createdAt":"2026-01-01T00:00:00Z"}],"limit":50,"skip":0,"total":1}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/audiences/aud_1":
			jsonResponse(w, `{"success":true,"data":{"id":"aud_1","name":"V","contactCount":0,"createdAt":"2026-01-01T00:00:00Z"}}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/v1/audiences/aud_1":
			jsonResponse(w, `{"success":true,"data":{"id":"aud_1","name":"V2","contactCount":0,"createdAt":"2026-01-01T00:00:00Z"}}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/audiences/aud_1":
			jsonResponse(w, `{"success":true,"data":{"id":"aud_1","deleted":true}}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	page, err := c.Audiences.List(ctx, ListAudiencesParams{Limit: 10})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if page.Total != 1 || len(page.Data) != 1 {
		t.Fatalf("unexpected page: %+v", page)
	}
	if _, err := c.Audiences.Get(ctx, "aud_1"); err != nil {
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
			jsonResponse(w, `{"success":true,"data":{"added":2,"contactCount":2}}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/v1/audiences/aud_1/contacts/c_1":
			jsonResponse(w, `{"success":true,"data":{"id":"c_1","to":"+1","variables":{},"addedAt":"2026-01-01T00:00:00Z"}}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/audiences/aud_1/contacts/c_1":
			jsonResponse(w, `{"success":true,"data":{"id":"c_1"}}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	res, err := c.Audiences.AppendContacts(ctx, "aud_1", []ContactInput{{To: "+1"}, {To: "+2"}})
	if err != nil {
		t.Fatalf("AppendContacts: %v", err)
	}
	if res.Added != 2 {
		t.Fatalf("expected added=2, got %d", res.Added)
	}
	if _, err := c.Audiences.UpdateContact(ctx, "aud_1", "c_1", UpdateAudienceContactParams{To: strptr("+1")}); err != nil {
		t.Fatalf("UpdateContact: %v", err)
	}
	del, err := c.Audiences.DeleteContact(ctx, "aud_1", "c_1")
	if err != nil {
		t.Fatalf("DeleteContact: %v", err)
	}
	if del.ID != "c_1" {
		t.Fatalf("unexpected deleted ref: %+v", del)
	}
}

func TestAudiencesUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Audiences.List(context.Background(), ListAudiencesParams{})
	expectAuthError(t, err)
}

func TestAudienceDecodeEmptyResponse(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &Audience{}); err != nil {
		t.Fatalf("Audience empty decode: %v", err)
	}
}

func strptr(s string) *string { return &s }
