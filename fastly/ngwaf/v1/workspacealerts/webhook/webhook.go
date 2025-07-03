package webhook

// Integration is the config object for integration type webhook.
type Integration struct {
	// Webhook is the Webhook URL (required).
	Webhook *string `json:"webhook,omitempty"`
}