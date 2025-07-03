package microsoftteams

// Integration is the config object for integration type microsoftteams.
type Integration struct {
	// Webhook is the Microsoft Teams webhook (required).
	Webhook *string `json:"webhook,omitempty"`
}