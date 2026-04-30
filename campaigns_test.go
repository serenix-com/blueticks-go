package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestCampaignsCRUDAndTransitions(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/campaigns":
			jsonResponse(w, `{"id":"cmp_1","name":"Q1","audience_id":"aud_1","status":"pending","total_count":0,"sent_count":0,"delivered_count":0,"read_count":0,"failed_count":0,"created_at":"2026-01-01T00:00:00Z"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/campaigns":
			jsonResponse(w, `{"data":[],"has_more":false}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/campaigns/cmp_1":
			jsonResponse(w, `{"id":"cmp_1","name":"Q1","audience_id":"aud_1","status":"running","total_count":1,"sent_count":1,"delivered_count":0,"read_count":0,"failed_count":0,"created_at":"2026-01-01T00:00:00Z"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/campaigns/cmp_1/pause":
			jsonResponse(w, `{"id":"cmp_1","name":"Q1","audience_id":"aud_1","status":"paused","total_count":1,"sent_count":1,"delivered_count":0,"read_count":0,"failed_count":0,"created_at":"2026-01-01T00:00:00Z"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/campaigns/cmp_1/resume":
			jsonResponse(w, `{"id":"cmp_1","name":"Q1","audience_id":"aud_1","status":"running","total_count":1,"sent_count":1,"delivered_count":0,"read_count":0,"failed_count":0,"created_at":"2026-01-01T00:00:00Z"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/campaigns/cmp_1/cancel":
			jsonResponse(w, `{"id":"cmp_1","name":"Q1","audience_id":"aud_1","status":"aborted","total_count":1,"sent_count":1,"delivered_count":0,"read_count":0,"failed_count":0,"created_at":"2026-01-01T00:00:00Z"}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	if _, err := c.Campaigns.Create(ctx, CreateCampaignParams{Name: "Q1", AudienceID: "aud_1", Text: "hi"}); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := c.Campaigns.List(ctx, ListCampaignsParams{}); err != nil {
		t.Fatalf("List: %v", err)
	}
	if _, err := c.Campaigns.Get(ctx, "cmp_1"); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if _, err := c.Campaigns.Pause(ctx, "cmp_1"); err != nil {
		t.Fatalf("Pause: %v", err)
	}
	if _, err := c.Campaigns.Resume(ctx, "cmp_1"); err != nil {
		t.Fatalf("Resume: %v", err)
	}
	if _, err := c.Campaigns.Cancel(ctx, "cmp_1"); err != nil {
		t.Fatalf("Cancel: %v", err)
	}
}

func TestCampaignsUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Campaigns.List(context.Background(), ListCampaignsParams{})
	expectAuthError(t, err)
}
