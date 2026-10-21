package providerconnection

// ProviderConnection is the API response structure for the create, get, update,
// and list operations.
type ProviderConnection struct {
	// ID is the unique identifier of the provider connection. (Read-only)
	ID string `json:"id"`
	// SecretID is an opaque identifier for the secret material backing this
	// connection. Not the secret itself. (Read-only)
	SecretID string `json:"secret_id,omitempty"`
	// Name is the display name of the provider connection.
	Name string `json:"name"`
	// Models is the AI model identifiers this connection is allowed to use.
	Models []string `json:"models"`
	// BaseURL is the base URL of the provider's API. Present only for `api-key`
	// connections.
	BaseURL string `json:"base_url,omitempty"`
	// AuthType is the authentication type for this connection. One of
	// `api-key`, `aws-iam`, `gcp-iam`.
	AuthType string `json:"auth_type"`
	// Region is the cloud region for API requests. Present only for `aws-iam`
	// and `gcp-iam` connections.
	Region string `json:"region,omitempty"`
	// VirtualKeyIDs is the IDs of the active virtual keys whose failover policy
	// references this connection. Empty if no virtual key references it.
	// Returned by the get operation. (Read-only)
	VirtualKeyIDs []string `json:"virtual_key_ids,omitempty"`
	// CreatedAt is the date and time the connection was created, in RFC 3339
	// format. (Read-only)
	CreatedAt string `json:"created_at,omitempty"`
	// UpdatedAt is the date and time the connection was last updated, in RFC 3339
	// format. (Read-only)
	UpdatedAt string `json:"updated_at,omitempty"`
	// DeletedAt is the date and time the connection was deleted, in RFC 3339
	// format. Present only for deleted connections, which are returned only when
	// `include_deleted` is `true`. (Read-only)
	DeletedAt *string `json:"deleted_at"`
}

// ProviderConnections is the API response structure for the list operation.
type ProviderConnections struct {
	// Data is the list of returned provider connections.
	Data []ProviderConnection `json:"data"`
	// Meta is metadata about the request.
	Meta Meta `json:"meta"`
}

// Meta contains metadata returned by the list operation.
type Meta struct {
	// Sort is the `sort` value used when making the request.
	Sort string `json:"sort"`
	// Total is the count of results matching the request.
	Total int `json:"total"`
}

// BedrockSetup is the API response structure for the GetBedrockSetup
// operation.
type BedrockSetup struct {
	// ArcRoleARN is the ARN of the AI Runtime Control IAM role that the
	// customer's AWS account IAM trust policy must trust.
	ArcRoleARN string `json:"arc_role_arn"`
	// ExternalID is the external ID to use in the IAM trust policy. This is
	// the customer ID.
	ExternalID string `json:"external_id"`
}

// VertexSetup is the API response structure for the GetVertexSetup operation.
type VertexSetup struct {
	// OIDCIssuerURL is the OIDC issuer URL for workload identity federation.
	OIDCIssuerURL string `json:"oidc_issuer_url"`
	// Audience is the audience value for the OIDC token.
	Audience string `json:"audience"`
	// CustomerID is the customer ID to use when configuring the workload
	// identity pool.
	CustomerID string `json:"customer_id"`
}
