package blueticks

// Response and request types for the Blueticks API (v2).
//
// Field naming follows the wire format (camelCase) via `json:` tags. Go field
// names use Go convention (PascalCase, idiomatic acronyms — URL, ID). Nullable
// wire fields use pointer types so a decoded nil is distinguishable from a
// present zero value; optional request inputs use pointers or `omitempty` so
// the encoder omits values the caller never set.
//
// Every 2xx response is wrapped in a `{success, data}` envelope. Single-object
// endpoints unwrap `data` in-method (a small anonymous struct); list endpoints
// return the flat offset-pagination envelope `{data, limit, skip, total}`
// directly into the typed <Resource>Page structs below (the `success` field is
// ignored on decode).
//
// Unknown fields are ignored at decode time (the default for encoding/json),
// so newer backend fields don't break older SDK consumers.

// --- Pagination wrappers (one per resource) ----------------------------------

// AudiencesPage is the offset-paginated response from AudiencesResource.List.
type AudiencesPage struct {
	Data  []Audience `json:"data"`
	Limit int        `json:"limit"`
	Skip  int        `json:"skip"`
	Total int        `json:"total"`
}

// CampaignsPage is the offset-paginated response from CampaignsResource.List.
type CampaignsPage struct {
	Data  []Campaign `json:"data"`
	Limit int        `json:"limit"`
	Skip  int        `json:"skip"`
	Total int        `json:"total"`
}

// ChatsPage is the offset-paginated response from ChatsResource.List.
type ChatsPage struct {
	Data  []Chat `json:"data"`
	Limit int    `json:"limit"`
	Skip  int    `json:"skip"`
	Total int    `json:"total"`
}

// ParticipantsPage is the offset-paginated response from ChatsResource.ListParticipants.
type ParticipantsPage struct {
	Data  []Participant `json:"data"`
	Limit int           `json:"limit"`
	Skip  int           `json:"skip"`
	Total int           `json:"total"`
}

// ContactsPage is the offset-paginated response from ContactsResource.List.
type ContactsPage struct {
	Data  []Contact `json:"data"`
	Limit int       `json:"limit"`
	Skip  int       `json:"skip"`
	Total int       `json:"total"`
}

// GroupsPage is the offset-paginated response from GroupsResource.List.
type GroupsPage struct {
	Data  []GroupListItem `json:"data"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Total int             `json:"total"`
}

// CommonGroupsPage is the offset-paginated response from ContactsResource.CommonGroups.
type CommonGroupsPage struct {
	Data  []GroupListItem `json:"data"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Total int             `json:"total"`
}

// NewslettersPage is the offset-paginated response from NewslettersResource.List.
type NewslettersPage struct {
	Data  []NewsletterListItem `json:"data"`
	Limit int                  `json:"limit"`
	Skip  int                  `json:"skip"`
	Total int                  `json:"total"`
}

// MessagesPage is the offset-paginated response from MessagesResource.List.
// HasMore is an extra hint (true when more messages exist in the engine store
// beyond this page) unique to the message-history endpoint.
type MessagesPage struct {
	Data    []Message `json:"data"`
	Limit   int       `json:"limit"`
	Skip    int       `json:"skip"`
	Total   int       `json:"total"`
	HasMore bool      `json:"hasMore,omitempty"`
}

