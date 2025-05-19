package virtualpatches

import "time"

// Workspace is the API response structure for the create, update, and get operations.
type Workspace struct {
	// WorkspaceID is the workspace identifier (UUID).
	WorkspaceID string `json:"id"`
	// Name is the name of the workspace.
	Name string `json:"name"`
	// Description is the description of the workspace.
	Description string `json:"description"`
	// Mode is the mode of the workspace.
	Mode string `json:"mode"`
	// AttackSignalThresholds are the parameters for system site alerts.
	AttackSignalThresholds AttackSignalThresholds `json:"attack_signal_thresholds"`
	// IPAnonymization is the selected option to anonymize IP addresses.
	IPAnonymization string `json:"ip_anonymization"`
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the date and time in ISO 8601 format.
	UpdatedAt time.Time `json:"updated_at"`
}

// AttackSignalThresholds are the parameters for system site alerts.
type AttackSignalThresholds struct {
	OneMinute  int  `json:"one_minute"`
	TenMinutes int  `json:"ten_minutes"`
	OneHour    int  `json:"one_hour"`
	Immediate  bool `json:"immediate"`
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

// VirtualPatch is the API response structure for the list of virtial patch operations.
type VirtualPatch struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
}
