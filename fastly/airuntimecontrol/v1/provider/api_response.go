package provider

// Model represents an AI model offered by a provider.
type Model struct {
	// ID is the model identifier.
	ID string `json:"id"`
	// DisplayName is the human-readable name of the model.
	DisplayName string `json:"display_name"`
	// ProviderID is the ID of the provider this model belongs to.
	ProviderID string `json:"provider_id"`
}

// Provider represents an AI provider supported by AI Runtime Control.
type Provider struct {
	// ID is the provider identifier (e.g., `anthropic`, `openai`).
	ID string `json:"id"`
	// DisplayName is the human-readable name of the provider.
	DisplayName string `json:"display_name"`
	// DefaultBaseURL is the default API base URL for the provider.
	DefaultBaseURL string `json:"default_base_url"`
	// Models is the models available from this provider.
	Models []Model `json:"models"`
}

// Providers is the API response structure for the List operation.
type Providers struct {
	// Data is the list of returned providers.
	Data []Provider `json:"data"`
	// Meta is metadata about the request.
	Meta Meta `json:"meta"`
}

// Models is the API response structure for the ListModels operation.
type Models struct {
	// Data is the list of returned models.
	Data []Model `json:"data"`
	// Meta is metadata about the request.
	Meta Meta `json:"meta"`
}

// Meta is the metadata returned by the list operations.
type Meta struct {
	// Total is the count of results returned.
	Total int `json:"total"`
}
