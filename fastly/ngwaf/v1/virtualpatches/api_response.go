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
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the date and time in ISO 8601 format.
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt is the date and time in ISO 8601 format.
	DeletedAt time.Time `json:"deleted_at"`
}
