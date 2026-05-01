package blueticks

// Response and request types for the Blueticks API.
//
// Field naming follows the wire format (snake_case) via `json:` tags. Go
// field names use Go convention (CamelCase). Nullable / optional fields use
// pointers OR `omitempty` — pointers preserve "null vs absent" semantics for
// fields where that matters (most timestamps, optional scalars), while
// `omitempty` is used for request-side optional inputs to avoid sending a
// zero value the caller never set.
//
// Unknown fields are ignored at decode time (the default for encoding/json),
// matching Python `extra="ignore"` and Zod `.strip()` — newer backend fields
// don't break older SDK consumers.

// Page is the cursor-paginated list envelope returned by every v1 list
// endpoint. Iterate Data for the current page; pass NextCursor back as the
// Cursor option of the next call to continue. When HasMore is false,
// NextCursor is empty and the iteration is complete.
//
// Go has no generics-aware envelope shared across resources without forcing
// callers into reflection — instead each resource defines its own typed
// page (e.g., MessagesPage). This keeps decoding direct and IDE-friendly.

// --- Pagination wrappers (one per resource) ----------------------------------

// MessagesPage is the paginated response from MessagesResource.List.
type MessagesPage struct {
	Data       []Message `json:"data"`
	HasMore    bool      `json:"has_more"`
	NextCursor *string   `json:"next_cursor,omitempty"`
}

// AudiencesPage is the paginated response from AudiencesResource.List.
type AudiencesPage struct {
	Data       []Audience `json:"data"`
	HasMore    bool       `json:"has_more"`
	NextCursor *string    `json:"next_cursor,omitempty"`
}

// CampaignsPage is the paginated response from CampaignsResource.List.
type CampaignsPage struct {
	Data       []Campaign `json:"data"`
	HasMore    bool       `json:"has_more"`
	NextCursor *string    `json:"next_cursor,omitempty"`
}

// WebhooksPage is the paginated response from WebhooksResource.List.
type WebhooksPage struct {
	Data       []Webhook `json:"data"`
	HasMore    bool      `json:"has_more"`
	NextCursor *string   `json:"next_cursor,omitempty"`
}

// ChatsPage is the paginated response from ChatsResource.List.
type ChatsPage struct {
	Data       []Chat  `json:"data"`
	HasMore    bool    `json:"has_more"`
	NextCursor *string `json:"next_cursor,omitempty"`
}

// ChatMessagesPage is the paginated response from ChatsResource.ListMessages.
type ChatMessagesPage struct {
	Data       []ChatMessage `json:"data"`
	HasMore    bool          `json:"has_more"`
	NextCursor *string       `json:"next_cursor,omitempty"`
}

// ParticipantsPage is the paginated response from ChatsResource.ListParticipants.
type ParticipantsPage struct {
	Data       []Participant `json:"data"`
	HasMore    bool          `json:"has_more"`
	NextCursor *string       `json:"next_cursor,omitempty"`
}

// ContactsPage is the paginated response from ContactsResource.List.
type ContactsPage struct {
	Data       []Contact `json:"data"`
	HasMore    bool      `json:"has_more"`
	NextCursor *string   `json:"next_cursor,omitempty"`
}

// ScheduledMessagesPage is the paginated response from
// ScheduledMessagesResource.List.
type ScheduledMessagesPage struct {
	Data       []ScheduledMessage `json:"data"`
	HasMore    bool               `json:"has_more"`
	NextCursor *string            `json:"next_cursor,omitempty"`
}

// --- Account / Ping ----------------------------------------------------------

// Account is the response from GET /v1/account.
//
// Timezone is a pointer because the OpenAPI declares it `string|null`;
// nil means the server returned `null`, not just an empty string.
type Account struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Timezone  *string `json:"timezone,omitempty"`
	CreatedAt string  `json:"created_at"`
}

// Ping is the response from GET /v1/ping.
type Ping struct {
	AccountID string   `json:"account_id"`
	KeyPrefix string   `json:"key_prefix"`
	Scopes    []string `json:"scopes"`
}

