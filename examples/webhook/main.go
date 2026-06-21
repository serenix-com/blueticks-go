// Package main demonstrates registering and listing webhook subscriptions.
//
//	BLUETICKS_API_KEY=sk_live_xxx go run ./examples/webhook
package main

import (
	"context"
	"fmt"
	"log"

	blueticks "github.com/serenix-com/blueticks-go"
)

func main() {
	c, err := blueticks.NewClient()
	if err != nil {
		log.Fatalf("init: %v", err)
	}
	ctx := context.Background()

	wh, err := c.Webhooks.Create(ctx, blueticks.CreateWebhookParams{
		URL:    "https://example.com/webhook",
		Events: []string{"message.delivered", "message.read"},
	})
	if err != nil {
		log.Fatalf("create webhook: %v", err)
	}
	fmt.Printf("created webhook %s -> %s\n", wh.ID, wh.URL)

	page, err := c.Webhooks.List(ctx, blueticks.ListWebhooksParams{})
	if err != nil {
		log.Fatalf("list webhooks: %v", err)
	}
	for _, w := range page.Data {
		fmt.Printf("webhook %s status=%s events=%v\n", w.ID, w.Status, w.Events)
	}
}
