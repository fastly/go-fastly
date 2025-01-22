package v1

import (
	"time"
)

// Collection is the API response structure for the list operation.
type Collection struct {
	// Data contains the API data.
	Data []Data `json:"data"`
	// Meta contains metadata related to paginating the full dataset.
	Meta Meta `json:"meta"`
}

// Data is a subset of the API response structure containing the specific API
// data itself.
type Data struct {
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// ID is the domain identifier (UUID).
	DomainID string `json:"id"`
	// FQDN is the fully-qualified domain name of the domain. Read-only
	// after creation.
	FQDN string `json:"fqdn"`
	// ServiceID is the service_id associated with the domain or nil if there
	// is no association.
	ServiceID *string `json:"service_id"`
	// UpdatedAt is the date and time in ISO 8601 format.
	UpdatedAt time.Time `json:"updated_at"`
}

// Meta is a subset of the API response structure containing metadata related to
// paginating the full dataset.
type Meta struct {
	// Limit is how many results are included in this response.
	Limit int `json:"limit"`
	// NextCursor is the cursor value used to retrieve the next page.
	NextCursor string `json:"next_cursor"`
	// Sort is the field used to order the response by.
	Sort string `json:"sort"`
	// Total is the total number of results.
	Total int `json:"total"`
}
