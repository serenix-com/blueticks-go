# Blueticks Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/serenix-com/blueticks-go.svg)](https://pkg.go.dev/github.com/serenix-com/blueticks-go)

Official Go client for the [Blueticks API](https://dev.blueticks.co/docs/api) — schedule WhatsApp messages, run campaigns, and verify webhooks from Go services.

## Install

```bash
go get github.com/serenix-com/blueticks-go
```

Requires Go 1.21 or newer.

## Quickstart

```go
package main

import (
	"context"
	"log"

	blueticks "github.com/serenix-com/blueticks-go"
)

func main() {
	c, err := blueticks.NewClient(blueticks.WithAPIKey("sk_live_xxx"))
	if err != nil { log.Fatal(err) }

	msg, err := c.Messages.Send(context.Background(), blueticks.SendMessageParams{
		To: "+15555550100", Text: "Hello!", IdempotencyKey: "msg-1",
	})
	if err != nil { log.Fatal(err) }
	log.Println("sent:", msg.ID)
}
```

`NewClient` falls back to the `BLUETICKS_API_KEY` and `BLUETICKS_BASE_URL`
environment variables when no explicit `WithAPIKey` / `WithBaseURL` option is
supplied.

## Resources

| Field          | Endpoint group                                           |
| -------------- | -------------------------------------------------------- |
| `client.Account`   | `/v1/account`                                       |
| `client.Ping`      | `/v1/ping`                                          |
| `client.Messages`  | `/v1/messages` — send, get, list                    |
| `client.Audiences` | `/v1/audiences` — CRUD + contact ops                |
| `client.Campaigns` | `/v1/campaigns` — create, list, pause/resume/cancel |
| `client.Webhooks`  | `/v1/webhooks` — register subscriptions, rotate secret |
| `client.Chats`     | `/v1/chats` — list, search, history, media          |
| `client.Contacts`  | `/v1/contacts`                                      |
| `client.Groups`    | `/v1/groups` — admin / membership                   |
| `client.Engines`   | `/v1/engines` — connection state, logout/reload     |
| `client.Utils`     | `/v1/utils` — phone validation, link previews       |

Every resource method takes `ctx context.Context` first, returns a typed
struct (or `*Page` envelope) plus an `error`.

## Webhook verification

```go
event, err := blueticks.VerifyWebhook(rawBody, r.Header, secret)
if errors.Is(err, blueticks.ErrWebhookVerification) {
    http.Error(w, "bad signature", http.StatusUnauthorized)
    return
}
```

Pass the *raw* request body — re-marshalling JSON will invalidate the HMAC
signature.

## Errors

All API failures return a typed `*APIError` carrying `StatusCode`, `Code`,
`Message`, and `RequestID`. Use `errors.Is` against the sentinel category to
branch:

```go
if errors.Is(err, blueticks.ErrAuthentication) { /* invalid key */ }
if errors.Is(err, blueticks.ErrRateLimit)       { /* back off */ }
if errors.Is(err, blueticks.ErrNotFound)        { /* missing resource */ }
```

The transport automatically retries 429/502/503/504 (with `Retry-After`
support) and connection errors on idempotent methods, with exponential
backoff and jitter. Cap defaults to 3 retries; override via
`WithMaxRetries`.

## Examples

- [`examples/send`](./examples/send) — minimal send example.
- [`examples/webhook`](./examples/webhook) — verify and dispatch incoming events.

## License

MIT — see [LICENSE](./LICENSE).
