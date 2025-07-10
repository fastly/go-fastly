package opsgenie

// ResponseConfig is the config object for integration type opsgenie in API responses.
type ResponseConfig struct {
	// Key is the Opsgenie integration key (required).
	Key *string `json:"key"`
}

// MetaAlerts is a subset of the Alerts response structure.
type MetaAlerts struct {
	// Limit is the limit of Alert.
	Limit int `json:"limit"`
	// Total is the sum of Alert.
	Total int `json:"total"`
}

// Alert is the API response structure for the create, get and update
// workspace alert operations.
type Alert struct {
	// Description is an optional description for the alert.
	Description string `json:"description,omitempty"`
	// ID is the workspace alert identifier.
	ID string `json:"id"`
	// Type is the type of workspace integration.
	Type string `json:"type"`
	// Config is the configuration associated with the workspace integration.
	Config ResponseConfig `json:"config"`
	// Events are the list of event types that trigger this webhook.
	Events []string `json:"events"`
	// CreatedAt is a time stamp of when the alert was created.
	CreatedAt string `json:"created_at"`
	// CreatedBy is the email of the user who created the alert.
	CreatedBy string `json:"created_by"`
	// LastStatusCode is the HTTP status code received during that last webhook attempt.
	LastStatusCode int `json:"last_status_code"`
}

// AlertEvent is a subset of the Alert response structure.
type AlertEvent struct {
	// Flag is the event flag.
	Flag string `json:"flag"`
}

// Alerts is the API response structure for the list workspace alert operation.
type Alerts struct {
	// Data is the list of returned workspace alerts.
	Data []Alert `json:"data"`
	// Meta is the information for total workspace alerts.
	Meta MetaAlerts `json:"meta"`
}
