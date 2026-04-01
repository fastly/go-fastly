package operations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v14/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string `json:"-"`
	// OperationID is the unique identifier of the operation (required).
	OperationID *string `json:"-"`
	// Description is an updated description for the operation.
	Description *string `json:"description,omitempty"`
	// TagIDs is an updated list of associated operation tag IDs.
	TagIDs []string `json:"tag_ids,omitempty"`
}

// Update partially updates an existing operation associated with a service.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Operation, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}
	if i.OperationID == nil {
		return nil, fastly.ErrMissingID
	}

	path := fastly.ToSafeURL(
		"api-security", "v1", "services", *i.ServiceID, "operations", *i.OperationID,
	)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
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
