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
		if r.URL.Query().Get("searchToken") != "a" {
			t.Errorf("missing searchToken param")
		}
		jsonResponse(w, `{"success":true,"data":[{"chatId":"c1@c.us","name":"Alice","isBusiness":false}],"limit":50,"skip":0,"total":1}`)
	})
	got, err := c.Contacts.List(context.Background(), ListContactsParams{SearchToken: "a"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].ChatID != "c1@c.us" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestContactsCommonGroups(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/contacts/c1@c.us/common_groups" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":[{"id":"g_1","name":"Team","description":null,"owner":null,"createdAt":null,"lastMessageAt":null,"participantCount":3,"restrict":null}],"limit":50,"skip":0,"total":1}`)
	})
	got, err := c.Contacts.CommonGroups(context.Background(), "c1@c.us")
	if err != nil {
		t.Fatalf("CommonGroups: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].ID != "g_1" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestContactsDecodeEmptyResponse(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &Contact{}); err != nil {
		t.Fatalf("Contact empty decode: %v", err)
	}
}

func TestContactsUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Contacts.List(context.Background(), ListContactsParams{})
	expectAuthError(t, err)
}

func TestContactsCommonGroupsWithParams(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/contacts/c1@c.us/common_groups" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "10" || r.URL.Query().Get("skip") != "20" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		jsonResponse(w, `{"success":true,"data":[{"id":"g1@g.us","name":"Team","description":null,"owner":null,"createdAt":null,"lastMessageAt":null,"participantCount":3,"restrict":null}],"limit":10,"skip":20,"total":21}`)
	})
	got, err := c.Contacts.CommonGroupsWithParams(context.Background(), "c1@c.us", CommonGroupsParams{Limit: 10, Skip: 20})
	if err != nil {
		t.Fatalf("CommonGroupsWithParams: %v", err)
	}
	if len(got.Data) != 1 || got.Total != 21 {
		t.Fatalf("unexpected: %+v", got)
	}
}
