package workspaces

import "time"

// Workspace is the API response structure for the create, update, and get operations.
type Workspace struct {
	// AttackSignalThresholds are the parameters for system site alerts.
	AttackSignalThresholds AttackSignalThresholds `json:"attack_signal_thresholds"`
	// ClientIPHeaders lists the request headers containing the client IP address.
	ClientIPHeaders []string `json:"client_ip_headers"`
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// DefaultBlockingResponseCode is the default response code.
	DefaultBlockingResponseCode int `json:"default_blocking_response_code"`
	// Description is the description of the workspace.
	Description string `json:"description"`
	// IPAnonymization is the selected option to anonymize IP addresses.
	IPAnonymization string `json:"ip_anonymization"`
	// Mode is the mode of the workspace.
	Mode string `json:"mode"`
	// Name is the name of the workspace.
	Name string `json:"name"`
	// UpdatedAt is the date and time in ISO 8601 format.
	UpdatedAt time.Time `json:"updated_at"`
	// WorkspaceID is the workspace identifier (UUID).
	WorkspaceID string `json:"id"`
}

// AttackSignalThresholds are the parameters for system site alerts.
type AttackSignalThresholds struct {
	Immediate  bool `json:"immediate"`
	OneHour    int  `json:"one_hour"`
	OneMinute  int  `json:"one_minute"`
	TenMinutes int  `json:"ten_minutes"`
}

// Workspaces is the API response structure for the list workspaces operation.
type Workspaces struct {
	// Data is the list of returned workspaces.
	Data []Workspace `json:"data"`
	// Meta is the information for total workspaces.
	Meta MetaWorkspaces `json:"meta"`
}

// MetaWorkspaces is a subset of the Workspaces response structure.
type MetaWorkspaces struct {
	// Limit is the limit of workspaces.
	Limit int `json:"limit"`
	// Total is the sum of workspaces.
	Total int `json:"total"`
}