// PinnedMessagesPage is the offset-paginated response from MessagesResource.ListPinned.
type PinnedMessagesPage struct {
	Data  []PinnedMessage `json:"data"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Total int             `json:"total"`
}

// BatchMessageAcksResult is the response from MessagesResource.BatchAcks. The
// server echoes the offset-pagination fields even though the query is a fixed
// batch of keys.
type BatchMessageAcksResult struct {
	Data  []BatchMessageAckEntry `json:"data"`
	Limit int                    `json:"limit"`
	Skip  int                    `json:"skip"`
	Total int                    `json:"total"`
}

// ScheduledMessagesPage is the offset-paginated response from
// ScheduledMessagesResource.List.
type ScheduledMessagesPage struct {
	Data  []SentMessage `json:"data"`
	Limit int           `json:"limit"`
	Skip  int           `json:"skip"`
	Total int           `json:"total"`
}

// WebhooksPage is the paginated response from WebhooksResource.List.
//
// NOTE: WebhooksResource is hand-maintained (webhooks.go is frozen) and still
// on the 1.x cursor-pagination shape. Kept unchanged so the frozen resource
// compiles and behaves as before.
type WebhooksPage struct {
	Data       []Webhook `json:"data"`
	HasMore    bool      `json:"has_more"`
	NextCursor *string   `json:"next_cursor,omitempty"`
}

// --- Ping --------------------------------------------------------------------

// Documented values for a connection/engine Type. The wire type is `string`.
const (
	EngineTypeGateway = "gateway"
	EngineTypeRegular = "regular"
)

// PingConnection is a single live WhatsApp engine entry in Ping.WhatsappConnections.
type PingConnection struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Connected bool   `json:"connected"`
}

// Ping is the response from GET /v1/ping.
//
// API is always "ok". WhatsappConnections lists the currently-connected engines
// (empty when no WhatsApp is connected); Message is present only in that empty
// state to explain it.
type Ping struct {
	API                 string           `json:"api"`
	AccountID           string           `json:"accountId"`
	WhatsappConnections []PingConnection `json:"whatsappConnections"`
	Message             *string          `json:"message,omitempty"`
}

// --- Account -----------------------------------------------------------------

// Account is the response from GET /v1/account.
//
// UserEmail and Timezone are pointers because the OpenAPI declares them
// `string|null`; nil means the server returned null.
type Account struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	UserEmail *string `json:"userEmail"`
	Timezone  *string `json:"timezone"`
	CreatedAt string  `json:"createdAt"`
}

// --- Audiences ---------------------------------------------------------------

// Audience is a reusable recipient list used as a target for campaigns.
type Audience struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ContactCount int    `json:"contactCount"`
	CreatedAt    string `json:"createdAt"`
}

// AudienceContact is a single contact within an audience (the response from
// AudiencesResource.UpdateContact).
type AudienceContact struct {
	ID        string            `json:"id"`
	To        string            `json:"to"`
	Variables map[string]string `json:"variables"`
	AddedAt   string            `json:"addedAt"`
}

// AppendContactsResult is the response from AudiencesResource.AppendContacts.
type AppendContactsResult struct {
	Added        int `json:"added"`
	ContactCount int `json:"contactCount"`
}

// ContactInput is the contact shape accepted on audience create / append.
//
// Variables are per-contact substitution values usable as {key} tokens in a
// campaign text.
type ContactInput struct {
	To        string            `json:"to"`
	Variables map[string]string `json:"variables,omitempty"`
}

// CreateAudienceParams is the body for AudiencesResource.Create.
type CreateAudienceParams struct {
	Name     string         `json:"name"`
	Contacts []ContactInput `json:"contacts,omitempty"`
}

// UpdateAudienceParams is the body for AudiencesResource.Update.
type UpdateAudienceParams struct {
	Name string `json:"name"`
}

// UpdateAudienceContactParams is the body for AudiencesResource.UpdateContact.
// Both fields are optional; supply at least one.
type UpdateAudienceContactParams struct {
	To        *string           `json:"to,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
}

// --- Campaigns ---------------------------------------------------------------

// Documented values for Campaign.Status. The wire type is `string`.
const (
	CampaignStatusPending           = "pending"
	CampaignStatusRunning           = "running"
	CampaignStatusPaused            = "paused"
	CampaignStatusCompleteSent      = "complete_sent"
	CampaignStatusCompleteDelivered = "complete_delivered"
	CampaignStatusAborted           = "aborted"
)

