package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestGroupsCRUD(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/groups":
			jsonResponse(w, `{"id":"g_1","name":"Team"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/groups/g_1":
			jsonResponse(w, `{"id":"g_1","name":"Team"}`)
		case r.Method == http.MethodPatch && r.URL.Path == "/v1/groups/g_1":
			jsonResponse(w, `{"id":"g_1","name":"Team Renamed"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/groups/g_1/members":
			jsonResponse(w, `{"ok":true}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/groups/g_1/members/c1@c.us":
			jsonResponse(w, `{"ok":true}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/groups/g_1/members/c1@c.us/admin":
			jsonResponse(w, `{"ok":true}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/groups/g_1/members/c1@c.us/admin":
			jsonResponse(w, `{"ok":true}`)
		case r.Method == http.MethodPut && r.URL.Path == "/v1/groups/g_1/picture":
			jsonResponse(w, `{"ok":true}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/groups/g_1/members/me":
			jsonResponse(w, `{"ok":true}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	if _, err := c.Groups.Create(ctx, CreateGroupParams{Name: "Team", Participants: []string{"c1@c.us"}}); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := c.Groups.Get(ctx, "g_1"); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if _, err := c.Groups.Update(ctx, "g_1", UpdateGroupParams{Name: "Team Renamed"}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if _, err := c.Groups.AddMember(ctx, "g_1", "c1@c.us"); err != nil {
		t.Fatalf("AddMember: %v", err)
	}
	if _, err := c.Groups.RemoveMember(ctx, "g_1", "c1@c.us"); err != nil {
		t.Fatalf("RemoveMember: %v", err)
	}
	if _, err := c.Groups.PromoteAdmin(ctx, "g_1", "c1@c.us"); err != nil {
		t.Fatalf("PromoteAdmin: %v", err)
	}
	if _, err := c.Groups.DemoteAdmin(ctx, "g_1", "c1@c.us"); err != nil {
		t.Fatalf("DemoteAdmin: %v", err)
	}
	if _, err := c.Groups.SetPicture(ctx, "g_1", SetGroupPictureParams{FileDataURL: "data:image/png;base64,AAA"}); err != nil {
		t.Fatalf("SetPicture: %v", err)
	}
	if _, err := c.Groups.Leave(ctx, "g_1"); err != nil {
		t.Fatalf("Leave: %v", err)
	}
}

func TestGroupsUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Groups.Get(context.Background(), "g_1")
	expectAuthError(t, err)
}
