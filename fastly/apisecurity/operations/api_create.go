package operations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v14/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string `json:"-"`
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

// Create creates a new operation associated with a service.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*Operation, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}
	if i.Method == nil {
		return nil, fastly.ErrMissingMethod
	}
	if i.Domain == nil {
		return nil, fastly.ErrMissingDomain
	}
	if i.Path == nil {
		return nil, fastly.ErrMissingPath
	}

	path := fastly.ToSafeURL("api-security", "v1", "services", *i.ServiceID, "operations")

	resp, err := c.PostJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var op *Operation
	if err := json.NewDecoder(resp.Body).Decode(&op); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return op, nil
}
