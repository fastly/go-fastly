package mailinglist

// Integration is the config object for integration type mailinglist.
type Integration struct {
	// Address An email address (required).
	Address *string `json:"address,omitempty"`
}