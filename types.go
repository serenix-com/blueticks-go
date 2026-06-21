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
// page (e.g., ScheduledMessagesPage). This keeps decoding direct and IDE-friendly.

// --- Pagination wrappers (one per resource) ----------------------------------

// AudiencesPage is the paginated response from AudiencesResource.List.
type AudiencesPage struct {
	Data       []Audience `json:"data"`
	HasMore    bool       `json:"hasMore"`
	NextCursor *string    `json:"nextCursor,omitempty"`
}

// CampaignsPage is the paginated response from CampaignsResource.List.
type CampaignsPage struct {
	Data       []Campaign `json:"data"`
	HasMore    bool       `json:"hasMore"`
	NextCursor *string    `json:"nextCursor,omitempty"`
}

// WebhooksPage is the paginated response from WebhooksResource.List.
type WebhooksPage struct {
	Data       []Webhook `json:"data"`
	HasMore    bool      `json:"hasMore"`
	NextCursor *string   `json:"nextCursor,omitempty"`
}

// ChatsPage is the paginated response from ChatsResource.List.
type ChatsPage struct {
	Data       []Chat  `json:"data"`
	HasMore    bool    `json:"hasMore"`
	NextCursor *string `json:"nextCursor,omitempty"`
}

// ChatMessagesPage is the paginated response from ChatsResource.ListMessages.
type ChatMessagesPage struct {
	Data       []ChatMessage `json:"data"`
	HasMore    bool          `json:"hasMore"`
	NextCursor *string       `json:"nextCursor,omitempty"`
}

// ParticipantsPage is the paginated response from ChatsResource.ListParticipants.
type ParticipantsPage struct {
	Data       []Participant `json:"data"`
	HasMore    bool          `json:"hasMore"`
	NextCursor *string       `json:"nextCursor,omitempty"`
}

// ContactsPage is the paginated response from ContactsResource.List.
type ContactsPage struct {
	Data       []Contact `json:"data"`
	HasMore    bool      `json:"hasMore"`
	NextCursor *string   `json:"nextCursor,omitempty"`
}

// GroupsPage is the paginated response from GroupsResource.List.
type GroupsPage struct {
	Data       []Group `json:"data"`
	HasMore    bool    `json:"hasMore"`
	NextCursor *string `json:"nextCursor,omitempty"`
}

// NewslettersPage is the paginated response from NewslettersResource.List.
type NewslettersPage struct {
	Data       []Newsletter `json:"data"`
	HasMore    bool         `json:"hasMore"`
	NextCursor *string      `json:"nextCursor,omitempty"`
}

// --- Ping --------------------------------------------------------------------

// Ping is the response from GET /v1/ping.
type Ping struct {
	AccountID string   `json:"accountId"`
	KeyPrefix string   `json:"keyPrefix"`
	Scopes    []string `json:"scopes"`
}

// --- Account -----------------------------------------------------------------

// Account is the response from GET /v1/account.
//
// Timezone is a pointer because the OpenAPI declares it `string|null`;
// nil means the server returned `null`, not just an empty string.
type Account struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Timezone  *string `json:"timezone,omitempty"`
	CreatedAt string  `json:"createdAt"`
}

// --- Newsletters -------------------------------------------------------------

// Documented values for Newsletter.Verification. The wire type is `string`.
const (
	NewsletterVerificationVerified   = "VERIFIED"
	NewsletterVerificationUnverified = "UNVERIFIED"
)

// Newsletter is the response from POST /v1/newsletters, GET /v1/newsletters,
// and GET /v1/newsletters/{id}.
//
// Pointer fields (Description, Owner, CreatedAt, Subscribers, Invite,
// Verification) are nullable on the wire — nil means the server returned null.
type Newsletter struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	Owner        *string `json:"owner"`
	CreatedAt    *string `json:"createdAt"`
	Subscribers  *int    `json:"subscribers"`
	Invite       *string `json:"invite"`
	Verification *string `json:"verification"`
}

// CreateNewsletterParams is the request body for NewslettersResource.Create.
//
// Name is required (1–100 chars). Description is optional (max 2048 chars).
type CreateNewsletterParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// ListNewslettersParams is the query-param bundle for NewslettersResource.List.
type ListNewslettersParams struct {
	Limit  int
	Cursor string
}

// --- Audiences ---------------------------------------------------------------

// Audience is an audience (recipient list) used for campaigns.
type Audience struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	ContactCount int               `json:"contactCount"`
	CreatedAt    string            `json:"createdAt"`
	Contacts     []AudienceContact `json:"contacts,omitempty"`
	Page         int               `json:"page,omitempty"`
	HasMore      bool              `json:"hasMore,omitempty"`
}

