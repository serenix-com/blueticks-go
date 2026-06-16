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
			jsonResponse(w, `{"id":"cmp_1","name":"Q1","audienceId":"aud_1","status":"pending","totalCount":0,"sentCount":0,"deliveredCount":0,"readCount":0,"failedCount":0,"createdAt":"2026-01-01T00:00:00Z"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/campaigns":
			jsonResponse(w, `{"data":[],"hasMore":false}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/campaigns/cmp_1":
			jsonResponse(w, `{"id":"cmp_1","name":"Q1","audienceId":"aud_1","status":"running","totalCount":1,"sentCount":1,"deliveredCount":0,"readCount":0,"failedCount":0,"createdAt":"2026-01-01T00:00:00Z"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/campaigns/cmp_1/pause":
			jsonResponse(w, `{"id":"cmp_1","name":"Q1","audienceId":"aud_1","status":"paused","totalCount":1,"sentCount":1,"deliveredCount":0,"readCount":0,"failedCount":0,"createdAt":"2026-01-01T00:00:00Z"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/campaigns/cmp_1/resume":
			jsonResponse(w, `{"id":"cmp_1","name":"Q1","audienceId":"aud_1","status":"running","totalCount":1,"sentCount":1,"deliveredCount":0,"readCount":0,"failedCount":0,"createdAt":"2026-01-01T00:00:00Z"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/campaigns/cmp_1/cancel":
			jsonResponse(w, `{"id":"cmp_1","name":"Q1","audienceId":"aud_1","status":"aborted","totalCount":1,"sentCount":1,"deliveredCount":0,"readCount":0,"failedCount":0,"createdAt":"2026-01-01T00:00:00Z"}`)
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
