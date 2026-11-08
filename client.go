package blueticks

import (
	"errors"
	"net/http"
	"os"
)

// Client is the entry point for the Blueticks API. Construct with NewClient
// and access resources via the exported fields:
//
//	c, err := blueticks.NewClient(blueticks.WithAPIKey("sk_..."))
//	if err != nil { ... }
//	msg, err := c.ScheduledMessages.Create(ctx, blueticks.CreateScheduledMessageParams{Type: blueticks.MessageTypeText, To: "+15555550100", Text: "hi"})
type Client struct {
	apiKey     string
	baseURL    string
	userAgent  string
	httpClient *http.Client
	maxRetries int

	// Resources. Each is a thin wrapper over the same transport.
	Ping              *PingResource
	Account           *AccountResource
	Audiences         *AudiencesResource
	Campaigns         *CampaignsResource
	Webhooks          *WebhooksResource
	Chats             *ChatsResource
	Contacts          *ContactsResource
	Groups            *GroupsResource
	Engines           *EnginesResource
	Newsletters       *NewslettersResource
	ScheduledMessages *ScheduledMessagesResource
	Messages          *MessagesResource
	Suno              *SunoResource
	Boards            *BoardsResource
	Identity          *IdentityResource
	LatestChats       *LatestChatsResource
	Uploads           *UploadsResource
	WaLabels          *WaLabelsResource
	Realtime          *RealtimeResource
}

// Option configures a Client at construction time.
type Option func(*Client)

// WithAPIKey sets the API key. Overrides BLUETICKS_API_KEY.
func WithAPIKey(apiKey string) Option {
	return func(c *Client) { c.apiKey = apiKey }
}

// WithBaseURL overrides the API base URL. Useful for staging /
// self-hosted deployments. Defaults to https://api.blueticks.co.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = baseURL }
}

// WithHTTPClient swaps in a custom *http.Client. Useful for plugging in
// instrumentation, custom transports, or test stubs. Caller is responsible
// for setting an appropriate timeout.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) { c.httpClient = httpClient }
}

// WithUserAgent overrides the default User-Agent header sent on every
// request. The default is "blueticks-go/<Version>".
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// WithMaxRetries overrides the per-call retry cap. Default is 3.
// Setting 0 disables retries.
func WithMaxRetries(n int) Option {
	return func(c *Client) {
		if n >= 0 {
			c.maxRetries = n
		}
	}
}

// NewClient constructs a configured *Client. Returns an error when no
// API key was supplied via WithAPIKey or BLUETICKS_API_KEY.
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		baseURL:    defaultBaseURL,
		maxRetries: defaultMaxRetries,
		userAgent:  "blueticks-go/" + Version,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}

	if v := os.Getenv("BLUETICKS_API_KEY"); v != "" {
		c.apiKey = v
	}
	if v := os.Getenv("BLUETICKS_BASE_URL"); v != "" {
		c.baseURL = v
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.apiKey == "" {
		return nil, errors.New("blueticks: api key is required (pass WithAPIKey or set BLUETICKS_API_KEY)")
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: defaultTimeout}
	}
	if c.userAgent == "" {
		c.userAgent = "blueticks-go/" + Version
	}

	c.Ping = &PingResource{client: c}
	c.Account = &AccountResource{client: c}
	c.Audiences = &AudiencesResource{client: c}
	c.Campaigns = &CampaignsResource{client: c}
	c.Webhooks = &WebhooksResource{client: c}
	c.Chats = &ChatsResource{client: c}
	c.Contacts = &ContactsResource{client: c}
	c.Groups = &GroupsResource{client: c}
	c.Engines = &EnginesResource{client: c}
	c.Newsletters = &NewslettersResource{client: c}
	c.ScheduledMessages = &ScheduledMessagesResource{client: c}
	c.Messages = &MessagesResource{client: c}
	c.Suno = &SunoResource{client: c}
	c.Boards = &BoardsResource{client: c}
	c.Identity = &IdentityResource{client: c}
	c.LatestChats = &LatestChatsResource{client: c}
	c.Uploads = &UploadsResource{client: c}
	c.WaLabels = &WaLabelsResource{client: c}
	c.Realtime = &RealtimeResource{client: c}

	return c, nil
}
