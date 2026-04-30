// Package main demonstrates sending a WhatsApp message via the Blueticks Go SDK.
//
//	BLUETICKS_API_KEY=sk_live_xxx go run ./examples/send
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	blueticks "github.com/serenix-com/blueticks-go"
)

func main() {
	to := os.Getenv("BLUETICKS_TO")
	if to == "" {
		to = "+15555550100"
	}

	c, err := blueticks.NewClient()
	if err != nil {
		log.Fatalf("init: %v", err)
	}

	msg, err := c.Messages.Send(context.Background(), blueticks.SendMessageParams{
		To:             to,
		Text:           "Hello from the Blueticks Go SDK!",
		IdempotencyKey: "example-send-1",
	})
	if err != nil {
		log.Fatalf("send: %v", err)
	}
	fmt.Printf("queued message %s (status=%s)\n", msg.ID, msg.Status)
}
