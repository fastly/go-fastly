package routingconfigs

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
	// ActivatedAt is the date and time the currently active version was
	// activated, in ISO 8601 format. Only present once the routing config has
	// been activated at least once.
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// Links contains hyperlinks to resources related to the routing config.
	// Not present on every response (e.g. the deactivate response omits it).
	Links *Links `json:"links,omitempty"`
	// Name is the user-defined name of the routing config.
	Name string `json:"name"`
	// RoutingConfigID is the routing config identifier.
	RoutingConfigID string `json:"id"`
	// State is the current lifecycle state of the routing config (e.g.
	// "draft-only", "active", or "active-with-draft").
	State string `json:"state,omitempty"`
	// UpdatedAt is the date and time in ISO 8601 format.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// Links contains hyperlinks to resources related to a routing config.
type Links struct {
	// Paths is the URL of the routing config's paths collection.
	Paths string `json:"paths,omitempty"`
	// Self is the URL of the routing config itself.
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