// --- Messages ----------------------------------------------------------------

// Documented values for Message.Status. The wire type is `string` (the OpenAPI
// schema enumerates these values but Go has no dedicated enum type) — compare
// the field directly against these constants.
const (
	MessageStatusScheduled = "scheduled"
	MessageStatusQueued    = "queued"
	MessageStatusSending   = "sending"
	MessageStatusSent      = "sent"
	MessageStatusDelivered = "delivered"
	MessageStatusRead      = "read"
	MessageStatusFailed    = "failed"
)

// Message is a message sent via the Blueticks API.
//
// Status: scheduled | queued | sending | sent | delivered | read | failed.
// Compare against the MessageStatus* constants above.
type Message struct {
	ID            string `json:"id"`
	To            string `json:"to"`
	From          string `json:"from,omitempty"`
	Text          string `json:"text,omitempty"`
	MediaURL      string `json:"media_url,omitempty"`
	Status        string `json:"status"`
	SendAt        string `json:"send_at,omitempty"`
	CreatedAt     string `json:"created_at"`
	SentAt        string `json:"sent_at,omitempty"`
	DeliveredAt   string `json:"delivered_at,omitempty"`
	ReadAt        string `json:"read_at,omitempty"`
	FailedAt      string `json:"failed_at,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
}

// SendMessageParams is the request body for MessagesResource.Send.
type SendMessageParams struct {
	To             string `json:"to"`
	Text           string `json:"text,omitempty"`
	MediaURL       string `json:"media_url,omitempty"`
	MediaCaption   string `json:"media_caption,omitempty"`
	SendAt         string `json:"send_at,omitempty"`
	From           string `json:"from,omitempty"`
	IdempotencyKey string `json:"-"`
}

// --- Audiences ---------------------------------------------------------------

// Audience is an audience (recipient list) used for campaigns.
type Audience struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	ContactCount int               `json:"contact_count"`
	CreatedAt    string            `json:"created_at"`
	Contacts     []AudienceContact `json:"contacts,omitempty"`
	Page         int               `json:"page,omitempty"`
	HasMore      bool              `json:"has_more,omitempty"`
}

// AudienceContact is a single contact within an audience.
type AudienceContact struct {
	ID        string            `json:"id"`
	To        string            `json:"to"`
	Variables map[string]string `json:"variables,omitempty"`
	AddedAt   string            `json:"added_at"`
}

// AppendContactsResult is the response from AudiencesResource.AppendContacts.
type AppendContactsResult struct {
	Added        int `json:"added"`
	ContactCount int `json:"contact_count"`
}

// CreateAudienceParams is the body for AudiencesResource.Create.
type CreateAudienceParams struct {
	Name     string                 `json:"name"`
	Contacts []AudienceContactInput `json:"contacts,omitempty"`
}

// AudienceContactInput is the contact shape accepted on create / append.
type AudienceContactInput struct {
	To        string            `json:"to"`
	Variables map[string]string `json:"variables,omitempty"`
}

// UpdateAudienceParams is the body for AudiencesResource.Update.
type UpdateAudienceParams struct {
	Name string `json:"name"`
}

// UpdateAudienceContactParams is the body for AudiencesResource.UpdateContact.
type UpdateAudienceContactParams struct {
	To        string            `json:"to,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
}

// --- Campaigns ---------------------------------------------------------------

// Campaign is a broadcast campaign sending the same message to every contact
// in an audience.
//
// Status: pending | running | paused | complete_sent | complete_delivered | aborted.
type Campaign struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	AudienceID     string `json:"audience_id"`
	Status         string `json:"status"`
	TotalCount     int    `json:"total_count"`
	SentCount      int    `json:"sent_count"`
	DeliveredCount int    `json:"delivered_count"`
	ReadCount      int    `json:"read_count"`
	FailedCount    int    `json:"failed_count"`
	From           string `json:"from,omitempty"`
	CreatedAt      string `json:"created_at"`
	StartedAt      string `json:"started_at,omitempty"`
	CompletedAt    string `json:"completed_at,omitempty"`
	AbortedAt      string `json:"aborted_at,omitempty"`
}

