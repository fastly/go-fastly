package eventmappings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListInput specifies the information needed for the List() function to
// perform the operation.
type ListInput struct {
	// IntegrationID filters results to mappings that reference the given
	// integration ID.
	IntegrationID *string
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

// List retrieves all event mappings matching the given filters,
// automatically paginating through all pages.
func List(ctx context.Context, c *fastly.Client, i *ListInput) ([]EventMapping, error) {
	var (
		out    []EventMapping
		cursor *string
	)
	for {
		page, err := listPage(ctx, c, i, cursor)
		if err != nil {
			return nil, err
		}
		out = append(out, page.Data...)
		if page.Meta.NextCursor == nil || *page.Meta.NextCursor == "" {
			break
		}
		cursor = page.Meta.NextCursor
	}
	return out, nil
}

// listPage retrieves a single page of event mappings.
func listPage(ctx context.Context, c *fastly.Client, i *ListInput, cursor *string) (*Collection, error) {
	path := fastly.ToSafeURL("notifications", "v1", "event-mappings")

	requestOptions := fastly.CreateRequestOptions()
	if cursor != nil {
		requestOptions.Params["cursor"] = *cursor
	}
	if i.IntegrationID != nil {
		requestOptions.Params["integration_id"] = *i.IntegrationID
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

	resp, err := c.GetJSON(ctx, path, requestOptions)
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
