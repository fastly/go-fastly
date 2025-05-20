package virtualpatches

import "time"

// MetaVirtualPatch is a subset of the VirtualPatch response structure.
type MetaVirtualPatch struct {
	// Limit is the limit of VirtualPatch.
	Limit int `json:"limit"`
	// Total is the sum of VirtualPatch.
	Total int `json:"total"`
}

// VirtualPatch is the API response structure for the list of virtial patch operations.
type VirtualPatch struct {
	// Enabled is the toggle status indicator of the VirtualPatch.
	Enabled string `json:"enabled"`
	// Description is the description of the workspace.
	Description string `json:"description"`
	// ID is the virtual patch identifier.
	ID string `json:"id"`
	// Mode is action to take when a signal for virtual patch is detected.
	Mode string `json:"mode"`
	// WorkspaceID is the workspace identifier (UUID).
	WorkspaceID string `json:"workspace_id"`
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// DeletedAt is the date and time in ISO 8601 format.
	DeletedAt time.Time `json:"deleted_at"`
	// UpdatedAt is the date and time in ISO 8601 format.
	UpdatedAt time.Time `json:"updated_at"`
}