// CreateCampaignParams is the body for CampaignsResource.Create.
type CreateCampaignParams struct {
	Name              string `json:"name"`
	AudienceID        string `json:"audience_id"`
	Text              string `json:"text,omitempty"`
	MediaURL          string `json:"media_url,omitempty"`
	MediaCaption      string `json:"media_caption,omitempty"`
	From              string `json:"from,omitempty"`
	OnMissingVariable string `json:"on_missing_variable,omitempty"`
}

// --- Webhooks ----------------------------------------------------------------

// Webhook is a registered webhook subscription.
type Webhook struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Events      []string `json:"events"`
	Description string   `json:"description,omitempty"`
	Status      string   `json:"status"`
	CreatedAt   string   `json:"created_at"`
}

// WebhookCreateResult is returned from Create / RotateSecret. It is a
// Webhook plus the (never-shown-again) Secret used to verify deliveries.
type WebhookCreateResult struct {
	Webhook
	Secret string `json:"secret"`
}

// WebhookEvent is the parsed body of a webhook delivery.
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	CreatedAt string                 `json:"created_at"`
	Data      map[string]interface{} `json:"data"`
}

// CreateWebhookParams is the body for WebhooksResource.Create.
type CreateWebhookParams struct {
	URL         string   `json:"url"`
	Events      []string `json:"events"`
	Description string   `json:"description,omitempty"`
}

// UpdateWebhookParams is the body for WebhooksResource.Update.
type UpdateWebhookParams struct {
	URL         string   `json:"url,omitempty"`
	Events      []string `json:"events,omitempty"`
	Description string   `json:"description,omitempty"`
	Status      string   `json:"status,omitempty"`
}

// --- Chats / Engine-sourced --------------------------------------------------

// Chat is a WhatsApp conversation.
//
// Pointer fields preserve the wire-level null vs. empty distinction
// (the OpenAPI declares them `string|null` / `integer|null`). nil ==
// "the server returned null"; a non-nil empty string means the field
// was present but blank.
type Chat struct {
	ID            string  `json:"id"`
	Name          *string `json:"name,omitempty"`
	IsGroup       bool    `json:"is_group"`
	LastMessageAt *string `json:"last_message_at,omitempty"`
	UnreadCount   *int    `json:"unread_count,omitempty"`
}

// Participant is a member of a group chat.
type Participant struct {
	ChatID       string `json:"chat_id"`
	IsAdmin      bool   `json:"is_admin"`
	IsSuperAdmin bool   `json:"is_super_admin,omitempty"`
}

// ChatMessage is a message returned by GET /v1/chats/{id}/messages.
type ChatMessage struct {
	Key       string `json:"key"`
	ChatID    string `json:"chat_id"`
	From      string `json:"from"`
	Timestamp string `json:"timestamp,omitempty"`
	Text      string `json:"text,omitempty"`
	Type      string `json:"type"`
	FromMe    bool   `json:"from_me"`
	Ack       int    `json:"ack,omitempty"`
	MediaURL  string `json:"media_url,omitempty"`
	Caption   string `json:"caption,omitempty"`
	Filename  string `json:"filename,omitempty"`
}

// Documented values for ChatMedia.MediaUnavailable. The wire type is `string`
// (the OpenAPI schema enumerates these values but Go has no dedicated enum
// type) — compare the field directly against these constants. An empty string
// means the fetch succeeded (or no fetch was needed).
const (
	MediaUnavailableExpired        = "expired"
	MediaUnavailableFetching       = "fetching"
	MediaUnavailableAwaitingSender = "awaiting_sender"
	MediaUnavailableError          = "error"
	MediaUnavailableNoMedia        = "no_media"
)

