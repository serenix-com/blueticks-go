# Changelog

All notable changes to the Blueticks Go SDK are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
