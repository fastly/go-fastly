package workspaces

// Workspace is the API response structure for the get, list, and update
// operations.
type Workspace struct {
	// ID is the unique identifier of the workspace.
	ID string `json:"id"`
	// Name is the display name of the workspace.
	Name string `json:"name"`
	// Description is a description of the workspace.
	Description string `json:"description,omitempty"`
	// ProtectionMode is the protection mode of the workspace. Must be one of
	// `off`, `log`, or `block`.
	ProtectionMode string `json:"protection_mode"`
	// Services is the list of IDs of the services attached to this workspace.
	Services []string `json:"services,omitempty"`
	// CreatedAt is the date and time the workspace was created, in RFC 3339
	// format.
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the date and time the workspace was last updated, in RFC
	// 3339 format.
	UpdatedAt string `json:"updated_at"`
}

// Workspaces is the API response structure for the list operation.
type Workspaces struct {
	// Data is the list of returned workspaces.
	Data []Workspace `json:"data"`
	// Meta contains pagination metadata.
	Meta Meta `json:"meta"`
}

// Meta is the pagination metadata returned by the list operation.
type Meta struct {
	// Limit is the `limit` value used when making the request.
	Limit *int `json:"limit"`
	// Total is the count of results matching the request.
	Total *int `json:"total"`
	// NextCursor is the cursor used to retrieve the next page of results. It
	// is empty when there are no further pages.
	NextCursor *string `json:"next_cursor"`
}
