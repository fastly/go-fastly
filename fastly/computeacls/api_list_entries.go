package computeacls

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v10/fastly"
)

// ListEntriesInput specifies the information needed for the ListEntries() function to perform
// the operation.
type ListEntriesInput struct {
	// ComputeACLID is an ACL Identifier (required).
	ComputeACLID *string
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Cursor is used for paginating through results.
	Cursor *string
	// Limit is the maximum number of entries included the response.
	Limit *int
}

// ListEntries
func ListEntries(c *fastly.Client, i *ListEntriesInput) (*ComputeACLEntries, error) {
	if i.ComputeACLID == nil {
		return nil, fastly.ErrMissingComputeACLID
	}

	requestOptions := fastly.CreateRequestOptions(i.Context)
	if i.Cursor != nil {
		requestOptions.Params["cursor"] = *i.Cursor
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}

	path := fastly.ToSafeURL("resources", "acls", *i.ComputeACLID, "entries")

	resp, err := c.Get(path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer fastly.CheckCloseForErr(resp.Body.Close)

	var entries *ComputeACLEntries
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return entries, nil
}