// AudienceContact is a single contact within an audience.
type AudienceContact struct {
	ID        string            `json:"id"`
	To        string            `json:"to"`
	Variables map[string]string `json:"variables,omitempty"`
	AddedAt   string            `json:"addedAt"`
}

// AppendContactsResult is the response from AudiencesResource.AppendContacts.
type AppendContactsResult struct {
	Added        int `json:"added"`
	ContactCount int `json:"contactCount"`
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
	AudienceID     string `json:"audienceId"`
	Status         string `json:"status"`
	TotalCount     int    `json:"totalCount"`
	SentCount      int    `json:"sentCount"`
	DeliveredCount int    `json:"deliveredCount"`
	ReadCount      int    `json:"readCount"`
	FailedCount    int    `json:"failedCount"`
	From           string `json:"from,omitempty"`
	CreatedAt      string `json:"createdAt"`
	StartedAt      string `json:"startedAt,omitempty"`
	CompletedAt    string `json:"completedAt,omitempty"`
	AbortedAt      string `json:"abortedAt,omitempty"`
}

// CreateCampaignParams is the body for CampaignsResource.Create.
type CreateCampaignParams struct {
	Name              string `json:"name"`
	AudienceID        string `json:"audienceId"`
	Text              string `json:"text,omitempty"`
	MediaURL          string `json:"mediaUrl,omitempty"`
	MediaCaption      string `json:"mediaCaption,omitempty"`
	From              string `json:"from,omitempty"`
	OnMissingVariable string `json:"onMissingVariable,omitempty"`
}

// --- Webhooks ----------------------------------------------------------------

// Webhook is a registered webhook subscription.
type Webhook struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Events      []string `json:"events"`
	Description string   `json:"description,omitempty"`
	Status      string   `json:"status"`
	CreatedAt   string   `json:"createdAt"`
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
	IsGroup       bool    `json:"isGroup"`
	IsNewsletter  bool    `json:"isNewsletter"`
	LastMessageAt *string `json:"lastMessageAt,omitempty"`
	UnreadCount   *int    `json:"unreadCount,omitempty"`
	MarkedUnread  bool    `json:"markedUnread"`
}

// Participant is a member of a group chat.
type Participant struct {
	ChatID       string `json:"chatId"`
	IsAdmin      bool   `json:"isAdmin"`
	IsSuperAdmin bool   `json:"isSuperAdmin,omitempty"`
}

