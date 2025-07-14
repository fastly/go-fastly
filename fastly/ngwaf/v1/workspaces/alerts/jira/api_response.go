package jira

// ResponseConfig is the config object for integration type jira in API responses.
type ResponseConfig struct {
	// Host is the name of the Jira instance (required).
	Host *string `json:"host"`
	// IssueType is the Jira issue type associated with the ticket.
	IssueType *string `json:"issue_type,omitempty"`
	// Key is the Jira API key / secret field (required).
	Key *string `json:"key"`
	// Project specifies the Jira project where the issue will be created (required).
	Project *string `json:"project"`
	// Username is the Jira username of the user who created the ticket (required).
	Username *string `json:"username"`
}

// MetaAlerts is a subset of the Alerts response structure.
type MetaAlerts struct {
	// Limit is the limit of Alert.
	Limit int `json:"limit"`
	// Total is the sum of Alert.
	Total int `json:"total"`
}

// Alert is the API response structure for the create, get and update
// alert operations.
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

// Alerts is the API response structure for the list alert operation.
type Alerts struct {
	// Data is the list of returned alerts.
	Data []Alert `json:"data"`
	// Meta is the information for total alerts.
	Meta MetaAlerts `json:"meta"`
}
