package virtualpatches

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v10/fastly"
)

// ListInput specifies the information needed for the List() function to perform
// the operation.
type ListInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Limit how many results are returned (optional).
	Limit *int
	// Page number of the collection to request (optional).
	Page *int
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// List retrieves a list of virtual patches, with optional filtering and pagination.
func List(c *fastly.Client, i *ListInput) (*VirtualPatches, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	requestOptions := fastly.CreateRequestOptions(i.Context)
	if i.Page != nil {
		requestOptions.Params["page"] = strconv.Itoa(*i.Page)
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "virtual-patches")

	resp, err := c.Get(path, requestOptions)
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
