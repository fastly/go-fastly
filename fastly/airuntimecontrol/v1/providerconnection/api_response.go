package providerconnection

// ProviderConnection is the API response structure for the create, get, update,
// and list operations.
type ProviderConnection struct {
	// ID is the unique identifier of the provider connection. (Read-only)
	ID string `json:"id"`
	// SecretID is an opaque identifier for the secret material backing this
	// connection. Not the secret itself. (Read-only)
	SecretID string `json:"secret_id,omitempty"`
	// Name is the display name of the provider connection.
	Name string `json:"name"`
	// Models is the AI model identifiers this connection is allowed to use.
	Models []string `json:"models"`
	// BaseURL is the base URL of the provider's API. Present only for `api-key`
	// connections.
	BaseURL string `json:"base_url"`
	// CreatedAt is the date and time the connection was created, in RFC 3339
	// format. (Read-only)
	CreatedAt string `json:"created_at,omitempty"`
	// UpdatedAt is the date and time the connection was last updated, in RFC 3339
	// format. (Read-only)
	UpdatedAt string `json:"updated_at,omitempty"`
}

// ProviderConnections is the API response structure for the list operation.
type ProviderConnections struct {
	// Data is the list of returned provider connections.
	Data []ProviderConnection `json:"data"`
	// Meta is metadata about the request.
	Meta Meta `json:"meta"`
}

// Meta contains metadata returned by the list operation.
type Meta struct {
	// Sort is the `sort` value used when making the request.
	Sort string `json:"sort"`
	// Total is the count of results matching the request.
	Total int `json:"total"`
}
