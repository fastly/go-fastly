package signals

import "time"

// Signal is the API response to signals and signals operations.
type Signal struct {
	// CreatedAt is the created date and time the event was
	// created in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// Description is the user created description of the signal
	Description string `json:"description"`
	// Name is the user created name of the signal.
	Name string `json:"name"`
	// ReferenceID is the reference ID of the signal.
	ReferenceID string `json:"reference_id"`
	// Scope is the scope that the signal applies to
	Scope Scope `json:"scope"`
	// SignalID is the ID of the signal (auto generated).
	SignalID string `json:"id"`
	// UpdatedAt is the date and time in ISO 8601 format.
	UpdatedAt time.Time `json:"updated_at"`
}

// Signals is the API response structure for the list Signals operation.
type Signals struct {
	// Data is the list of returned signals.
	Data []Signal `json:"data"`
	// Meta is the information for total signals.
	Meta MetaSignals `json:"meta"`
}

// MetaSignals is a subset of the signals response structure.
type MetaSignals struct {
	// Limit is the limit of signals.
	Limit int `json:"limit"`
	// Total is the sum of signals.
	Total int `json:"total"`
}

// Scope is the definition of the scope that a signal applies to.
type Scope struct {
	// Type is the type of scope
	Type string `json:"type"`
	// AppliesTo defines what scope the signal applies to.
	AppliesTo []string `json:"applies_to"`
}

// Reason is the signal that corresponds to the reason an event was
// triggered.
type Reason struct {
	// Signal ID is the ID of the signal that triggered the event
	SignalID string `json:"signal_id"`
	// Count is the number of times this signal was detected
	Count int `json:"count"`
}
