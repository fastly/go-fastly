package computeacls

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v12/fastly"
)

// ListEntriesInput specifies the information needed for the ListEntries() function to perform
// the operation.
type ListEntriesInput struct {
	// ComputeACLID is an ACL Identifier (required).
	ComputeACLID *string
	// Cursor is used for paginating through results.
	Cursor *string
	// Limit is the maximum number of entries included the response.
	Limit *int
}

// ListEntries.
func ListEntries(ctx context.Context, c *fastly.Client, i *ListEntriesInput) (*ComputeACLEntries, error) {
	if i.ComputeACLID == nil {
		return nil, fastly.ErrMissingComputeACLID
	}

	requestOptions := fastly.CreateRequestOptions()
	if i.Cursor != nil {
		requestOptions.Params["cursor"] = *i.Cursor
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}

	path := fastly.ToSafeURL("resources", "acls", *i.ComputeACLID, "entries")

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var entries *ComputeACLEntries
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return entries, nil
}
