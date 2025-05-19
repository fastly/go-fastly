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
	// ID is the workspace VirtualPatch (UUID).
	ID string `json:"id"`
	// Mode is the mode of the VirtualPatch.
	Enabled string `json:"enabled"`
	// WorkspaceID is the workspace identifier (UUID).
	Mode string `json:"mode"`
	// WorkspaceID is the workspace identifier (UUID).
	WorkspaceID string `json:"workspace_id"`
	// Description is the description of the workspace.
	Description string `json:"description"`
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the date and time in ISO 8601 format.
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt is the date and time in ISO 8601 format.
	DeletedAt time.Time `json:"deleted_at"`
}
