package session

// Log contains the request/response data for a single session log entry.
type Log struct {
	// Request is the raw request data.
	Request string `json:"request"`
	// Response is the raw response data.
	Response string `json:"response"`
	// Attrs is additional attributes recorded for this request/response pair.
	Attrs map[string]any `json:"attrs,omitempty"`
}

// Session represents a session log including request/response data and metadata.
type Session struct {
	// ID is the unique identifier of the session. (Read-only)
	ID string `json:"id"`
	// VirtualKeyID is the ID of the virtual key used for this session.
	VirtualKeyID string `json:"virtual_key_id"`
	// VirtualKeyName is the display name of the virtual key used for this session.
	VirtualKeyName string `json:"virtual_key_name"`
	// Model is the AI model identifier used.
	Model string `json:"model"`
	// Provider is the ID of the provider connection that served the session.
	Provider string `json:"provider"`
	// PrimaryProvider is the ID of the provider connection that was originally
	// requested. Differs from Provider when the request failed over.
	PrimaryProvider string `json:"primary_provider"`
	// Requests is the number of requests made in this session.
	Requests int `json:"requests"`
	// InputTokens is the total input tokens consumed in this session.
	InputTokens int `json:"input_tokens"`
	// OutputTokens is the total output tokens generated in this session.
	OutputTokens int `json:"output_tokens"`
	// CreatedAt is the date and time the session was created, in RFC 3339 format.
	// (Read-only)
	CreatedAt string `json:"created_at,omitempty"`
	// UpdatedAt is the date and time the session was last updated, in RFC 3339
	// format. (Read-only)
	UpdatedAt string `json:"updated_at,omitempty"`
	// Logs is the request and response pairs recorded for this session.
	Logs []Log `json:"logs,omitempty"`
}

// Sessions is the API response structure for the List operation.
type Sessions struct {
	// Data is the list of returned sessions.
	Data []Session `json:"data"`
	// Meta is metadata about the request.
	Meta Meta `json:"meta"`
}

// Meta is the pagination metadata returned by the list operation.
type Meta struct {
	// NextCursor is the cursor to use to request the next page of results. `null`
	// when there are no more results.
	NextCursor string `json:"next_cursor"`
	// Limit is the `limit` value used when making the request.
	Limit int `json:"limit"`
	// Sort is the `sort` value used when making the request.
	Sort string `json:"sort"`
	// Total is the count of results matching the request.
	Total int `json:"total"`
}