// Campaign is a broadcast campaign sending the same message to every contact in
// an audience.
//
// Started/Completed/Aborted timestamps are nullable on the wire — nil means the
// server returned null.
type Campaign struct {
	CmpID          string  `json:"cmpId"`
	Name           string  `json:"name"`
	AudienceID     string  `json:"audienceId"`
	Status         string  `json:"status"`
	TotalCount     int     `json:"totalCount"`
	SentCount      int     `json:"sentCount"`
	DeliveredCount int     `json:"deliveredCount"`
	ReadCount      int     `json:"readCount"`
	FailedCount    int     `json:"failedCount"`
	CreatedAt      string  `json:"createdAt"`
	StartedAt      *string `json:"startedAt"`
	CompletedAt    *string `json:"completedAt"`
	AbortedAt      *string `json:"abortedAt"`
}

// Documented values for CreateCampaignParams.OnMissingVariable.
const (
	OnMissingVariableFail = "fail"
	OnMissingVariableSkip = "skip"
)

// CreateCampaignParams is the body for CampaignsResource.Create.
type CreateCampaignParams struct {
	Name              string `json:"name"`
	AudienceID        string `json:"audienceId"`
	Text              string `json:"text,omitempty"`
	MediaURL          string `json:"mediaUrl,omitempty"`
	MediaCaption      string `json:"mediaCaption,omitempty"`
	OnMissingVariable string `json:"onMissingVariable,omitempty"`
}

// --- Webhooks (FROZEN — webhooks.go / webhooks_verify.go are hand-maintained) -
//
// These types back the frozen WebhooksResource and are intentionally left on
// the 1.x snake_case shape. Do not migrate them to 2.0 here — that must be done
// alongside the frozen resource files.

// Webhook is a registered webhook subscription.
type Webhook struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Events      []string `json:"events"`
	Description string   `json:"description,omitempty"`
	Status      string   `json:"status"`
	CreatedAt   string   `json:"created_at"`
}

// WebhookCreateResult is returned from Create. It is a Webhook plus the
// (never-shown-again) Secret used to verify deliveries.
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

// --- Chats -------------------------------------------------------------------

// Documented values for Chat.ChatType. The wire type is `string`.
const (
	ChatTypeContact    = "contact"
	ChatTypeGroup      = "group"
	ChatTypeNewsletter = "newsletter"
)

// Chat is a WhatsApp conversation — a 1:1 contact, a group, or a channel.
//
// Pointer fields preserve the wire-level null vs. present distinction. Name,
// LastMessageAt, UnreadCount are nullable; LastMessageText / LastMessageFromMe
// are present only when includeLastMessage=true.
type Chat struct {
	ID                string  `json:"chatId"`
	Name              *string `json:"name"`
	ChatType          string  `json:"chatType"`
	Pinned            *bool   `json:"pinned,omitempty"`
	Archived          *bool   `json:"archived,omitempty"`
	LastMessageAt     *string `json:"lastMessageAt"`
	UnreadCount       *int    `json:"unreadCount"`
	MarkedUnread      bool    `json:"markedUnread"`
	LastMessageText   *string `json:"lastMessageText,omitempty"`
	LastMessageFromMe *bool   `json:"lastMessageFromMe,omitempty"`
}

// Participant is a member of a group chat.
type Participant struct {
	ChatID       string  `json:"chatId"`
	Name         *string `json:"name,omitempty"`
	IsAdmin      bool    `json:"isAdmin"`
	IsSuperAdmin *bool   `json:"isSuperAdmin,omitempty"`
}

// OkResult is the {"ok": true} payload returned by chat/message write endpoints
// (archive, mark_read, pin, reactions, …) that have no richer body to surface.
type OkResult struct {
	Ok bool `json:"ok"`
}

// --- Contacts ----------------------------------------------------------------

// Contact is a WhatsApp contact in the connected account's address book.
type Contact struct {
	ChatID     string  `json:"chatId"`
	Name       *string `json:"name,omitempty"`
	IsBusiness *bool   `json:"isBusiness,omitempty"`
}

// --- Groups ------------------------------------------------------------------

// Group is a WhatsApp group's full metadata (GroupsResource.Get / Create /
// Update / member ops). Pointer fields are nullable on the wire; Participants is
// non-nil only when the group was fetched with Include="participants".
type Group struct {
	ID               string             `json:"id"`
	Name             *string            `json:"name"`
	Description      *string            `json:"description"`
	Owner            *string            `json:"owner"`
	CreatedAt        *string            `json:"createdAt"`
	LastMessageAt    *string            `json:"lastMessageAt"`
	ParticipantCount *int               `json:"participantCount"`
	Announce         *bool              `json:"announce"`
	Restrict         *bool              `json:"restrict"`
	Participants     []GroupParticipant `json:"participants,omitempty"`
}

