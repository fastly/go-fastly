package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v13/fastly"
)

// ListTagsInput specifies the information needed for the ListTags() function to
// perform the operation.
type ListTagsInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string
	// Limit is the maximum number of tags to return per page.
	Limit *int
	// Page is the page number to return.
	Page *int
}

// ListTags lists all operation tags associated with a service.
func ListTags(ctx context.Context, c *fastly.Client, i *ListTagsInput) (*OperationTags, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}

	opts := fastly.CreateRequestOptions()
	if i.Limit != nil {
		opts.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Page != nil {
		opts.Params["page"] = strconv.Itoa(*i.Page)
	}

	path := fastly.ToSafeURL("api-security", "v1", "services", *i.ServiceID, "tags")

	resp, err := c.Get(ctx, path, opts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tags *OperationTags
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return tags, nil
}
