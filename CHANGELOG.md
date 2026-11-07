# Changelog

All notable changes to the Blueticks Go SDK are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v2.1.0 — 2026-11-07

### Added

- Group picture: `Groups.GetPicture` (get the current picture URL) and `Groups.RemovePicture` (remove the picture; Gateway/Baileys engine only — Chrome-extension-only accounts get `501 unsupported_command`). `Groups.SetPicture` already existed.

## v1.2.0 — 2026-05-22

OpenAPI parity pass. The SDK now matches `backend/openapi.json`
operation-for-operation; an engineless drift check
(`.github/workflows/sdk-spec-drift.yml`) gates future regressions. The
`/v1/*` surface is pre-release — none of these changes affect production
callers yet.

### Changed

- `Messages.Send` now takes the discriminated body shape matching the
  backend's strict `anyOf` (BE#50). `SendMessageParams` accepts all variant
  fields; caller sets `Type` to one of `MessageTypeText` / `MessageTypeMedia`
  / `MessageTypePoll`:
  ```go
  c.Messages.Send(ctx, blueticks.SendMessageParams{
      To: "+1...", Type: blueticks.MessageTypeText, Text: ptr.String("hi"),
  })
  c.Messages.Send(ctx, blueticks.SendMessageParams{
      To: "+1...", Type: blueticks.MessageTypeMedia,
      Media: &blueticks.SendMediaObject{URL: "...", Kind: "image"},
  })
  c.Messages.Send(ctx, blueticks.SendMessageParams{
      To: "+1...", Type: blueticks.MessageTypePoll,
      Poll: &blueticks.SendPollObject{Question: "...", Options: []string{...}},
  })
  ```
  Explicit typed convenience structs (`SendMessageTextParams`,
  `SendMessageMediaParams`, `SendMessagePollParams`) are also available.
- `Newsletter` struct: 3 fields → 8 fields matching the spec (`ID`, `Name`,
  `Description *string`, `Owner *string`, `CreatedAt *string`,
  `Subscribers *int`, `Invite *string`, `Verification *string`).
- `CreateNewsletterParams.Description` is now `*string` (was `string`), to
  distinguish "absent" from the empty string.

### Added

- `Newsletters.List(ctx, ListNewslettersParams)` — `GET /v1/newsletters` (returns `*NewslettersPage`)
- `Newsletters.Retrieve(ctx, id)` — `GET /v1/newsletters/{id}`
- `Newsletters.Create(ctx, CreateNewsletterParams)` — now returns typed `*Newsletter`
- `Ping.Retrieve(ctx)` — typed `*Ping` (`AccountID`, `KeyPrefix`, `Scopes`).
  The `ping` resource had been deleted in a prior interim regen; restored here.
- `Groups.List` — `GET /v1/groups` (was documented but absent for ~9 days)
- `Message` now exposes `Key *string`, `Type string`, `MediaKind *string`,
  `PollQuestion *string`, `LinkPreview *MessageLinkPreview`

### Removed

- `Engines.Me`, `Engines.Logout`, `Engines.Reload`
- `Contacts.GetProfilePicture`
- `Utils.ValidatePhone`, `Utils.LinkPreview`

## v1.1.0 — 2026-05-01

### Added

- New `ScheduledMessages` resource backing `/v1/scheduled-messages`:
  `List`, `Retrieve`, `Update`, `Delete`. Plus `ScheduledMessage`,
  `ScheduledMessagesPage`, `ListScheduledMessagesParams`,
  `UpdateScheduledMessageParams` types.
- `Group` gains `LastMessageAt *string` and `Participants []GroupParticipant`
  (returned when the engine includes participant data).
- New `GroupParticipant` type for the `Group.Participants` array.
- `ChatMedia.MediaUnavailable` enum gains `"awaiting_sender"` (WA's CDN
  asked the sender's device to reupload — caller should not poll, the
  sender's app must open). New `MediaUnavailable*` const block documents
  the full enum.
- `Message.Status` enum gains `"read"` and `"sent"`. New `MessageStatus*`
  const block documents the full enum.

### Changed (semver-minor on the type level, but additive at runtime)

- `Group` nullable fields migrated to pointers to match the v1.0.1
  policy (preserve null vs empty distinction):
  - `Name`, `Description`, `Owner`, `CreatedAt`: `string` → `*string`
  - `ParticipantCount`: `int` → `*int`
  - `Announce`, `Restrict`: `bool` → `*bool`

  Previously these collapsed `null`, `""`, and `0` into the same zero
  value, hiding the wire distinction. Code that read these fields
  directly will need to deref / nil-check; code that round-tripped via
  JSON keeps working unchanged.

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
