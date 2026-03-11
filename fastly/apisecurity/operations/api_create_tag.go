package operations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v13/fastly"
)

// CreateTagInput specifies the information needed for the CreateTag() function
// to perform the operation.
type CreateTagInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string `json:"-"`
	// Name is the name of the operation tag (required).
	Name *string `json:"name"`
	// Description describes the tag.
	Description *string `json:"description,omitempty"`
}

// CreateTag creates a new operation tag associated with a service.
func CreateTag(ctx context.Context, c *fastly.Client, i *CreateTagInput) (*OperationTag, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}

	path := fastly.ToSafeURL("api-security", "v1", "services", *i.ServiceID, "tags")

	resp, err := c.PostJSON(ctx, path, i, fastly.CreateRequestOptions())
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
