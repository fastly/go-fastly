package operations

// Operation is the API response structure for operation endpoints.
type Operation struct {
	// ID is the unique identifier of the operation.
	ID string `json:"id"`
	// Method is the HTTP method for the operation.
	Method string `json:"method"`
	// Domain is the domain for the operation.
	Domain string `json:"domain"`
	// Path is the path for the operation.
	Path string `json:"path"`
	// Description describes what the operation does.
	Description string `json:"description,omitempty"`
	// Status is the discovery status of the operation (when present).
	Status string `json:"status,omitempty"`
	// TagIDs is the list of associated operation tag IDs.
	TagIDs []string `json:"tag_ids,omitempty"`
	// CreatedAt is when the operation was created.
	CreatedAt string `json:"created_at,omitempty"`
	// UpdatedAt is when the operation was last updated.
	UpdatedAt string `json:"updated_at,omitempty"`
	// LastSeenAt is when the operation was last seen in traffic.
	LastSeenAt string `json:"last_seen_at,omitempty"`
	// RPS is the observed requests per second for this operation.
	RPS float64 `json:"rps,omitempty"`
}

// Operations is the API response structure for listing operations.
type Operations struct {
	// Data is the list of returned operations.
	Data []Operation `json:"data"`
	// Meta contains pagination metadata.
	Meta Meta `json:"meta"`
}

// DiscoveredOperation is the API response structure for discovered operations.
type DiscoveredOperation struct {
	// ID is the unique identifier of the discovered operation.
	ID string `json:"id"`
	// Method is the HTTP method for the operation.
	Method string `json:"method"`
	// Domain is the domain for the operation.
	Domain string `json:"domain"`
	// Path is the path for the operation.
	Path string `json:"path"`
	// Status is the current status of the discovered operation.
	Status string `json:"status,omitempty"`
	// UpdatedAt is when the operation was last updated.
	UpdatedAt string `json:"updated_at,omitempty"`
	// LastSeenAt is when the operation was last seen in traffic.
	LastSeenAt string `json:"last_seen_at,omitempty"`
	// RPS is the observed requests per second for this discovered operation.
	RPS float64 `json:"rps,omitempty"`
}

// DiscoveredOperations is the API response structure for listing discovered operations.
type DiscoveredOperations struct {
	// Data is the list of returned discovered operations.
	Data []DiscoveredOperation `json:"data"`
	// Meta contains pagination metadata.
	Meta Meta `json:"meta"`
}

// OperationTag is the API response structure for operation tag endpoints.
type OperationTag struct {
	// ID is the unique identifier of the tag.
	ID string `json:"id"`
	// Name is the name of the tag.
	Name string `json:"name"`
	// Description describes the tag.
	Description string `json:"description,omitempty"`
	// Count is the number of operations associated with this tag.
	Count int `json:"count,omitempty"`
	// CreatedAt is when the tag was created.
	CreatedAt string `json:"created_at,omitempty"`
	// UpdatedAt is when the tag was last updated.
	UpdatedAt string `json:"updated_at,omitempty"`
}

// OperationTags is the API response structure for listing operation tags.
type OperationTags struct {
	// Data is the list of returned tags.
	Data []OperationTag `json:"data"`
	// Meta contains pagination metadata.
	Meta Meta `json:"meta"`
}

// Meta is a subset of pagination metadata returned by the API.
type Meta struct {
	// Limit is the maximum number of results returned.
	Limit int `json:"limit"`
	// Total is the total number of results.
	Total int `json:"total"`
}

// DiscoveredOperationStatusUpdate is the request body for updating a discovered operation status.
type DiscoveredOperationStatusUpdate struct {
	// Status is the new status to apply to the discovered operation.
	Status *string `json:"status"`
}

// DiscoveredOperationBulkStatusUpdate is the request body for bulk updating discovered operation status.
type DiscoveredOperationBulkStatusUpdate struct {
	// OperationIDs is the list of discovered operation IDs to update.
	OperationIDs []string `json:"operation_ids"`
	// Status is the new status to apply to the operations.
	Status *string `json:"status"`
}

// OperationBulkCreateItem is a single operation entry in a bulk create request.
type OperationBulkCreateItem struct {
	// Method is the HTTP method for the operation (required).
	Method *string `json:"method"`
	// Domain is the domain for the operation (required).
	Domain *string `json:"domain"`
	// Path is the path for the operation (required).
	Path *string `json:"path"`
	// Description is a description of the operation.
	Description *string `json:"description,omitempty"`
	// TagIDs is a list of associated operation tag IDs.
	TagIDs []string `json:"tag_ids,omitempty"`
}

// OperationBulkCreate is the request body for bulk creating operations.
type OperationBulkCreate struct {
	// Operations is the list of operations to create.
	Operations []OperationBulkCreateItem `json:"operations"`
}

// OperationBulkAddTags is the request body for bulk adding tags to operations.
type OperationBulkAddTags struct {
	// OperationIDs is the list of operation IDs to add tags to.
	OperationIDs []string `json:"operation_ids"`
	// TagIDs is the list of tag IDs to add to operations.
	TagIDs []string `json:"tag_ids"`
}

// BulkOperationResult is a single result item in bulk operations responses.
type BulkOperationResult struct {
	// ID is the operation ID.
	ID string `json:"id,omitempty"`
	// StatusCode is the HTTP status code for this operation.
	StatusCode int `json:"status_code,omitempty"`
	// Reason is the error reason if the operation failed.
	Reason string `json:"reason,omitempty"`
}

// BulkOperationResultsResponse is the response structure for bulk endpoints that return per-item results.
type BulkOperationResultsResponse struct {
	Data []BulkOperationResult `json:"data"`
}

// BulkCreateOperationsResult is a result entry for bulk-create operations.
type BulkCreateOperationsResult struct {
	BulkOperationResult
	// Operation is present when the operation creation succeeded.
	Operation *Operation `json:"operation,omitempty"`
}

// BulkCreateOperationsResponse is the response structure for bulk-create operations.
type BulkCreateOperationsResponse struct {
	Data []BulkCreateOperationsResult `json:"data"`
}
