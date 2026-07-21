package blueticks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestSunoAccount(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/suno/account" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":{"creditsLeft":250.5,"monthlyLimit":500,"monthlyUsage":249.5,"plan":"Pro Plan"}}`)
	})
	got, err := c.Suno.Account(context.Background())
	if err != nil {
		t.Fatalf("Account: %v", err)
	}
	if got.CreditsLeft != 250.5 || got.Plan == nil || *got.Plan != "Pro Plan" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestSunoGenerateSong(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/suno/songs" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["lyrics"] != "la la" || body["style"] != "pop" {
			t.Errorf("unexpected body: %+v", body)
		}
		jsonResponse(w, `{"success":true,"data":{"jobId":"job_1","clips":[{"id":"clip_1","status":"queued","audioUrl":null,"imageUrl":null,"title":null,"durationSec":null,"model":"v5.5"}]}}`)
	})
	got, err := c.Suno.GenerateSong(context.Background(), GenerateSongParams{Lyrics: "la la", Style: "pop"})
	if err != nil {
		t.Fatalf("GenerateSong: %v", err)
	}
	if got.JobID != "job_1" || len(got.Clips) != 1 || got.Clips[0].ID != "clip_1" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestSunoGetSong(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/suno/songs/clip_1" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":{"id":"clip_1","status":"complete","audioUrl":"https://ex.com/s.mp3","imageUrl":null,"title":"My Song","durationSec":123.4,"model":"v5.5"}}`)
	})
	got, err := c.Suno.GetSong(context.Background(), "clip_1")
	if err != nil {
		t.Fatalf("GetSong: %v", err)
	}
	if got.Status != SongStatusComplete || got.AudioURL == nil || *got.AudioURL != "https://ex.com/s.mp3" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestSunoCreateUpload(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/suno/uploads" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		jsonResponse(w, `{"success":true,"data":{"uploadId":"up_1","status":"complete"}}`)
	})
	got, err := c.Suno.CreateUpload(context.Background(), CreateUploadParams{AudioURL: "https://ex.com/a.mp3"})
	if err != nil {
		t.Fatalf("CreateUpload: %v", err)
	}
	if got.UploadID != "up_1" || got.Status != "complete" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestSunoDecodesEmptyTypedResponses(t *testing.T) {
	if err := decodeJSON([]byte(`{}`), &SunoAccount{}); err != nil {
		t.Fatalf("SunoAccount: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &SongClip{}); err != nil {
		t.Fatalf("SongClip: %v", err)
	}
	if err := decodeJSON([]byte(`{}`), &GenerateSongResult{}); err != nil {
		t.Fatalf("GenerateSongResult: %v", err)
	}
}

func TestSunoUnauthorized(t *testing.T) {
	c, _ := newTestServer(t, unauthorizedHandler)
	_, err := c.Suno.Account(context.Background())
	expectAuthError(t, err)
}
