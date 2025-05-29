package virtualpatches

import "time"

// VirtualPatch is the API response structure for the patch virtial patch operations.
type VirtualPatch struct {
	// Description is the description of the workspace.
	Description string `json:"description"`
	// Enabled is the toggle status indicator of the VirtualPatch.
	Enabled bool `json:"enabled"`
	// ID is the virtual patch identifier.
	ID string `json:"id"`
	// Mode is action to take when a signal for virtual patch is detected.
	Mode string `json:"mode"`
}

// MetaVirtualPatch is a subset of the VirtualPatch response structure.
type MetaVirtualPatches struct {
	// Limit is the limit of VirtualPatch.
	Limit int `json:"limit"`
	// Total is the sum of VirtualPatch.
	Total int `json:"total"`
}

// VirtualPatch is the API response structure for the get and list virtial patch operations.
type VirtualPatches struct {
	// Data is the list of returned virtual patches.
	Data []VirtualPatch `json:"data"`
	// Meta is the information for total virtual patches.
	Meta MetaVirtualPatches `json:"meta"`
}

// AttackSignalThresholds are the parameters for system site alerts.
type AttackSignalThresholds struct {
	OneMinute  int  `json:"one_minute"`
	TenMinutes int  `json:"ten_minutes"`
	OneHour    int  `json:"one_hour"`
	Immediate  bool `json:"immediate"`
}

// Workspace is the API response structure for the create operation.
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