// ChatMedia is the response from GET /v1/chats/{id}/messages/{key}/media.
//
// MediaUnavailable is one of: expired | fetching | awaiting_sender | error |
// no_media (or empty on success). Compare against the MediaUnavailable*
// constants above.
type ChatMedia struct {
	URL              string `json:"url,omitempty"`
	Mimetype         string `json:"mimetype,omitempty"`
	Filename         string `json:"filename,omitempty"`
	DataBase64       string `json:"data_base64,omitempty"`
	OriginalQuality  *bool  `json:"original_quality,omitempty"`
	MediaUnavailable string `json:"media_unavailable,omitempty"`
}

// ListChatMessagesParams is the query-param bundle for ChatsResource.ListMessages.
type ListChatMessagesParams struct {
	Mode         string // "latest" or "history"; default "latest"
	Query        string
	Since        string
	Until        string
	MessageTypes []string // joined comma-separated on the wire
	Limit        int
	Cursor       string
}

// OkResponse is the `{"ok": true}` envelope returned by chat write endpoints
// (mark_read, reactions) that don't have a richer payload to surface.
type OkResponse struct {
	Ok bool `json:"ok"`
}

// ChatRef is the `{"chat_id": "..."}` envelope returned by ChatsResource.Open.
type ChatRef struct {
	ChatID string `json:"chat_id"`
}

// MessageAck is the response from GET /v1/chats/{chat_id}/messages/{key}/ack.
//
// Ack semantics: -1=error, 0=pending, 1=server, 2=device, 3=read, 4=played.
// nil means no engine response was available.
type MessageAck struct {
	Ack *int `json:"ack"`
}

// LoadOlderMessagesResponse is the response from
// POST /v1/chats/{chat_id}/messages/load_older.
//
// TotalMessages and Added are pointers because the OpenAPI declares them
// nullable; CanLoadMore is required.
type LoadOlderMessagesResponse struct {
	TotalMessages *int `json:"total_messages"`
	Added         *int `json:"added"`
	CanLoadMore   bool `json:"can_load_more"`
}

// MediaURLResponse is the response from
// GET /v1/chats/{chat_id}/messages/{key}/media_url.
//
// URL is a pointer because the OpenAPI declares it nullable; nil means the
// engine could not produce a hosted URL.
type MediaURLResponse struct {
	URL *string `json:"url"`
}

// BatchMessageAcksParams is the body for ChatsResource.BatchMessageAcks.
type BatchMessageAcksParams struct {
	MessageKeys []string `json:"message_keys"`
}

// BatchMessageAcksResponse is the response from POST /v1/chats/message_acks.
//
// Each Data entry is a free-form per-message ack envelope (the server permits
// additional properties, so the SDK keeps the entry as a map for forward
// compatibility).
type BatchMessageAcksResponse struct {
	Data []map[string]interface{} `json:"data"`
}

// --- Contacts ----------------------------------------------------------------

// Contact is a WhatsApp contact in the user's address book.
type Contact struct {
	ChatID     string `json:"chat_id"`
	Name       string `json:"name,omitempty"`
	IsBusiness bool   `json:"is_business"`
}

// ProfilePicture is the response from GET /v1/contacts/{chat_id}/profile_picture.
type ProfilePicture struct {
	URL string `json:"url,omitempty"`
}

// --- Groups ------------------------------------------------------------------

// Group is a WhatsApp group's metadata.
//
// Pointer fields preserve the wire-level null vs. empty distinction (the
// OpenAPI declares them nullable). nil == "the server returned null"; a
// non-nil zero value means the field was present but blank / zero.
type Group struct {
	ID               string             `json:"id"`
	Name             *string            `json:"name,omitempty"`
	Description      *string            `json:"description,omitempty"`
	Owner            *string            `json:"owner,omitempty"`
	CreatedAt        *string            `json:"created_at,omitempty"`
	LastMessageAt    *string            `json:"last_message_at,omitempty"`
	ParticipantCount *int               `json:"participant_count,omitempty"`
	Announce         *bool              `json:"announce,omitempty"`
	Restrict         *bool              `json:"restrict,omitempty"`
	Participants     []GroupParticipant `json:"participants,omitempty"`
}

