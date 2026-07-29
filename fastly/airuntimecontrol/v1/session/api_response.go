package session

// Log contains the request/response data for a single session log entry.
type Log struct {
	// Request is the raw request data.
	Request string `json:"request"`
	// Response is the raw response data.
	Response string `json:"response"`
	// Attrs contains additional attributes.
	Attrs map[string]any `json:"attrs,omitempty"`
}

// Session represents a session log including request/response data and metadata.
type Session struct {
	// ID is the unique session identifier.
	ID string `json:"id"`
	// VirtualKeyID is the ID of the virtual key used.
	VirtualKeyID string `json:"virtual_key_id"`
	// VirtualKeyName is the human-readable name of the virtual key.
	VirtualKeyName string `json:"virtual_key_name"`
	// Model is the AI model identifier.
	Model string `json:"model"`
	// Provider is the AI provider name.
	Provider string `json:"provider"`
	// Requests is the number of requests in the session.
	Requests int `json:"requests"`
	// InputTokens is the total input tokens consumed.
	InputTokens int `json:"input_tokens"`
	// OutputTokens is the total output tokens generated.
	OutputTokens int `json:"output_tokens"`
	// CreatedAt is when the session was created.
	CreatedAt string `json:"created_at,omitempty"`
	// UpdatedAt is when the session was last updated.
	UpdatedAt string `json:"updated_at,omitempty"`
	// Logs is the list of request/response log entries for the session.
	Logs []Log `json:"logs,omitempty"`
}

// Sessions is the API response structure for the List operation.
type Sessions struct {
	// Data is the list of returned sessions.
	Data []Session `json:"data"`
	// Meta contains pagination metadata.
	Meta Meta `json:"meta"`
}

// Meta is the pagination metadata returned by the list operation.
type Meta struct {
	// NextCursor is the cursor used to retrieve the next page of results.
	// It is empty when there are no further pages.
	NextCursor string `json:"next_cursor"`
	// Limit is the maximum number of results returned per page.
	Limit int `json:"limit"`
	// Sort is the sort field applied to the results.
	Sort string `json:"sort"`
	// Total is the total number of results.
	Total int `json:"total"`
}
