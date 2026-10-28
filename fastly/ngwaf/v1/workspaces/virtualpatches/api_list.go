package virtualpatches

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v18/fastly"
)

// ListInput specifies the information needed for the List() function
// to perform the operation.
type ListInput struct {
	// Limit is the maximum number of results to return.
	Limit *int
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// maxPageSize is the largest page size the API accepts.
const maxPageSize = 1000

// List retrieves all virtual patches, automatically paginating through every
// page.
func List(ctx context.Context, c *fastly.Client, i *ListInput) ([]VirtualPatch, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	out := []VirtualPatch{}
	for page := 1; ; page++ {
		vps, err := listPage(ctx, c, i, page)
		if err != nil {
			return nil, err
		}
		out = append(out, vps.Data...)

		// Stop on an empty page, once Limit results have been
		// collected, or once meta.total is reached.
		if i.Limit != nil && *i.Limit > 0 && len(out) >= *i.Limit {
			out = out[:*i.Limit]
			break
		}
		if len(vps.Data) == 0 || len(out) >= vps.Meta.Total {
			break
		}
	}

	return out, nil
}

// listPage retrieves a single page of virtual patches.
func listPage(ctx context.Context, c *fastly.Client, i *ListInput, page int) (*VirtualPatches, error) {
	// Request the largest page the API allows, unless Limit asks for fewer.
	perPage := maxPageSize
	if i.Limit != nil && *i.Limit > 0 && *i.Limit < perPage {
		perPage = *i.Limit
	}

	requestOptions := fastly.CreateRequestOptions()
	requestOptions.Params["limit"] = strconv.Itoa(perPage)
	if page > 1 {
		requestOptions.Params["page"] = strconv.Itoa(page)
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "virtual-patches")

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vps *VirtualPatches
	if err := json.NewDecoder(resp.Body).Decode(&vps); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return vps, nil
}
