package paths

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
	// Links contains hyperlinks to resources related to the path.
	Links *Links `json:"links,omitempty"`
	// Path is the URL path pattern (max 2048 characters, starts with "/").
	Path string `json:"path"`
	// PathID is the path identifier.
	PathID string `json:"id"`
	// UpdatedAt is the date and time in ISO 8601 format.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// Links contains hyperlinks to resources related to a path.
type Links struct {
	// Rules is the URL of the path's rules collection.
	Rules string `json:"rules,omitempty"`
	// Self is the URL of the path itself.
	Self string `json:"self,omitempty"`
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
