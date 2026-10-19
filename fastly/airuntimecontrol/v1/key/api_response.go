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
	// RPMLimit is the requests-per-minute limit for this key. `null` if there is
	// no limit.
	RPMLimit *int `json:"rpm_limit"`
	// TPMLimit is the tokens-per-minute limit for this key. `null` if there is
	// no limit.
	TPMLimit *int `json:"tpm_limit"`
	// SpendLimitEnabled is whether a spend limit is enabled for this key.
	SpendLimitEnabled bool `json:"spend_limit_enabled"`
	// SpendLimitEnabledOn is the date and time `spend_limit_enabled` was last
	// set to `true`, in RFC 3339 format. `null` if the spend limit has never been
	// enabled. This timestamp is not cleared when `spend_limit_enabled` is later
	// set back to `false`, so it can be non-null even when the spend limit is
	// currently disabled. (Read-only)
	SpendLimitEnabledOn *string `json:"spend_limit_enabled_on"`
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
	// SecurityEnabled is whether AI Firewall security is enabled for requests made
	// with this key. `null` if it has not been set.
	SecurityEnabled *bool `json:"security_enabled"`
	// SecurityAction is the action AI Firewall takes when it detects a violation
	// on a request made with this key. One of `log`, `block`, `replace`. `null`
	// if it has not been set.
	SecurityAction *string `json:"security_action"`
	// RPMLimit is the requests-per-minute limit for this key. `null` if there is
	// no limit.
	RPMLimit *int `json:"rpm_limit"`
	// TPMLimit is the tokens-per-minute limit for this key. `null` if there is
	// no limit.
	TPMLimit *int `json:"tpm_limit"`
	// SpendLimitEnabled is whether a spend limit is enabled for this key.
	SpendLimitEnabled bool `json:"spend_limit_enabled"`
	// SpendLimitEnabledOn is the date and time `spend_limit_enabled` was last
	// set to `true`, in RFC 3339 format. `null` if the spend limit has never been
	// enabled. This timestamp is not cleared when `spend_limit_enabled` is later
	// set back to `false`, so it can be non-null even when the spend limit is
	// currently disabled. (Read-only)
	SpendLimitEnabledOn *string `json:"spend_limit_enabled_on"`
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

// FailoverTarget is a provider connection in a virtual key's failover policy.
type FailoverTarget struct {
	// ProviderConnectionID is the ID of the provider connection. The first entry
	// in a policy's targets is the primary connection; subsequent entries are
	// fallbacks, tried in order.
	ProviderConnectionID string `json:"provider_connection_id"`
}

// FailoverPolicy is the API response structure for the get failover policy and
// update failover policy operations.
type FailoverPolicy struct {
	// Targets is the ordered list of provider connections requests are sent to.
	Targets []FailoverTarget `json:"targets"`
	// TimeoutEnabled is whether timeouts trigger failover.
	TimeoutEnabled bool `json:"timeout_enabled"`
	// RateLimit429Enabled is whether 429 rate-limit responses trigger failover.
	RateLimit429Enabled bool `json:"rate_limit_429_enabled"`
	// ServerError5xxEnabled is whether 5xx server errors trigger failover.
	ServerError5xxEnabled bool `json:"server_error_5xx_enabled"`
}
