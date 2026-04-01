package operations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v14/fastly"
)

// UpdateDiscoveredStatusInput specifies the information needed for the
// UpdateDiscoveredStatus() function to perform the operation.
type UpdateDiscoveredStatusInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string `json:"-"`
	// OperationID is the unique identifier of the discovered operation (required).
	OperationID *string `json:"-"`
	// Status is the new status to apply (required).
	Status *string `json:"status"`
}

// UpdateDiscoveredStatus updates the status of a single discovered operation.
func UpdateDiscoveredStatus(ctx context.Context, c *fastly.Client, i *UpdateDiscoveredStatusInput) (*DiscoveredOperation, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}
	if i.OperationID == nil {
		return nil, fastly.ErrMissingID
	}
	if i.Status == nil || *i.Status == "" {
		return nil, fastly.ErrMissingStatus
	}

	path := fastly.ToSafeURL("api-security", "v1", "services", *i.ServiceID, "discovered-operations", *i.OperationID)

	resp, err := c.PatchJSON(ctx, path, &DiscoveredOperationStatusUpdate{Status: i.Status}, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var op *DiscoveredOperation
	if err := json.NewDecoder(resp.Body).Decode(&op); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return op, nil
}
