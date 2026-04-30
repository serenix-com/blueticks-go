// Package main demonstrates verifying an incoming webhook delivery and
// dispatching on its event type.
//
//	BLUETICKS_WEBHOOK_SECRET=whsec_xxx go run ./examples/webhook
//
// Then point your Blueticks webhook subscription at this server (use ngrok or
// similar to expose the port).
package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	blueticks "github.com/serenix-com/blueticks-go"
)

func main() {
	secret := os.Getenv("BLUETICKS_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("BLUETICKS_WEBHOOK_SECRET is required")
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "read body", http.StatusBadRequest)
			return
		}
		event, err := blueticks.VerifyWebhook(body, r.Header, secret)
		if err != nil {
			if errors.Is(err, blueticks.ErrWebhookVerification) {
				log.Printf("rejected webhook: %v", err)
				http.Error(w, "bad signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}

		switch event.Type {
		case "message.delivered":
			log.Printf("delivered: %v", event.Data)
		case "message.read":
			log.Printf("read: %v", event.Data)
		default:
			log.Printf("event %s: %v", event.Type, event.Data)
		}
		w.WriteHeader(http.StatusOK)
	})

	addr := ":8081"
	log.Printf("listening on %s/webhook", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
