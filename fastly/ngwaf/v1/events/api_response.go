package events

import (
	"time"

	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/requests"
)

// Event is the API response structure for the get, list, and expire operations.
type Event struct {
	// Action value can be 'flagged' (requests will be blocked), 'info' (requests will be logged), or 'template'.
	Action string `json:"action"`
	// BlockSignals is the list of block signals.
	BlockSignals []string `json:"block_signals"`
	// BlockedRequestCount is the number of blocked reqests
	BlockedRequestCount int `json:"blocked_request_count"`
	// Country is the country code.
	Country string `json:"country"`
	// CreatedAt is the created date and time the event was created in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// DetectedAt is the date and time the event was detected in ISO 8601 format.
	DetectedAt time.Time `json:"detected_at"`
	// ExpiresAt is the date and time the event expires in ISO 8601 format.
	ExpiresAt time.Time `json:"expires_at"`
	// EventID is the event identifier.
	EventID string `json:"id"`
	// FlaggedRequestCount is the number of flagged reqests
	FlaggedRequestCount int `json:"flagged_request_count"`
	// IsExpired if true, the event should be set to expired.
	IsExpired bool `json:"is_expired"`
	// Reasons is a list of signals and their counts.
	Reasons []Reason `json:"reasons"`
	// RemoteHostname is the hostname of the event
	RemoteHostname string `json:"remote_hostname"`
	// RequestCount is the total numer of requests.
	RequestCount int `json:"request_count"`
	// SampleRequest is an example of a request that triggered the event.
	SampleRequest requests.Request `json:"sample_request"`
	// Source is the IP address of the source of the event.
	Source string `json:"source"`
	// Type is the type of event
	Type string `json:"type"`
	// UserAgents is a list of user agents contained in the event requests.
	UserAgents []string `json:"user_agents"`
	// Window is the time in seconds where the items were detected.
	Window int `json:"window"`
}

// Reason is the signal that corresponds to the reason an event was triggered.
type Reason struct {
	// Signal ID is the ID of the signal that triggered the event
	SignalID string `json:"signal_id"`
	// Count is the number of times this signal was detected
	Count int `json:"count"`
}

// Events is the API response structure for the list events operation.
type Events struct {
	// Data is the list of returned workspaces.
	Data []Event `json:"data"`
	// Meta is the information for total workspaces.
	Meta MetaEvents `json:"meta"`
}

// MetaEvents is a subset of the Event response structure.
type MetaEvents struct {
	// Limit is the limit of events.
	Limit int `json:"limit"`
	// Total is the sum of events.
	Total int `json:"total"`
}
