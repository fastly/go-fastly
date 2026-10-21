package usagemetrics

// UsageMetric represents a single AI usage metric record.
type UsageMetric struct {
	// Date is the date the usage was measured, in `YYYY-MM-DD` format.
	Date string `json:"date"`
	// UsageType is the type of usage being measured (one of `requests`,
	// `sessions`, `input_tokens`, `output_tokens`).
	UsageType string `json:"usage_type"`
	// Quantity is the quantity of the usage type for this date.
	Quantity int `json:"quantity"`
	// VirtualKeyID is the ID of the virtual key this usage is attributed to.
	VirtualKeyID string `json:"virtual_key_id"`
	// VirtualKeyName is the display name of the virtual key this usage is
	// attributed to.
	VirtualKeyName string `json:"virtual_key_name"`
	// Provider is the ID of the provider connection that served the request.
	Provider string `json:"provider"`
	// PrimaryProvider is the ID of the provider connection that was originally
	// requested. Differs from Provider when the request failed over.
	PrimaryProvider string `json:"primary_provider"`
	// Model is the AI model identifier used.
	Model string `json:"model"`
}

// UsageMetrics is the API response structure for the List operation.
type UsageMetrics struct {
	// Data is the list of returned usage metrics.
	Data []UsageMetric `json:"data"`
	// Meta is metadata about the request.
	Meta Meta `json:"meta"`
}

// Meta is the pagination metadata returned by the list operation.
type Meta struct {
	// NextCursor is the cursor to use to request the next page of results. `null`
	// when there are no more results.
	NextCursor string `json:"next_cursor"`
	// Limit is the `limit` value used when making the request.
	Limit int `json:"limit"`
	// Sort is the `sort` value used when making the request.
	Sort string `json:"sort"`
	// Total is the count of results matching the request.
	Total int `json:"total"`
}
