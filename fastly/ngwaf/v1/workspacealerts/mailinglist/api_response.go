package mailinglist

// MetaWorkspaceAlerts is a subset of the WorkspaceAlerts response structure.
type MetaWorkspaceAlerts struct {
	// Limit is the limit of WorkspaceAlert.
	Limit int `json:"limit"`
	// Total is the sum of WorkspaceAlert.
	Total int `json:"total"`
}

// WorkspaceAlert is the API response structure for the create, get and update
// workspace alert operations.
type WorkspaceAlert struct {
	// Description is an optional description for the alert.
	Description string `json:"description"`
	// ID is the workspace alert identifier.
	ID string `json:"id"`
	// Type is the type of workspace integration.
	Type string `json:"type"`
	// Config is the configuration associated with the workspace integration.
	Config Config `json:"config"`
	// Events are the list of event types that trigger this webhook.
	Events []WorkspaceAlertEvent `json:"events"`
	// CreatedAt is a time stamp of when the alert was created.
	CreatedAt string `json:"created_at"`
	// CreatedBy is the email of the user who created the alert.
	CreatedBy string `json:"created_by"`
	// LastStatusCode is the HTTP status code received during that last webhook attempt.
	LastStatusCode string `json:"last_status_code"`
}

// WorkspaceAlertEvent is a subset of the WorkspaceAlert response structure.
type WorkspaceAlertEvent struct {
	// Flag is the event flag.
	Flag string `json:"flag"`
}

// WorkspaceAlerts is the API response structure for the list workspace alert operation.
type WorkspaceAlerts struct {
	// Data is the list of returned workspace alerts.
	Data []WorkspaceAlert `json:"data"`
	// Meta is the information for total workspace alerts.
	Meta MetaWorkspaceAlerts `json:"meta"`
}

// WorkspaceAlertsKey is the API response structure for the get and rotate workspace alert
// signing key operations.
type WorkspaceAlertsKey struct {
	// SigningKey is the details of a workspace alert signing key.
	SigningKey string `json:"signing_key"`
}