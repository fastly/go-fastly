package key

// VirtualKey is the API response structure for the update operation.
type VirtualKey struct {
	// ID is the unique identifier of the virtual key.
	ID string `json:"id"`
	// Name is the human-readable name of the virtual key.
	Name string `json:"name"`
	// UserID is the ID of the user who created the key.
	UserID string `json:"user_id"`
	// CustomerID is the ID of the customer that owns the key.
	CustomerID string `json:"customer_id"`
	// Model is the AI model identifier.
	Model string `json:"model"`
	// Provider is the AI provider name.
	Provider string `json:"provider"`
	// ExpiresAt is the expiration timestamp of the key, if any.
	ExpiresAt *string `json:"expires_at"`
	// CreatedAt is when the key was created.
	CreatedAt string `json:"created_at,omitempty"`
	// UpdatedAt is when the key was last updated.
	UpdatedAt string `json:"updated_at,omitempty"`
}

// VirtualKeyWithToken is the API response structure for the create and rotate
// operations. The access token is only returned by those operations.
type VirtualKeyWithToken struct {
	VirtualKey
	// UserName is the display name of the user who created the key.
	UserName string `json:"user_name,omitempty"`
	// AccessToken is the generated access token (only returned on create/rotate).
	AccessToken string `json:"access_token,omitempty"`
}

// VirtualKeyListItem is the API response structure for the get and list
// operations.
type VirtualKeyListItem struct {
	// ID is the unique identifier of the virtual key.
	ID string `json:"id"`
	// Type is the type of the virtual key.
	Type string `json:"type,omitempty"`
	// Name is the human-readable name of the virtual key.
	Name string `json:"name"`
	// Model is the AI model identifier.
	Model string `json:"model"`
	// Provider is the AI provider name.
	Provider string `json:"provider"`
	// UserID is the ID of the user who created the key.
	UserID string `json:"user_id,omitempty"`
	// UserName is the display name of the user who created the key. Returned by
	// the list operation.
	UserName string `json:"user_name,omitempty"`
	// CustomerID is the ID of the customer that owns the key. Returned by the
	// list operation.
	CustomerID string `json:"customer_id,omitempty"`
	// CreatedAt is when the key was created.
	CreatedAt string `json:"created_at,omitempty"`
	// UpdatedAt is when the key was last updated.
	UpdatedAt string `json:"updated_at,omitempty"`
	// ExpiresAt is the expiration timestamp of the key, if any.
	ExpiresAt *string `json:"expires_at"`
	// DeletedAt is when the key was deleted, if applicable.
	DeletedAt *string `json:"deleted_at"`
	// LastUsedAt is when the key was last used, if applicable.
	LastUsedAt *string `json:"last_used_at"`
	// CreatedBy is the display name of the user who created the key.
	CreatedBy string `json:"created_by,omitempty"`
}

// VirtualKeys is the API response structure for the list operation.
type VirtualKeys struct {
	// Data is the list of returned virtual keys.
	Data []VirtualKeyListItem `json:"data"`
	// Meta contains pagination metadata.
	Meta Meta `json:"meta"`
}

// Meta is the pagination metadata returned by the list operation.
type Meta struct {
	// NextCursor is the cursor used to retrieve the next page of results.
	// It is empty when there are no further pages.
	NextCursor string `json:"next_cursor"`
	// Limit is the maximum number of results returned per page.
	Limit int `json:"limit"`
	// Sort is the sort field applied to the results.
	Sort string `json:"sort"`
	// Total is the total number of results.
	Total int `json:"total"`
}