// GroupListItem is the lighter group shape returned by GroupsResource.List and
// ContactsResource.CommonGroups (no announce flag, no participants array).
type GroupListItem struct {
	ID               string  `json:"id"`
	Name             *string `json:"name"`
	Description      *string `json:"description"`
	Owner            *string `json:"owner"`
	CreatedAt        *string `json:"createdAt"`
	LastMessageAt    *string `json:"lastMessageAt"`
	ParticipantCount *int    `json:"participantCount"`
	Restrict         *bool   `json:"restrict"`
}

// GroupParticipant is a single member of a Group (in Group.Participants).
type GroupParticipant struct {
	ChatID       string  `json:"chatId"`
	IsAdmin      bool    `json:"isAdmin"`
	IsSuperAdmin bool    `json:"isSuperAdmin"`
	Name         *string `json:"name"`
}

// CreateGroupParams is the body for GroupsResource.Create.
type CreateGroupParams struct {
	Name         string   `json:"name"`
	Participants []string `json:"participants"`
}

// GroupSettings is the optional settings block on UpdateGroupParams.
type GroupSettings struct {
	Announce           *bool   `json:"announce,omitempty"`
	Restrict           *bool   `json:"restrict,omitempty"`
	EditInfoAdminsOnly *bool   `json:"editInfoAdminsOnly,omitempty"`
	Description        *string `json:"description,omitempty"`
}

// UpdateGroupParams is the body for GroupsResource.Update.
type UpdateGroupParams struct {
	Name     string         `json:"name,omitempty"`
	Settings *GroupSettings `json:"settings,omitempty"`
}

// AddMemberParams is the body for GroupsResource.AddMember. Supply ChatID for a
// single member, or Participants for a batch (up to 256).
type AddMemberParams struct {
	ChatID       string   `json:"chatId,omitempty"`
	Participants []string `json:"participants,omitempty"`
}

// SetGroupPictureParams is the body for GroupsResource.SetPicture. Provide a
// base64 data URL via FileDataURL, or a public https URL via URL.
type SetGroupPictureParams struct {
	FileDataURL  string `json:"fileDataUrl,omitempty"`
	URL          string `json:"url,omitempty"`
	FileName     string `json:"fileName,omitempty"`
	FileMimeType string `json:"fileMimeType,omitempty"`
}

// GetGroupParams are the optional query params for GroupsResource.Get. Set
// Include="participants" to inline the members array.
type GetGroupParams struct {
	Include string
}

// --- Messages ----------------------------------------------------------------

// Documented values for a message Type discriminator. The wire type is `string`.
const (
	MessageTypeText  = "text"
	MessageTypeMedia = "media"
	MessageTypePoll  = "poll"
)

// Documented values for a media Kind. The wire type is `string`.
const (
	MediaKindImage    = "image"
	MediaKindVideo    = "video"
	MediaKindAudio    = "audio"
	MediaKindDocument = "document"
	MediaKindSticker  = "sticker"
	MediaKindVoice    = "voice"
	MediaKindGif      = "gif"
)

// Documented values for SentMessage.Status. The wire type is `string`.
const (
	MessageStatusPending        = "pending"
	MessageStatusBeforeWaSend   = "before-wa-send"
	MessageStatusBtSent         = "bt-sent"
	MessageStatusSending        = "sending"
	MessageStatusSent           = "sent"
	MessageStatusSentPendingAck = "sent_pending_ack"
	MessageStatusConfirmed      = "confirmed"
	MessageStatusDelivered      = "delivered"
	MessageStatusReceived       = "received"
	MessageStatusRead           = "read"
	MessageStatusPlayed         = "played"
	MessageStatusCancelled      = "cancelled"
	MessageStatusError          = "error"
	MessageStatusFailed         = "failed"
	MessageStatusExpired        = "expired"
)

