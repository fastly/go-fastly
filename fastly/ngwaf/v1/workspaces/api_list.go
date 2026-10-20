package workspaces

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
	// Limit is the maximum number of results to return.
	Limit *int
	// Mode filter results based on mode.
	Mode *string
}

// maxPageSize is the largest page size the API accepts.
const maxPageSize = 1000

// List retrieves all workspaces, automatically paginating through every page.
func List(ctx context.Context, c *fastly.Client, i *ListInput) ([]Workspace, error) {
	out := []Workspace{}
	for page := 1; ; page++ {
		ws, err := listPage(ctx, c, i, page)
		if err != nil {
			return nil, err
		}
		out = append(out, ws.Data...)

		// Stop on an empty page, once Limit results have been
		// collected, or once meta.total is reached.
		if i.Limit != nil && *i.Limit > 0 && len(out) >= *i.Limit {
			out = out[:*i.Limit]
			break
		}
		if len(ws.Data) == 0 || len(out) >= ws.Meta.Total {
			break
		}
	}

	return out, nil
}

// listPage retrieves a single page of workspaces.
func listPage(ctx context.Context, c *fastly.Client, i *ListInput, page int) (*Workspaces, error) {
	// Request the largest page the API allows, unless Limit asks for fewer.
	perPage := maxPageSize
	if i.Limit != nil && *i.Limit > 0 && *i.Limit < perPage {
		perPage = *i.Limit
	}

	requestOptions := fastly.CreateRequestOptions()
	requestOptions.Params["limit"] = strconv.Itoa(perPage)
	if i.Mode != nil {
		requestOptions.Params["mode"] = *i.Mode
	}
	if page > 1 {
		requestOptions.Params["page"] = strconv.Itoa(page)
	}

	resp, err := c.Get(ctx, "/ngwaf/v1/workspaces", requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ws *Workspaces
	if err := json.NewDecoder(resp.Body).Decode(&ws); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return ws, nil
}
