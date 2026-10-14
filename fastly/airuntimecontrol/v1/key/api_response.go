package key

// VirtualKey is the API response structure for the update operation.
type VirtualKey struct {
	// ID is the unique identifier of the virtual key. (Read-only)
	ID string `json:"id"`
	// Name is the display name of the virtual key.
	Name string `json:"name"`
	// UserID is the ID of the user who created the key. (Read-only)
	UserID string `json:"user_id"`
	// CustomerID is the ID of the customer account that owns the key. (Read-only)
	CustomerID string `json:"customer_id"`
	// Model is the AI model identifier.
	Model string `json:"model"`
	// Provider is the AI provider name.
	Provider string `json:"provider"`
	// ExpiresAt is the date and time the key expires, in RFC 3339 format. `null`
	// if the key doesn't expire.
	ExpiresAt *string `json:"expires_at"`
	// CreatedAt is the date and time the key was created, in RFC 3339 format.
	// (Read-only)
	CreatedAt string `json:"created_at,omitempty"`
	// UpdatedAt is the date and time the key was last updated, in RFC 3339 format.
	// (Read-only)
	UpdatedAt string `json:"updated_at,omitempty"`
}

// VirtualKeyWithToken is the API response structure for the create and rotate
// operations. The access token is only returned by those operations.
type VirtualKeyWithToken struct {
	VirtualKey
	// UserName is the display name of the user who created the key. (Read-only)
	UserName string `json:"user_name,omitempty"`
	// AccessToken is the access token for this virtual key. This is only returned
	// when the key is created or rotated, and cannot be retrieved again afterward.
	// (Read-only)
	AccessToken string `json:"access_token,omitempty"`
}

// VirtualKeyListItem is the API response structure for the get and list
// operations.
type VirtualKeyListItem struct {
	// ID is the unique identifier of the virtual key. (Read-only)
	ID string `json:"id"`
	// Type is the type of the resource. (Read-only)
	Type string `json:"type,omitempty"`
	// Name is the display name of the virtual key.
	Name string `json:"name"`
	// Model is the AI model identifier.
	Model string `json:"model"`
	// Provider is the AI provider name.
	Provider string `json:"provider"`
	// UserID is the ID of the user who created the key. (Read-only)
	UserID string `json:"user_id,omitempty"`
	// UserName is the display name of the user who created the key. Returned by
	// the list operation. (Read-only)
	UserName string `json:"user_name,omitempty"`
	// CustomerID is the ID of the customer account that owns the key. Returned by
	// the list operation. (Read-only)
	CustomerID string `json:"customer_id,omitempty"`
	// CreatedAt is the date and time the key was created, in RFC 3339 format.
	// (Read-only)
	CreatedAt string `json:"created_at,omitempty"`
	// UpdatedAt is the date and time the key was last updated, in RFC 3339 format.
	// (Read-only)
	UpdatedAt string `json:"updated_at,omitempty"`
	// ExpiresAt is the date and time the key expires, in RFC 3339 format. `null`
	// if the key doesn't expire.
	ExpiresAt *string `json:"expires_at"`
	// DeletedAt is the date and time the key was deleted, in RFC 3339 format.
	// `null` unless the key was deleted and `include_deleted` was used to retrieve
	// it. (Read-only)
	DeletedAt *string `json:"deleted_at"`
	// LastUsedAt is the date and time the key was last used to make a request, in
	// RFC 3339 format. `null` if the key has never been used. (Read-only)
	LastUsedAt *string `json:"last_used_at"`
	// CreatedBy is the display name of the user who created the key. (Read-only)
	CreatedBy string `json:"created_by,omitempty"`
}

// VirtualKeys is the API response structure for the list operation.
type VirtualKeys struct {
	// Data is the list of returned virtual keys.
	Data []VirtualKeyListItem `json:"data"`
	// Meta is metadata about the request.
	Meta Meta `json:"meta"`
}

// Meta is the pagination metadata returned by the list operation.
type Meta struct {
	// NextCursor is the cursor to use to request the next page of results. `null`
	// when there are no more results.
	NextCursor string `json:"next_cursor"`
	// Limit is the `limit` value used when making the request.
	Limit int `json:"limit"`
	// Sort is the `sort` value used when making the request.
	Sort string `json:"sort"`
	// Total is the count of results matching the request.
	Total int `json:"total"`
}
