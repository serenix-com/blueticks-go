# Changelog

All notable changes to the Blueticks Go SDK are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v1.0.1 — 2026-04-30

### Changed (semver-minor on the type level, but additive at runtime)

- Nullable scalar fields now use pointers to preserve the wire-level
  `null` vs `""`/`0` distinction declared in the OpenAPI spec:
  - `Account.Timezone`: `string` → `*string`
  - `Chat.Name`, `Chat.LastMessageAt`: `string` → `*string`
  - `Chat.UnreadCount`: `int` → `*int`
  - `*Page.NextCursor` (all 8 paginated envelopes): `string` → `*string`

  Previously, a server-returned `null` decoded to the zero value, losing
  information ("server said null" vs "server said empty string"). Now
  `nil` means null. Callers reading these fields need to deref or
  nil-check.

  Spotted by the integration smoke against stg-api.blueticks.co — Ruby
  SDK preserved the distinction (`nil` vs `""`); Go did not. Now they
  match.

## v1.0.0 — 2026-04-30

Initial release.

### Added

- `blueticks.NewClient` with `WithAPIKey`, `WithBaseURL`, `WithHTTPClient`,
  `WithUserAgent`, `WithMaxRetries` options.
- Environment fallbacks: `BLUETICKS_API_KEY`, `BLUETICKS_BASE_URL`.
- Resource clients: `Account`, `Ping`, `Messages`, `Audiences`, `Campaigns`,
  `Webhooks`, `Chats`, `Contacts`, `Groups`, `Engines`, `Utils`.
- `VerifyWebhook` helper — HMAC-SHA256 over `"{timestamp}.{rawBody}"` with
  `v1=` scheme detection, replay-window check, and case-insensitive header lookup.
- Typed errors (`APIError`, `WebhookVerificationError`, `APIConnectionError`,
  `ValidationError`) plus sentinel categories (`ErrAuthentication`,
  `ErrPermissionDenied`, `ErrNotFound`, `ErrBadRequest`, `ErrRateLimit`,
  `ErrAPI`, `ErrAPIConnection`, `ErrWebhookVerification`).
- Retry transport with exponential backoff + jitter, `Retry-After` honoring,
  and method-aware idempotency.
- Examples for send + webhook flows.
