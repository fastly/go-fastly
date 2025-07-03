package slack

// Integration is the config object for integration type slack.
type Integration struct {
	// Webhook is the Slack webhook (required).
	Webhook *string `json:"webhook,omitempty"`
}