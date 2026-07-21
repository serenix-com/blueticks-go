package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func campaignFixture(status string) string {
	return `{"cmpId":"cmp_1","name":"Q1","audienceId":"aud_1","status":"` + status + `","totalCount":1,"sentCount":1,"deliveredCount":0,"readCount":0,"failedCount":0,"createdAt":"2026-01-01T00:00:00Z","startedAt":null,"completedAt":null,"abortedAt":null}`
}

func TestCampaignsCRUDAndTransitions(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/campaigns":
			jsonResponse(w, `{"success":true,"data":`+campaignFixture("pending")+`}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/campaigns":
			jsonResponse(w, `{"success":true,"data":[`+campaignFixture("running")+`],"limit":50,"skip":0,"total":1}`)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/campaigns/cmp_1":
			jsonResponse(w, `{"success":true,"data":`+campaignFixture("running")+`}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/campaigns/cmp_1/pause":
			jsonResponse(w, `{"success":true,"data":`+campaignFixture("paused")+`}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/campaigns/cmp_1/resume":
			jsonResponse(w, `{"success":true,"data":`+campaignFixture("running")+`}`)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/campaigns/cmp_1/cancel":
			jsonResponse(w, `{"success":true,"data":`+campaignFixture("aborted")+`}`)
		default:
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
	})
	ctx := context.Background()
	created, err := c.Campaigns.Create(ctx, CreateCampaignParams{Name: "Q1", AudienceID: "aud_1", Text: "hi"})
	if err != nil || created.CmpID != "cmp_1" {
		t.Fatalf("Create: %v %+v", err, created)
	}
	page, err := c.Campaigns.List(ctx, ListCampaignsParams{})
	if err != nil || page.Total != 1 {
		t.Fatalf("List: %v %+v", err, page)
	}
	if _, err := c.Campaigns.Get(ctx, "cmp_1"); err != nil {
		t.Fatalf("Get: %v", err)
	}
	paused, err := c.Campaigns.Pause(ctx, "cmp_1")
	if err != nil || paused.Status != CampaignStatusPaused {
		t.Fatalf("Pause: %v %+v", err, paused)
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

func TestCampaignDecodeEmptyResponse(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &Campaign{}); err != nil {
		t.Fatalf("Campaign empty decode: %v", err)
	}
}
