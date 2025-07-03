package datadog

// Integration is the config object for integration type datadog.
type Integration struct {
	// Key is the Datadog integration key (required).
	Key *string `json:"key,omitempty"`
	// Site is the Datadog site (required).
	Site *string `json:"site,omitempty"`
}