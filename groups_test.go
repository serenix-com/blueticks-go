package blueticks

import (
	"context"
	"net/http"
	"testing"
)

const groupFixture = `{"id":"g_1","name":"Team","description":null,"owner":null,"createdAt":null,"lastMessageAt":null,"participantCount":2,"announce":null,"restrict":null,"participants":[]}`

const groupListItemFixture = `{"id":"g_1","name":"Team","description":null,"owner":null,"createdAt":null,"lastMessageAt":null,"participantCount":2,"restrict":null}`

func TestGroupsList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/groups" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Errorf("limit query: %q", got)
		}
		jsonResponse(w, `{"success":true,"data":[`+groupListItemFixture+`],"limit":10,"skip":0,"total":1}`)
	})
	got, err := c.Groups.List(context.Background(), ListGroupsParams{Limit: 10})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].ID != "g_1" || got.Total != 1 {
		t.Fatalf("unexpected page: %+v", got)
	}
}

func TestGroupsListUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Groups.List(context.Background(), ListGroupsParams{})
	expectAuthError(t, err)
}

func TestGroupsCRUD(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		wrapped := `{"success":true,"data":` + groupFixture + `}`
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/groups":
			jsonResponse(w, wrapped)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/groups/g_1":
			if r.URL.Query().Get("include") != "participants" {
				t.Errorf("expected include=participants, got %q", r.URL.Query().Get("include"))
			}
			jsonResponse(w, wrapped)
		case r.Method == http.MethodPatch && r.URL.Path == "/v1/groups/g_1":
			jsonResponse(w, wrapped)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/groups/g_1/members":
			jsonResponse(w, wrapped)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/groups/g_1/members/c1@c.us":
			jsonResponse(w, wrapped)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/groups/g_1/members/c1@c.us/admin":
			jsonResponse(w, wrapped)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/groups/g_1/members/c1@c.us/admin":
			jsonResponse(w, wrapped)
		case r.Method == http.MethodPut && r.URL.Path == "/v1/groups/g_1/picture":
			jsonResponse(w, wrapped)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/groups/g_1/members/me":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	g, err := c.Groups.Create(ctx, CreateGroupParams{Name: "Team", Participants: []string{"c1@c.us"}})
	if err != nil || g.ID != "g_1" {
		t.Fatalf("Create: %v %+v", err, g)
	}
	if _, err := c.Groups.Get(ctx, "g_1", GetGroupParams{Include: "participants"}); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if _, err := c.Groups.Update(ctx, "g_1", UpdateGroupParams{Name: "Team Renamed"}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if _, err := c.Groups.AddMember(ctx, "g_1", AddMemberParams{ChatID: "c1@c.us"}); err != nil {
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
	if err := c.Groups.Leave(ctx, "g_1"); err != nil {
		t.Fatalf("Leave: %v", err)
	}
}

func TestGroupsDecodesEmptyTypedResponse(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &Group{}); err != nil {
		t.Fatalf("Group: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &GroupListItem{}); err != nil {
		t.Fatalf("GroupListItem: %v", err)
	}
}

func TestGroupsUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Groups.Get(context.Background(), "g_1", GetGroupParams{})
	expectAuthError(t, err)
}

func TestGroupsGetPicture(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/groups/g_1/picture" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":{"url":"https://pps.whatsapp.net/p.jpg"}}`)
	})
	got, err := c.Groups.GetPicture(context.Background(), "g_1")
	if err != nil {
		t.Fatalf("GetPicture: %v", err)
	}
	if got.URL == nil || *got.URL != "https://pps.whatsapp.net/p.jpg" {
		t.Fatalf("unexpected picture: %+v", got)
	}
}

func TestGroupsGetPictureNull(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, `{"success":true,"data":{"url":null}}`)
	})
	got, err := c.Groups.GetPicture(context.Background(), "g_1")
	if err != nil {
		t.Fatalf("GetPicture: %v", err)
	}
	if got.URL != nil {
		t.Fatalf("expected nil URL, got %q", *got.URL)
	}
}

func TestGroupsRemovePicture(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/v1/groups/g_1/picture" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":{"ok":true}}`)
	})
	got, err := c.Groups.RemovePicture(context.Background(), "g_1")
	if err != nil {
		t.Fatalf("RemovePicture: %v", err)
	}
	if !got.OK {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestGroupsPictureUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Groups.GetPicture(context.Background(), "g_1")
	expectAuthError(t, err)
	_, err = c.Groups.RemovePicture(context.Background(), "g_1")
	expectAuthError(t, err)
}