// GroupParticipant is a single member of a Group as returned in
// Group.Participants.
type GroupParticipant struct {
	ChatID       string  `json:"chat_id"`
	IsAdmin      bool    `json:"is_admin"`
	IsSuperAdmin bool    `json:"is_super_admin"`
	Name         *string `json:"name,omitempty"`
}

// CreateGroupParams is the body for GroupsResource.Create.
type CreateGroupParams struct {
	Name         string   `json:"name"`
	Participants []string `json:"participants"`
}

// UpdateGroupParams is the body for GroupsResource.Update.
type UpdateGroupParams struct {
	Name     string          `json:"name,omitempty"`
	Settings map[string]bool `json:"settings,omitempty"`
}

// --- Scheduled messages ------------------------------------------------------

// ScheduledMessage is a one-time or recurring scheduled message returned by
// /v1/scheduled-messages.
//
// All non-id fields are nullable per the OpenAPI schema — pointers preserve
// the wire-level null vs. empty distinction. Status mirrors Message.Status
// (compare against the MessageStatus* constants).
type ScheduledMessage struct {
	ID             string  `json:"id"`
	To             *string `json:"to"`
	Text           *string `json:"text"`
	MediaURL       *string `json:"media_url"`
	MediaCaption   *string `json:"media_caption"`
	MediaFilename  *string `json:"media_filename"`
	MediaMimeType  *string `json:"media_mime_type"`
	SendAt         *string `json:"send_at"`
	Status         *string `json:"status"`
	IsRecurring    bool    `json:"is_recurring"`
	RecurrenceRule *string `json:"recurrence_rule"`
	CreatedAt      *string `json:"created_at"`
	UpdatedAt      *string `json:"updated_at"`
}

// ListScheduledMessagesParams is the query-param bundle for
// ScheduledMessagesResource.List.
type ListScheduledMessagesParams struct {
	Limit  int
	Cursor string
}

// UpdateScheduledMessageParams is the body for ScheduledMessagesResource.Update.
//
// All fields are optional; omit (leave nil / zero) any you don't want to
// change.
type UpdateScheduledMessageParams struct {
	Text         *string `json:"text,omitempty"`
	MediaURL     *string `json:"media_url,omitempty"`
	MediaCaption *string `json:"media_caption,omitempty"`
	SendAt       *string `json:"send_at,omitempty"`
}

// --- Engines -----------------------------------------------------------------

// EngineStatus is the response from GET /v1/engines.
type EngineStatus struct {
	Connected bool   `json:"connected"`
	State     string `json:"state,omitempty"`
	Stream    string `json:"stream,omitempty"`
	HasSynced *bool  `json:"has_synced,omitempty"`
}

// WhatsAppMe is the authenticated WhatsApp account profile from
// GET /v1/engines/me.
type WhatsAppMe struct {
	Phone    string `json:"phone,omitempty"`
	Name     string `json:"name,omitempty"`
	Platform string `json:"platform,omitempty"`
}

// --- Utils -------------------------------------------------------------------

// PhoneValidation is the response from POST /v1/utils/validate_phone.
type PhoneValidation struct {
	Valid           bool   `json:"valid"`
	FormattedChatID string `json:"formatted_chat_id,omitempty"`
}

// LinkPreview is the response from GET /v1/utils/link_preview.
type LinkPreview struct {
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	Thumbnail    string `json:"thumbnail,omitempty"`
	CanonicalURL string `json:"canonical_url,omitempty"`
}

// --- ID-only delete envelope -------------------------------------------------

// DeletedRef is the typical {"id": "..."} envelope returned by DELETE endpoints.
type DeletedRef struct {
	ID string `json:"id"`
}
