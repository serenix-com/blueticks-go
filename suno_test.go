package blueticks

import (
	"context"
	"net/http"
	"testing"
)

func TestSunoGenerate(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/suno/songs" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"jobId":"job_1","clips":[{"id":"clip_1","status":"submitted","audioUrl":null,"imageUrl":null,"title":null,"durationSec":null,"model":null},{"id":"clip_2","status":"queued","audioUrl":null,"imageUrl":null,"title":null,"durationSec":null,"model":null}]}`)
	})
	got, err := c.Suno.Generate(context.Background(), SunoGenerateParams{Lyrics: "la la", Style: "pop"})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if got.JobID != "job_1" || len(got.Clips) != 2 || got.Clips[0].ID != "clip_1" {
		t.Fatalf("unexpected payload: %+v", got)
	}
}

func TestSunoGenerateUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Suno.Generate(context.Background(), SunoGenerateParams{Lyrics: "la", Style: "pop"})
	expectAuthError(t, err)
}

func TestSunoGetSong(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/suno/songs/clip_1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"id":"clip_1","status":"complete","audioUrl":"https://cdn/x.mp3","imageUrl":"https://cdn/x.png","title":"My Song","durationSec":123.5,"model":"v5.5"}`)
	})
	got, err := c.Suno.GetSong(context.Background(), "clip_1")
	if err != nil {
		t.Fatalf("GetSong: %v", err)
	}
	if got.ID != "clip_1" || got.Status != SunoStatusComplete || got.AudioURL == nil || *got.AudioURL != "https://cdn/x.mp3" {
		t.Fatalf("unexpected payload: %+v", got)
	}
}

func TestSunoGetSongUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Suno.GetSong(context.Background(), "clip_1")
	expectAuthError(t, err)
}

func TestSunoUpload(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/suno/uploads" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"uploadId":"up_1","status":"complete"}`)
	})
	got, err := c.Suno.Upload(context.Background(), SunoUploadParams{AudioURL: "https://example.com/a.mp3"})
	if err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if got.UploadID != "up_1" || got.Status != "complete" {
		t.Fatalf("unexpected payload: %+v", got)
	}
}

func TestSunoUploadUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Suno.Upload(context.Background(), SunoUploadParams{AudioURL: "https://example.com/a.mp3"})
	expectAuthError(t, err)
}

func TestSunoGetAccount(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/suno/account" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"creditsLeft":250,"monthlyLimit":2500,"monthlyUsage":100,"plan":"Pro Plan"}`)
	})
	got, err := c.Suno.GetAccount(context.Background())
	if err != nil {
		t.Fatalf("GetAccount: %v", err)
	}
	if got.CreditsLeft != 250 || got.Plan == nil || *got.Plan != "Pro Plan" {
		t.Fatalf("unexpected payload: %+v", got)
	}
}

func TestSunoGetAccountUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Suno.GetAccount(context.Background())
	expectAuthError(t, err)
}

// Lock in encoding/json's lenient decode — an empty `{}` decodes each Suno
// response type to zero values without panicking.
func TestSunoDecodesEmptyTypedResponses(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &SunoGenerateResult{}); err != nil {
		t.Fatalf("SunoGenerateResult: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &SunoClip{}); err != nil {
		t.Fatalf("SunoClip: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &SunoUpload{}); err != nil {
		t.Fatalf("SunoUpload: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &SunoAccount{}); err != nil {
		t.Fatalf("SunoAccount: %v", err)
	}
}
