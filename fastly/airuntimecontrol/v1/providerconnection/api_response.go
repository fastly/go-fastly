package providerconnection

// ProviderConnection is the API response structure for the create, get, update,
// and list operations.
type ProviderConnection struct {
	// ID is the unique identifier of the provider connection.
	ID string `json:"id"`
	// Name is the human-readable name of the provider.
	Name string `json:"name"`
	// Models is the list of allowed AI model identifiers.
	Models []string `json:"models"`
	// BaseURL is the base URL for the provider's API.
	BaseURL string `json:"base_url"`
	// CreatedAt is when the connection was created.
	CreatedAt string `json:"created_at,omitempty"`
	// UpdatedAt is when the connection was last updated.
	UpdatedAt string `json:"updated_at,omitempty"`
}

// ProviderConnections is the API response structure for the list operation.
type ProviderConnections struct {
	// Data is the list of returned provider connections.
	Data []ProviderConnection `json:"data"`
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
