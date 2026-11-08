package blueticks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestUploadsCreate(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/uploads" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var sent map[string]any
		_ = json.Unmarshal(raw, &sent)
		if sent["contentType"] != "application/pdf" || sent["sizeBytes"] != float64(482113) || sent["filename"] != "invoice.pdf" {
			t.Errorf("unexpected body: %s", raw)
		}
		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, `{"success":true,"data":{"fileId":"f_abc","uploadUrl":"https://storage.example/put?sig=x","method":"PUT","headers":{"Content-Type":"application/pdf","x-goog-content-length-range":"0,26214400"},"expiresAt":"2026-01-01T00:15:00Z","deleteAt":"2026-01-02T00:00:00Z","maxBytes":26214400,"contentType":"application/pdf"}}`)
	})
	got, err := c.Uploads.Create(context.Background(), CreateFileUploadParams{
		ContentType: "application/pdf", SizeBytes: 482113, Filename: "invoice.pdf",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.FileID != "f_abc" || got.Method != "PUT" || got.Headers["Content-Type"] != "application/pdf" || got.MaxBytes != 26214400 {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestUploadsCreateOmitsUnsetOptionals(t *testing.T) {
	c, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var sent map[string]any
		_ = json.Unmarshal(raw, &sent)
		if _, ok := sent["sizeBytes"]; ok {
			t.Errorf("unset sizeBytes must be omitted: %s", raw)
		}
		if _, ok := sent["filename"]; ok {
			t.Errorf("unset filename must be omitted: %s", raw)
		}
		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, `{"success":true,"data":{"fileId":"f_1","uploadUrl":"https://storage.example/put","method":"PUT","headers":{},"expiresAt":"2026-01-01T00:15:00Z","deleteAt":"2026-01-02T00:00:00Z","maxBytes":1,"contentType":"image/png"}}`)
	})
	if _, err := c.Uploads.Create(context.Background(), CreateFileUploadParams{ContentType: "image/png"}); err != nil {
		t.Fatalf("Create: %v", err)
	}
}
