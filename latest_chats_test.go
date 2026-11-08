package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestLatestChatsList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/latest-chats" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("numberOfChats") != "10" || r.URL.Query().Get("numberOfMessages") != "0" {
			t.Errorf("explicit zero must be sent; query: %v", r.URL.RawQuery)
		}
		jsonResponse(w, `{"success":true,"data":{"chats":[{"chatId":"c1@c.us","name":"Dana","chatType":"contact","pinned":true,"lastMessageAt":"2026-01-01T00:00:00Z","unreadCount":2,"markedUnread":false,"messages":[{"waMessageKey":{"fromMe":false,"id":"k1"},"chatId":"c1@c.us","from":"c1@c.us","type":"chat","fromMe":false,"text":"hi","source":"vault","secret":"s1"}]}],"totalChats":1,"numberOfChats":10,"numberOfMessages":0}}`)
	})
	n, m := 10, 0
	got, err := c.LatestChats.List(context.Background(), LatestChatsParams{NumberOfChats: &n, NumberOfMessages: &m})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if got.TotalChats != 1 || len(got.Chats) != 1 || got.Chats[0].ChatID != "c1@c.us" || len(got.Chats[0].Messages) != 1 {
		t.Fatalf("unexpected: %+v", got)
	}
	msg := got.Chats[0].Messages[0]
	if msg.Source == nil || *msg.Source != "vault" || msg.Secret == nil || *msg.Secret != "s1" {
		t.Fatalf("message source/secret: %+v", msg)
	}
}

func TestLatestChatsListDefaultsOmitQuery(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query, got %q", r.URL.RawQuery)
		}
		jsonResponse(w, `{"success":true,"data":{"chats":[],"totalChats":0,"numberOfChats":50,"numberOfMessages":20}}`)
	})
	if _, err := c.LatestChats.List(context.Background(), LatestChatsParams{}); err != nil {
		t.Fatalf("List: %v", err)
	}
}

func TestLatestChatsUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.LatestChats.List(context.Background(), LatestChatsParams{})
	expectAuthError(t, err)
}

func TestLatestChatsDecodeEmptyResponse(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) { jsonResponse(w, `{}`) })
	got, err := c.LatestChats.List(context.Background(), LatestChatsParams{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if got.Chats != nil || got.TotalChats != 0 {
		t.Fatalf("expected zero value, got %+v", got)
	}
}

func TestLatestChatsEngineTargeting(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("engineId") != "eng_2" || r.URL.Query().Get("fanoutLimit") != "2" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		jsonResponse(w, `{"success":true,"data":{"chats":[],"totalChats":0,"numberOfChats":50,"numberOfMessages":20}}`)
	})
	if _, err := c.LatestChats.List(context.Background(), LatestChatsParams{EngineTarget: EngineTarget{EngineID: "eng_2", FanoutLimit: 2}}); err != nil {
		t.Fatalf("List: %v", err)
	}
}