// Documented values for Media.MediaUnavailable. The wire type is `string`; an
// empty value means the fetch succeeded (or no fetch was needed).
const (
	MediaUnavailableExpired        = "expired"
	MediaUnavailableFetching       = "fetching"
	MediaUnavailableAwaitingSender = "awaiting_sender"
	MediaUnavailableError          = "error"
	MediaUnavailableNoMedia        = "no_media"
)

// WaMessageKey is WhatsApp's structured message key.
type WaMessageKey struct {
	FromMe      bool   `json:"fromMe"`
	Remote      string `json:"remote,omitempty"`
	ID          string `json:"id,omitempty"`
	Serialized  string `json:"_serialized,omitempty"`
	Participant string `json:"participant,omitempty"`
}

// MessageLinkPreview is the optional rich-preview card embedded in a message.
// All fields are nullable on the wire.
type MessageLinkPreview struct {
	Title        *string `json:"title,omitempty"`
	Description  *string `json:"description,omitempty"`
	CanonicalURL *string `json:"canonicalUrl,omitempty"`
	Thumbnail    *string `json:"thumbnail,omitempty"`
}

// QuotedMessage is the message a Message quote-replies to.
type QuotedMessage struct {
	WaMessageKey WaMessageKey `json:"waMessageKey"`
	Text         string       `json:"text,omitempty"`
	Type         string       `json:"type,omitempty"`
	Author       string       `json:"author,omitempty"`
}

// Message is a chat-history message returned by MessagesResource.List / Get.
type Message struct {
	WaMessageKey  WaMessageKey        `json:"waMessageKey"`
	ChatID        string              `json:"chatId"`
	From          string              `json:"from"`
	Author        *string             `json:"author,omitempty"`
	SenderName    *string             `json:"senderName,omitempty"`
	Timestamp     *string             `json:"timestamp,omitempty"`
	Text          *string             `json:"text,omitempty"`
	Type          string              `json:"type"`
	FromMe        bool                `json:"fromMe"`
	Ack           *int                `json:"ack,omitempty"`
	MediaURL      *string             `json:"mediaUrl,omitempty"`
	Filename      *string             `json:"filename,omitempty"`
	LinkPreview   *MessageLinkPreview `json:"linkPreview,omitempty"`
	QuotedMessage *QuotedMessage      `json:"quotedMessage,omitempty"`
}

// SentMessage is the queued/scheduled/sent message record returned by
// MessagesResource.Send and every ScheduledMessagesResource endpoint
// (the MessageResponse schema).
//
// ID is the primary handle: the queue id for scheduled sends, or the serialized
// WhatsApp key for direct sends. Type: text | media | poll (MessageType*
// constants). Status: MessageStatus* constants. Nullable fields use pointers.
type SentMessage struct {
	ID                string              `json:"id,omitempty"`
	WaMessageKey      *WaMessageKey       `json:"waMessageKey"`
	To                string              `json:"to"`
	Type              string              `json:"type"`
	Text              *string             `json:"text"`
	MediaURL          *string             `json:"mediaUrl"`
	MediaKind         *string             `json:"mediaKind"`
	PollQuestion      *string             `json:"pollQuestion"`
	PollOptions       []string            `json:"pollOptions"`
	PollAllowMultiple *bool               `json:"pollAllowMultiple"`
	Status            string              `json:"status"`
	SendAt            *string             `json:"sendAt"`
	CreatedAt         string              `json:"createdAt"`
	ConfirmedAt       *string             `json:"confirmedAt"`
	ReceivedAt        *string             `json:"receivedAt"`
	ReadAt            *string             `json:"readAt"`
	PlayedAt          *string             `json:"playedAt"`
	FailedAt          *string             `json:"failedAt"`
	FailureReason     *string             `json:"failureReason"`
	Secret            *string             `json:"secret"`
	LinkPreview       *MessageLinkPreview `json:"linkPreview,omitempty"`
}

