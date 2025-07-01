package workspacealerts

// MetaWorkspaceAlert is a subset of the WorkspaceAlerts response structure.
type MetaWorkspaceAlerts struct {
	// Limit is the limit of WorkspaceAlert.
	Limit int `json:"limit"`
	// Total is the sum of WorkspaceAlert.
	Total int `json:"total"`
}

// WorkspaceAlerts is the API response structure for the create, get and update
// workspace alert operations.
type WorkspaceAlert struct {
	// Description is an optional description for the alert.
	Description string `json:"description"`
	// ID is the workspace alert identifier.
	ID string `json:"id"`
	// Type is the tyoe of workspace integration.
	Type string `json:"type"`
	// Config is the configuration associated with the workspace integration.
	Config WorkspaceAlertConfig `json:"config"`
	// Events are the list of event types that trigger this webhook.
	Events []WorkspaceAlertConfig `json:"events"`
	// CreatedAt is a time stamp of when the alert was created.
	CreatedAt string `json:"created_at"`
	// CreatedBy is the email of the user who created the alert.
	CreatedBy string `json:"created_by"`
	// LastStatusCode is the HTTP status code recieved during that last webhook attempt.
	LastStatusCode string `json:"last_status_code"`
}

// WorkspaceAlertConfig is a subset of the WorkspaceAlert response structure.
type WorkspaceAlertConfig struct {
	// Webhook is the URL of the specified webhook.
	Flag string `json:"flag"`
}

// WorkspaceAlertEvemts is a subset of the WorkspaceAlert response structure.
type WorkspaceAlertEvemts struct {
	// Webhook is the URL of the specified webhook.
	Webhook string `json:"webhook"`
}

// WorkspaceAlerts is the API response structure for the list workspace alert operation.
type WorkspaceAlerts struct {
	// Data is the list of returned virtual patches.
	Data []WorkspaceAlert `json:"data"`
	// Meta is the information for total virtual patches.
	Meta MetaWorkspaceAlerts `json:"meta"`
}

// WorkspaceAlertsLey is the API response structure for the get and rotate workspace alert
// signing key operations.
type WorkspaceAlertsKey struct {
	// SigningKey is the details of a workspace alert signing key.
	SigningKey string `json:"signing_key"`
}
