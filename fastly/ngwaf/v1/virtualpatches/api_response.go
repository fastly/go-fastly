package virtualpatches

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
