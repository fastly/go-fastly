package redactions

import "time"

// Redaction is the API response structure for the create, update, and
// get operations.
type Redaction struct {
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// Field is the field being redacted. You cannot create multiple
	Field string `json:"field"`
	// Type is the type of field being redacted. One of
	// `request_parameter`, `request_header`, or
	// `response_header`.
	Type string `json:"type"`
	// RedactionID is the redaction identifier (UUID).
	RedactionID string `json:"id"`
}

// Redactions is the API response structure for the list Redactions
// operation.
type Redactions struct {
	// Data is the list of returned redactions.
	Data []Redaction `json:"data"`
	// Meta is the information for total redactions.
	Meta MetaRedactions `json:"meta"`
}

// MetaRedactions is a subset of the redactions response structure.
type MetaRedactions struct {
	// Limit is the limit of redactions.
	Limit int `json:"limit"`
	// Total is the sum of redactions.
	Total int `json:"total"`
}
