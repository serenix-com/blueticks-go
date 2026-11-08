package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestWaLabelsList(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/wa-labels" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":[{"id":"1","name":"VIP","color":"#00ff00"},{"id":"2","name":"Lead","color":"#cccccc"}]}`)
	})
	got, err := c.WaLabels.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 || got[0].ID != "1" || got[1].Color != "#cccccc" {
		t.Fatalf("unexpected: %+v", got)
	}
}