// ChatMessage is a message returned by GET /v1/chats/{id}/messages.
type ChatMessage struct {
	Key       string `json:"key"`
	ChatID    string `json:"chatId"`
	From      string `json:"from"`
	Timestamp string `json:"timestamp,omitempty"`
	Text      string `json:"text,omitempty"`
	Type      string `json:"type"`
	FromMe    bool   `json:"fromMe"`
	Ack       int    `json:"ack,omitempty"`
	MediaURL  string `json:"mediaUrl,omitempty"`
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

// ChatMedia is the response from GET /v1/messages/media/{waMessageKey}.
//
// MediaUnavailable is one of: expired | fetching | awaiting_sender | error |
// no_media (or empty on success). Compare against the MediaUnavailable*
// constants above.
type ChatMedia struct {
	URL              string `json:"url,omitempty"`
	Mimetype         string `json:"mimetype,omitempty"`
	Filename         string `json:"filename,omitempty"`
	DataBase64       string `json:"dataBase64,omitempty"`
	OriginalQuality  *bool  `json:"originalQuality,omitempty"`
	MediaUnavailable string `json:"mediaUnavailable,omitempty"`
}

// ListChatMessagesParams is the query-param bundle for ChatsResource.ListMessages.
type ListChatMessagesParams struct {
	Order        string // "asc" (oldest-first) or "desc" (newest-first, default)
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

// ChatRef is the `{"chatId": "..."}` envelope returned by ChatsResource.Open.
type ChatRef struct {
	ChatID string `json:"chatId"`
}

// MessageAck is the response from GET /v1/messages/ack/{waMessageKey}.
//
// Ack semantics: -1=error, 0=pending, 1=server, 2=device, 3=read, 4=played.
// nil means no engine response was available.
type MessageAck struct {
	Ack *int `json:"ack"`
}

// LoadOlderMessagesResponse is the response from
// POST /v1/messages/load_older/{chatId}.
//
// TotalMessages and Added are pointers because the OpenAPI declares them
// nullable; CanLoadMore is required.
type LoadOlderMessagesResponse struct {
	TotalMessages *int `json:"totalMessages"`
	Added         *int `json:"added"`
	CanLoadMore   bool `json:"canLoadMore"`
}

// MediaURLResponse is the response from
// GET /v1/messages/media_url/{waMessageKey}.
//
// URL is a pointer because the OpenAPI declares it nullable; nil means the
// engine could not produce a hosted URL.
type MediaURLResponse struct {
	URL *string `json:"url"`
}

// BatchMessageAcksParams is the body for ChatsResource.BatchMessageAcks.
type BatchMessageAcksParams struct {
	MessageKeys []string `json:"messageKeys"`
}

// BatchMessageAckEntry is a single per-message ack envelope inside
// BatchMessageAcksResponse.Data.
//
// Ack semantics mirror MessageAck.Ack: -1=error, 0=pending, 1=server,
// 2=device, 3=read, 4=played; nil when no engine response was available.
type BatchMessageAckEntry struct {
	Key string `json:"key"`
	Ack *int   `json:"ack"`
}

// BatchMessageAcksResponse is the response from POST /v1/messages/acks.
type BatchMessageAcksResponse struct {
	Data []BatchMessageAckEntry `json:"data"`
}

// SendMessageInChatParams is the request body for ChatsResource.SendMessage
// (POST /v1/messages/{chatId}).
//
// The wire schema is a discriminated union on the `type` field — same shape
// as CreateScheduledMessageParams minus `to` (derived from the URL path) and
// `sendAt` (this endpoint is fire-and-forget):
//
//	params := blueticks.SendMessageInChatParams{
//	    Type: blueticks.MessageTypeText,
//	    Text: "hello",
//	}
//
// The server validates the type/variant pair (Text for "text", Media for
// "media", Poll for "poll").
type SendMessageInChatParams struct {
	// Type is the discriminator: "text", "media", or "poll".
	Type string `json:"type"`
	// Text is required when Type == "text".
	Text string `json:"text,omitempty"`
	// LinkPreview is optional when Type == "text".
	// true (default) = auto-fetch; false = never attach; object = custom card.
	LinkPreview interface{} `json:"linkPreview,omitempty"`
	// Media is required when Type == "media".
	Media *SendMediaObject `json:"media,omitempty"`
	// Poll is required when Type == "poll".
	Poll    *SendPollObject `json:"poll,omitempty"`
	From    string          `json:"from,omitempty"`
	ReplyTo string          `json:"replyTo,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header, not in the body.
	IdempotencyKey string `json:"-"`
}

// --- Contacts ----------------------------------------------------------------

// Contact is a WhatsApp contact in the user's address book.
type Contact struct {
	ChatID     string `json:"chatId"`
	Name       string `json:"name,omitempty"`
	IsBusiness bool   `json:"isBusiness"`
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
	CreatedAt        *string            `json:"createdAt,omitempty"`
	LastMessageAt    *string            `json:"lastMessageAt,omitempty"`
	ParticipantCount *int               `json:"participantCount,omitempty"`
	Announce         *bool              `json:"announce,omitempty"`
	Restrict         *bool              `json:"restrict,omitempty"`
	Participants     []GroupParticipant `json:"participants,omitempty"`
}

// GroupParticipant is a single member of a Group as returned in
// Group.Participants.
type GroupParticipant struct {
	ChatID       string  `json:"chatId"`
	IsAdmin      bool    `json:"isAdmin"`
	IsSuperAdmin bool    `json:"isSuperAdmin"`
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

// --- Scheduled Messages ------------------------------------------------------

// Documented values for ScheduledMessage.Status. The wire type is `string`
// (the OpenAPI schema enumerates these values but Go has no dedicated enum
// type) — compare the field directly against these constants.
const (
	MessageStatusPending   = "pending"
	MessageStatusConfirmed = "confirmed"
	MessageStatusReceived  = "received"
	MessageStatusRead      = "read"
	MessageStatusPlayed    = "played"
	MessageStatusFailed    = "failed"
)

// Documented values for ScheduledMessage.Type. The wire type is `string`.
const (
	MessageTypeText  = "text"
	MessageTypeMedia = "media"
	MessageTypePoll  = "poll"
)

// ScheduledMessagesPage is the paginated response from
// ScheduledMessagesResource.List.
type ScheduledMessagesPage struct {
	Data       []ScheduledMessage `json:"data"`
	HasMore    bool               `json:"hasMore"`
	NextCursor *string            `json:"nextCursor,omitempty"`
}

// MessageLinkPreview is the optional rich-preview card embedded in a
// ScheduledMessage.
//
// All fields are nullable on the wire — nil means the server returned null.
type MessageLinkPreview struct {
	Title        *string `json:"title"`
	Description  *string `json:"description"`
	CanonicalURL *string `json:"canonicalUrl"`
	Thumbnail    *string `json:"thumbnail"`
}

// ScheduledMessage is a message created via POST /v1/scheduled-messages and
// returned by all /v1/scheduled-messages* endpoints.
//
// Type: text | media | poll. Compare against the MessageType* constants.
// Status: pending | confirmed | received | read | played | failed. Compare
// against the MessageStatus* constants.
//
// Nullable fields use pointer types — nil means the server returned null
// (not absent).
type ScheduledMessage struct {
	ID            string              `json:"id"`
	Key           *string             `json:"key"`
	To            string              `json:"to"`
	From          *string             `json:"from"`
	Type          string              `json:"type"`
	Text          *string             `json:"text"`
	MediaURL      *string             `json:"mediaUrl"`
	MediaKind     *string             `json:"mediaKind"`
	PollQuestion  *string             `json:"pollQuestion"`
	Status        string              `json:"status"`
	SendAt        *string             `json:"sendAt"`
	CreatedAt     string              `json:"createdAt"`
	ConfirmedAt   *string             `json:"confirmedAt"`
	ReceivedAt    *string             `json:"receivedAt"`
	ReadAt        *string             `json:"readAt"`
	PlayedAt      *string             `json:"playedAt"`
	FailedAt      *string             `json:"failedAt"`
	FailureReason *string             `json:"failureReason"`
	LinkPreview   *MessageLinkPreview `json:"linkPreview"`
}

// SendMediaObject is the nested media descriptor inside
// CreateScheduledMessageParams when Type == "media".
type SendMediaObject struct {
	URL      string `json:"url,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename,omitempty"`
}

// SendPollObject is the nested poll descriptor inside
// CreateScheduledMessageParams when Type == "poll".
type SendPollObject struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	AllowMultiple bool     `json:"allowMultiple,omitempty"`
}

// CreateScheduledMessageParams is the request body for
// ScheduledMessagesResource.Create.
//
// The POST body is a discriminated union on the `type` field:
//
//	params := blueticks.CreateScheduledMessageParams{
//	    Type: blueticks.MessageTypeText,
//	    To:   "+15551234567",
//	    Text: "hello",
//	}
//
// The server validates the type/variant pair (Text for "text", Media for
// "media", Poll for "poll").
type CreateScheduledMessageParams struct {
	// Type is the discriminator: "text", "media", or "poll".
	Type string `json:"type"`
	To   string `json:"to"`
	// Text is required when Type == "text".
	Text string `json:"text,omitempty"`
	// LinkPreview is optional when Type == "text".
	// true (default) = auto-fetch; false = never attach; object = custom card.
	LinkPreview interface{} `json:"linkPreview,omitempty"`
	// Media is required when Type == "media".
	Media *SendMediaObject `json:"media,omitempty"`
	// Poll is required when Type == "poll".
	Poll    *SendPollObject `json:"poll,omitempty"`
	SendAt  string          `json:"sendAt,omitempty"`
	From    string          `json:"from,omitempty"`
	ReplyTo string          `json:"replyTo,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header, not in the body.
	IdempotencyKey string `json:"-"`
}

// ListScheduledMessagesParams is the query-param bundle for
// ScheduledMessagesResource.List.
type ListScheduledMessagesParams struct {
	Limit  int
	Cursor string
	ChatID string
	Status string
	Q      string
}

// UpdateScheduledMessageParams is the request body for
// ScheduledMessagesResource.Update (PATCH /v1/scheduled-messages/{id}). All
// fields are optional but at least one must be provided. Replaces the
// matching field on the pending message; rejected once the message has
// advanced past the editable window.
type UpdateScheduledMessageParams struct {
	Text         *string `json:"text,omitempty"`
	MediaURL     *string `json:"mediaUrl,omitempty"`
	MediaCaption *string `json:"mediaCaption,omitempty"`
	SendAt       *string `json:"sendAt,omitempty"`
}

// --- Engines -----------------------------------------------------------------

// EngineStatus is the response from GET /v1/engines.
type EngineStatus struct {
	Connected bool   `json:"connected"`
	State     string `json:"state,omitempty"`
	Stream    string `json:"stream,omitempty"`
	HasSynced *bool  `json:"hasSynced,omitempty"`
}

// --- ID-only delete envelope -------------------------------------------------

// DeletedRef is the typical {"id": "..."} envelope returned by DELETE endpoints.
type DeletedRef struct {
	ID string `json:"id"`
}

// DeletedResource is the {"id": "...", "deleted": true} envelope returned by
// soft-delete endpoints (audiences, scheduled-messages, webhooks). The
// `deleted` flag is always true on a successful response — its presence
// distinguishes a delete confirmation from a regular resource fetch.
type DeletedResource struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}
