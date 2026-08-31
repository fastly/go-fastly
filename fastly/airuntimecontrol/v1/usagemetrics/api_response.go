package usagemetrics

// UsageMetric represents a single AI usage metric record.
type UsageMetric struct {
	// Date is the date of the usage record.
	Date string `json:"date"`
	// UsageType is the type of usage being measured (one of "requests",
	// "sessions", "input_tokens", "output_tokens").
	UsageType string `json:"usage_type"`
	// Quantity is the quantity of the usage type.
	Quantity int `json:"quantity"`
	// VirtualKeyID is the ID of the virtual key.
	VirtualKeyID string `json:"virtual_key_id"`
	// VirtualKeyName is the human-readable name of the virtual key.
	VirtualKeyName string `json:"virtual_key_name"`
	// Provider is the AI provider name.
	Provider string `json:"provider"`
	// Model is the AI model identifier.
	Model string `json:"model"`
}

// UsageMetrics is the API response structure for the List operation.
type UsageMetrics struct {
	// Data is the list of returned usage metrics.
	Data []UsageMetric `json:"data"`
	// Meta contains pagination metadata.
	Meta Meta `json:"meta"`
}

// Meta contains metadata returned by the list operation.
type Meta struct {
	// Sort is the sort field applied to the results.
	Sort string `json:"sort"`
	// Total is the total number of results.
	Total int `json:"total"`
}
