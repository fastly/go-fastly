package versions

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
	// ActivatedAt is the date and time the version was last activated, in
	// ISO 8601 format.
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
	// Comment is the descriptive note associated with the version.
	Comment string `json:"comment,omitempty"`
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// VersionID is the version identifier.
	VersionID string `json:"id"`
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