// MessageAck is the response from GET /v1/messages/ack/{waMessageKey}.
//
// Ack: -1=error, 0=pending, 1=server, 2=device, 3=read, 4=played; nil when no
// engine response was available.
type MessageAck struct {
	Ack *int `json:"ack"`
}

// BatchMessageAckEntry is a single per-message ack in BatchMessageAcksResult.
//
// Ack semantics mirror MessageAck.Ack; nil when undetermined or the key was not
// found (see Found).
type BatchMessageAckEntry struct {
	Key   string `json:"key"`
	Ack   *int   `json:"ack"`
	Found bool   `json:"found"`
}

// PinnedMessage is a single pinned message in PinnedMessagesPage.
type PinnedMessage struct {
	Key    string  `json:"key"`
	ChatID string  `json:"chatId"`
	Text   *string `json:"text"`
}

// Media is the response from GET /v1/messages/media/{waMessageKey}.
//
// MediaUnavailable (when set) is one of the MediaUnavailable* constants;
// OriginalQuality is false when the bytes are a synthesized lower-quality
// preview rather than the sender's original.
type Media struct {
	URL              *string `json:"url,omitempty"`
	Mimetype         *string `json:"mimetype,omitempty"`
	Filename         *string `json:"filename,omitempty"`
	DataBase64       *string `json:"dataBase64,omitempty"`
	OriginalQuality  *bool   `json:"originalQuality,omitempty"`
	MediaUnavailable *string `json:"mediaUnavailable,omitempty"`
}

// LoadOlderResult is the response from POST /v1/messages/load_older/{chatId}.
type LoadOlderResult struct {
	TotalMessages      *int    `json:"totalMessages"`
	Added              *int    `json:"added"`
	CanLoadMore        bool    `json:"canLoadMore"`
	HistoryUnavailable bool    `json:"historyUnavailable"`
	Error              *string `json:"error"`
}

