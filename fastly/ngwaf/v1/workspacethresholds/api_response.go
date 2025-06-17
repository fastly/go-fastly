package workspacethresholds

// DeleteWorkspaceThresholds
type DeleteWorkspaceThresholds struct {
	//
}

// MetaWorkspaceThresholds is a subset of the WorkspaceThresholds response structure.
type MetaWorkspaceThresholds struct {
	// Total is the sum of WorkspaceThresholds.
	Total int `json:"total"`
}

// UpdateWorkspaceThreshold is the API response structure for the update workspace thresholds operation.
type UpdateWorkspaceThreshold struct {
	// Action is the action to take when the theshold is exceeded.
	Action string `json:"action"`
	// DontNotify determines whether to silence notifications when action is taken.
	DontNotify bool `json:"dont_notify"`
	// Duration is the duration of the action in place.
	Duration int `json:"duration"`
	// Enabled is whether this threshold is active.
	Enabled bool `json:"enabled"`
	// Interval is the threshold interval in seconds.
	Interval int `json:"interval"`
	// Limit is the threshold limit.
	Limit int `json:"limit"`
	// Name is the threshold name.
	Name string `json:"name"`
	// Signal is the name of the signal this threshold is acting on.
	Signal string `json:"signal"`
}

// WorkspaceThreshold is the API response structure for the create and get workspace thresholds operations.
type WorkspaceThreshold struct {
	// Action is the action to take when the theshold is exceeded.
	Action string `json:"action"`
	// DontNotify determines whether to silence notifications when action is taken.
	DontNotify bool `json:"dont_notify"`
	// Duration is the duration of the action in place.
	Duration int `json:"duration"`
	// CreatedAt is when the workspace threshold was created.
	CreatedAt string `json:"created_at"`
	// Enabled is whether this threshold is active.
	Enabled bool `json:"enabled"`
	// ID is the workspace threshold identifier.
	ID string `json:"id"`
	// Interval is the threshold interval in seconds.
	Interval int `json:"interval"`
	// Limit is the threshold limit.
	Limit int `json:"limit"`
	// Name is the threshold name.
	Name string `json:"name"`
	// Signal is the name of the signal this threshold is acting on.
	Signal string `json:"signal"`
}

// WorkspaceThresholds is the API response structure for the list workspace thresholds operations.
type WorkspaceThresholds struct {
	// Data is the list of returned workspace thresholds.
	Data []WorkspaceThreshold `json:"data"`
	// Meta is the information for total workspace thresholds.
	Meta MetaWorkspaceThresholds `json:"meta"`
}
