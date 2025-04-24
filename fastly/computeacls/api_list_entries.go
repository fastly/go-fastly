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
	// Context is a context.Context object that will be set to the Request's context.
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

	ro := &fastly.RequestOptions{
		Context: i.Context,
		Params:  map[string]string{},
	}
	if i.Cursor != nil {
		ro.Params["cursor"] = *i.Cursor
	}
	if i.Limit != nil {
		ro.Params["limit"] = strconv.Itoa(*i.Limit)
	}

	path := fastly.ToSafeURL("resources", "acls", *i.ComputeACLID, "entries")

	resp, err := c.Get(path, ro)
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