// SendInChatParams is the body for MessagesResource.Send
// (POST /v1/messages/{chatId}) — a discriminated union on Type.
//
//	params := blueticks.SendInChatParams{Type: blueticks.MessageTypeText, Text: "hi"}
//
// The server validates the type/variant pair. IdempotencyKey is sent as the
// Idempotency-Key header, not in the body.
type SendInChatParams struct {
	Type              string   `json:"type"`
	ReplyTo           string   `json:"replyTo,omitempty"`
	Text              string   `json:"text,omitempty"`
	MediaURL          string   `json:"mediaUrl,omitempty"`
	MediaBase64       string   `json:"mediaBase64,omitempty"`
	MediaKind         string   `json:"mediaKind,omitempty"`
	MediaFilename     string   `json:"mediaFilename,omitempty"`
	PollQuestion      string   `json:"pollQuestion,omitempty"`
	PollOptions       []string `json:"pollOptions,omitempty"`
	PollAllowMultiple *bool    `json:"pollAllowMultiple,omitempty"`
	Secret            string   `json:"secret,omitempty"`
	WithTyping        *bool    `json:"withTyping,omitempty"`
	TypingSeconds     *float64 `json:"typingSeconds,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header, not in the body.
	IdempotencyKey string `json:"-"`
}

// GetMessageParams are the optional query params for MessagesResource.Get.
type GetMessageParams struct {
	ChatID string
}

// MessageAckParams are the optional query params for MessagesResource.GetAck.
type MessageAckParams struct {
	ChatID string
}

// GetMediaParams are the optional query params for MessagesResource.GetMedia.
type GetMediaParams struct {
	ChatID      string
	MaxAttempts int
}

// PinMessageParams configure MessagesResource.Pin. Duration (seconds) is
// optional; ChatID is an optional locator query param.
type PinMessageParams struct {
	ChatID   string
	Duration *int
}

// UnpinMessageParams are the optional query params for MessagesResource.Unpin.
type UnpinMessageParams struct {
	ChatID string
}

// ReactParams configure MessagesResource.React. Emoji is required (empty clears
// the reaction); ChatID is an optional locator query param.
type ReactParams struct {
	ChatID string
	Emoji  string
}

// BatchMessageAcksParams is the body for MessagesResource.BatchAcks.
type BatchMessageAcksParams struct {
	MessageKeys []string `json:"messageKeys"`
	ChatID      string   `json:"chatId,omitempty"`
}

// ListMessagesParams is the query-param bundle for MessagesResource.List.
// MessageTypes is sent as a single comma-separated value.
type ListMessagesParams struct {
	ChatID                string
	SearchToken           string
	Order                 string // "asc" or "desc"
	Since                 string
	Until                 string
	MessageTypes          []string
	LoadFromPhoneIfNeeded *bool
	IncludeMediaContent   *bool
	Skip                  int
	Limit                 int
}

// --- Scheduled Messages ------------------------------------------------------

// SendMessageParams is the body for ScheduledMessagesResource.Create
// (POST /v1/scheduled-messages/{chatId}) — a discriminated union on Type. Set
// SendAt to schedule; omit for immediate queueing. IdempotencyKey is sent as
// the Idempotency-Key header.
type SendMessageParams struct {
	Type              string   `json:"type"`
	ReplyTo           string   `json:"replyTo,omitempty"`
	Text              string   `json:"text,omitempty"`
	MediaURL          string   `json:"mediaUrl,omitempty"`
	MediaBase64       string   `json:"mediaBase64,omitempty"`
	MediaKind         string   `json:"mediaKind,omitempty"`
	MediaFilename     string   `json:"mediaFilename,omitempty"`
	PollQuestion      string   `json:"pollQuestion,omitempty"`
	PollOptions       []string `json:"pollOptions,omitempty"`
	PollAllowMultiple *bool    `json:"pollAllowMultiple,omitempty"`
	Secret            string   `json:"secret,omitempty"`
	SendAt            string   `json:"sendAt,omitempty"`
	// IdempotencyKey is sent as the Idempotency-Key header, not in the body.
	IdempotencyKey string `json:"-"`
}

// UpdateScheduledMessageParams is the body for ScheduledMessagesResource.Update
// (PATCH /v1/scheduled-messages/{id}). All fields optional; supply at least one.
// Pointer fields distinguish "leave unchanged" (nil) from an explicit value.
type UpdateScheduledMessageParams struct {
	Type              *string  `json:"type,omitempty"`
	Text              *string  `json:"text,omitempty"`
	MediaURL          *string  `json:"mediaUrl,omitempty"`
	MediaBase64       *string  `json:"mediaBase64,omitempty"`
	MediaKind         *string  `json:"mediaKind,omitempty"`
	MediaFilename     *string  `json:"mediaFilename,omitempty"`
	SendAt            *string  `json:"sendAt,omitempty"`
	ReplyTo           *string  `json:"replyTo,omitempty"`
	Secret            *string  `json:"secret,omitempty"`
	PollQuestion      *string  `json:"pollQuestion,omitempty"`
	PollAllowMultiple *bool    `json:"pollAllowMultiple,omitempty"`
	PollOptions       []string `json:"pollOptions,omitempty"`
}

// ListScheduledMessagesParams is the query-param bundle for
// ScheduledMessagesResource.List.
type ListScheduledMessagesParams struct {
	ChatID      string
	SearchToken string
	Status      string
	Order       string // "asc" or "desc"
	Skip        int
	Limit       int
}

// --- Engines -----------------------------------------------------------------

// EngineStatus is a single connected engine from GET /v1/engines.
//
// Type is one of the EngineType* constants. State / Stream / HasSynced are
// nullable on the wire (nil when unknown).
type EngineStatus struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Connected bool    `json:"connected"`
	State     *string `json:"state"`
	Stream    *string `json:"stream"`
	HasSynced *bool   `json:"hasSynced"`
}

// --- Newsletters -------------------------------------------------------------

// Documented values for a newsletter Verification. The wire type is `string`.
const (
	NewsletterVerificationVerified   = "VERIFIED"
	NewsletterVerificationUnverified = "UNVERIFIED"
)

// Newsletter is the response from POST /v1/newsletters and GET
// /v1/newsletters/{id} (the NewsletterResponse schema). Pointer fields are
// nullable on the wire.
type Newsletter struct {
	NewsletterID string  `json:"newsletterId"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	CreatedAt    *string `json:"createdAt"`
	Subscribers  *int    `json:"subscribers"`
	Invite       *string `json:"invite"`
	Verification *string `json:"verification"`
}

