package opsgenie

// Integration is the config object for integration type opsgenie.
type Integration struct {
	// Key is the Opsgenie integration key (required).
	Key *string `json:"key,omitempty"`
}