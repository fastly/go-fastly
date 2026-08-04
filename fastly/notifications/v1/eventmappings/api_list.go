package eventmappings

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListInput specifies the information needed for the List() function to
// perform the operation.
type ListInput struct {
	// Cursor is the pagination cursor from the NextCursor field of a
	// previous response, used to retrieve the next page. To request the
	// first page, this should be nil.
	Cursor *string
	// IntegrationID filters results to mappings that reference the given
	// integration ID.
	IntegrationID *string
	// Limit is the maximum number of results to return per page (1-100).
	Limit *int
	// MappingStatus filters results by mapping status: MappingStatusActive
	// or MappingStatusInactive.
	MappingStatus *string
	// Name filters results to mappings whose name contains the given string
	// (case-insensitive).
	Name *string
	// ScopeID filters results to mappings that apply to the given service or
	// workspace ID. Cannot be combined with ScopeType set to
	// ScopeTypeAccount.
	ScopeID *string
	// ScopeType filters results to the given scope type: ScopeTypeAccount,
	// ScopeTypeVCL, ScopeTypeWasm, or ScopeTypeNGWAF.
	ScopeType *string
	// Sort is the order in which to return results by creation date:
	// "created_at" (oldest first) or "-created_at" (newest first). Defaults
	// to "created_at".
	Sort *string
}

// List retrieves the event mappings matching the given filters, with
// cursor-based pagination.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Collection, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.Cursor != nil {
		requestOptions.Params["cursor"] = *i.Cursor
	}
	if i.IntegrationID != nil {
		requestOptions.Params["integration_id"] = *i.IntegrationID
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.MappingStatus != nil {
		requestOptions.Params["mapping_status"] = *i.MappingStatus
	}
	if i.Name != nil {
		requestOptions.Params["name"] = *i.Name
	}
	if i.ScopeID != nil {
		requestOptions.Params["scope_id"] = *i.ScopeID
	}
	if i.ScopeType != nil {
		requestOptions.Params["scope_type"] = *i.ScopeType
	}
	if i.Sort != nil {
		requestOptions.Params["sort"] = *i.Sort
	}

	path := fastly.ToSafeURL("notifications", "v1", "event-mappings")

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cl *Collection
	if err := json.NewDecoder(resp.Body).Decode(&cl); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return cl, nil
}
