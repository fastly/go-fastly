package operations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v14/fastly"
)

// DescribeInput specifies the information needed for the Describe() function to
// perform the operation.
type DescribeInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string
	// OperationID is the unique identifier of the operation (required).
	OperationID *string
}

// Describe retrieves a specific operation associated with a service.
func Describe(ctx context.Context, c *fastly.Client, i *DescribeInput) (*Operation, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}
	if i.OperationID == nil {
		return nil, fastly.ErrMissingID
	}

	path := fastly.ToSafeURL(
		"api-security", "v1", "services", *i.ServiceID, "operations", *i.OperationID,
	)

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
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
