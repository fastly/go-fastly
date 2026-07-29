package provider

// Model represents an AI model offered by a provider.
type Model struct {
	// ID is the model identifier.
	ID string `json:"id"`
	// DisplayName is the human-readable model name.
	DisplayName string `json:"display_name"`
	// ProviderID is the ID of the provider this model belongs to.
	ProviderID string `json:"provider_id"`
}

// Provider represents an AI provider supported by ARC.
type Provider struct {
	// ID is the provider identifier (e.g. "anthropic", "openai").
	ID string `json:"id"`
	// DisplayName is the human-readable provider name.
	DisplayName string `json:"display_name"`
	// DefaultBaseURL is the default API base URL for the provider.
	DefaultBaseURL string `json:"default_base_url"`
	// Models is the list of models available for the provider.
	Models []Model `json:"models"`
}

// Providers is the API response structure for the List operation.
type Providers struct {
	// Data is the list of returned providers.
	Data []Provider `json:"data"`
	// Meta contains total-count metadata.
	Meta Meta `json:"meta"`
}

// Models is the API response structure for the ListModels operation.
type Models struct {
	// Data is the list of returned models.
	Data []Model `json:"data"`
	// Meta contains total-count metadata.
	Meta Meta `json:"meta"`
}

// Meta is the metadata returned by the list operations.
type Meta struct {
	// Total is the total number of results.
	Total int `json:"total"`
}
