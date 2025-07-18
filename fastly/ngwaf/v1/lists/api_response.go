package lists

import "time"

// List is the API response structure for the create, update, and get
// operations.
type List struct {
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// Description is the description of the list.
	Description string `json:"description"`
	// Entries are the entries of the list.
	Entries []string `json:"entries"`
	// ListID is the list identifier (UUID).
	ListID string `json:"id"`
	// Name is the name of the list.
	Name string `json:"name"`
	// ReferenceID is the reference ID of the list.
	ReferenceID string `json:"reference_id"`
	// Scope is the scope of the list.
	Scope Scope `json:"scope"`
	// Type is the type of the list.
	Type string `json:"type"`
	// UpdatedAt is the date and time in ISO 8601 format.
	UpdatedAt time.Time `json:"updated_at"`
}

// Scope is the API response structure for the scope of the list.
type Scope struct {
	// Type is the type of the scope.
	Type string `json:"type"`
}

// Lists is the API response structure for the list lists operation.
type Lists struct {
	// Data is the list of returned lists.
	Data []List `json:"data"`
	// Meta is the information for total lists.
	Meta MetaLists `json:"meta"`
}

// MetaLists is a subset of the Lists response structure.
type MetaLists struct {
	// Total is the sum of lists.
	Total int `json:"total"`
}
