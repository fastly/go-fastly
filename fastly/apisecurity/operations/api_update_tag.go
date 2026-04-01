package operations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v14/fastly"
)

// UpdateTagInput specifies the information needed for the UpdateTag() function
// to perform the operation.
type UpdateTagInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string `json:"-"`
	// TagID is the unique identifier of the operation tag (required).
	TagID *string `json:"-"`
	// Name is the updated name of the tag.
	Name *string `json:"name,omitempty"`
	// Description is an updated description of the tag.
	Description *string `json:"description,omitempty"`
}

// UpdateTag partially updates an existing operation tag.
func UpdateTag(ctx context.Context, c *fastly.Client, i *UpdateTagInput) (*OperationTag, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}
	if i.TagID == nil {
		return nil, fastly.ErrMissingID
	}

	path := fastly.ToSafeURL(
		"api-security", "v1", "services", *i.ServiceID, "tags", *i.TagID,
	)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
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