// NewsletterListItem is a newsletter row from GET /v1/newsletters. It keys on
// ChatID (the newsletter JID) rather than newsletterId.
type NewsletterListItem struct {
	ChatID       string  `json:"chatId"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	CreatedAt    *string `json:"createdAt"`
	Subscribers  *int    `json:"subscribers"`
	Invite       *string `json:"invite"`
	Verification *string `json:"verification"`
}

// CreateNewsletterParams is the body for NewslettersResource.Create.
type CreateNewsletterParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// --- Suno --------------------------------------------------------------------

// SunoAccount is the response from GET /v1/suno/account (AccountUsageResponse).
type SunoAccount struct {
	CreditsLeft  float64  `json:"creditsLeft"`
	MonthlyLimit *float64 `json:"monthlyLimit"`
	MonthlyUsage *float64 `json:"monthlyUsage"`
	Plan         *string  `json:"plan"`
}

// Documented values for SongClip.Status. The wire type is `string`; terminal
// states are "complete" and "error".
const (
	SongStatusSubmitted = "submitted"
	SongStatusQueued    = "queued"
	SongStatusRunning   = "running"
	SongStatusStreaming = "streaming"
	SongStatusComplete  = "complete"
	SongStatusError     = "error"
)

// SongClip is a single generated song variant (GetSongResponse / a clip in
// GenerateSongResult).
type SongClip struct {
	ID           string   `json:"id"`
	Status       string   `json:"status"`
	AudioURL     *string  `json:"audioUrl"`
	ImageURL     *string  `json:"imageUrl"`
	Title        *string  `json:"title"`
	DurationSec  *float64 `json:"durationSec"`
	Model        *string  `json:"model"`
	ErrorType    *string  `json:"errorType,omitempty"`
	ErrorMessage *string  `json:"errorMessage,omitempty"`
}

// GenerateSongResult is the response from POST /v1/suno/songs.
type GenerateSongResult struct {
	JobID string     `json:"jobId"`
	Clips []SongClip `json:"clips"`
}

// CreateUploadResult is the response from POST /v1/suno/uploads.
type CreateUploadResult struct {
	UploadID string `json:"uploadId"`
	Status   string `json:"status"`
}

// GenerateSongParams is the body for SunoResource.GenerateSong. Lyrics and Style
// are required; the rest are optional tuning knobs.
type GenerateSongParams struct {
	Lyrics         string   `json:"lyrics"`
	Style          string   `json:"style"`
	NegativeStyle  string   `json:"negativeStyle,omitempty"`
	VocalGender    string   `json:"vocalGender,omitempty"` // "m" or "f"
	Weirdness      *float64 `json:"weirdness,omitempty"`
	StyleInfluence *float64 `json:"styleInfluence,omitempty"`
	AudioInfluence *float64 `json:"audioInfluence,omitempty"`
	Instrumental   *bool    `json:"instrumental,omitempty"`
	Model          string   `json:"model,omitempty"`
	Title          string   `json:"title,omitempty"`
	UploadID       string   `json:"uploadId,omitempty"`
	CaptchaToken   string   `json:"captchaToken,omitempty"`
}

// CreateUploadParams is the body for SunoResource.CreateUpload. Provide AudioURL
// or AudioBase64.
type CreateUploadParams struct {
	AudioURL    string `json:"audioUrl,omitempty"`
	AudioBase64 string `json:"audioBase64,omitempty"`
	FileName    string `json:"fileName,omitempty"`
}

// --- Delete envelopes --------------------------------------------------------

// DeletedRef is the {"id": "..."} envelope returned by contact-delete endpoints.
type DeletedRef struct {
	ID string `json:"id"`
}

// DeletedResource is the {"id": "...", "deleted": true} envelope returned by
// soft-delete endpoints (audiences, scheduled-messages, webhooks).
type DeletedResource struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}
