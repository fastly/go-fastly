package operations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v13/fastly"
)

// DescribeTagInput specifies the information needed for the DescribeTag()
// function to perform the operation.
type DescribeTagInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string
	// TagID is the unique identifier of the operation tag (required).
	TagID *string
}

// DescribeTag retrieves a specific operation tag associated with a service.
func DescribeTag(ctx context.Context, c *fastly.Client, i *DescribeTagInput) (*OperationTag, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}
	if i.TagID == nil {
		return nil, fastly.ErrMissingID
	}

	path := fastly.ToSafeURL(
		"api-security", "v1", "services", *i.ServiceID, "tags", *i.TagID,
	)

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tag *OperationTag
	if err := json.NewDecoder(resp.Body).Decode(&tag); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return tag, nil
}
