package redactions

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListInput specifies the information needed for the List() function
// to perform the operation.
type ListInput struct {
	// Limit how many results are returned (1-100). Defaults to 100.
	Limit *int
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// maxLimit is the largest page size the API accepts. The endpoint has no
// paging, so it is also the most redactions one call can return.
const maxLimit = 100

// List retrieves a list of redactions, up to 100 by default (the API's
// own default is 10).
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Redactions, error) {
	requestOptions := fastly.CreateRequestOptions()
	limit := maxLimit
	if i.Limit != nil {
		limit = *i.Limit
	}
	requestOptions.Params["limit"] = strconv.Itoa(limit)

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "redactions")

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var redactions *Redactions
	if err := json.NewDecoder(resp.Body).Decode(&redactions); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return redactions, nil
}
