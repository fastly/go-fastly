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
	// Limit how many results are returned per page.
	Limit *int
	// Mode filter results based on mode.
	Mode *string
}

// List retrieves all workspaces, automatically paginating through every page.
func List(ctx context.Context, c *fastly.Client, i *ListInput) ([]Workspace, error) {
	out := []Workspace{}
	for page := 1; ; page++ {
		ws, err := listPage(ctx, c, i, page)
		if err != nil {
			return nil, err
		}
		out = append(out, ws.Data...)

		// Stop on an empty page or once meta.total is reached.
		if len(ws.Data) == 0 || len(out) >= ws.Meta.Total {
			break
		}
	}

	return out, nil
}

// listPage retrieves a single page of workspaces.
func listPage(ctx context.Context, c *fastly.Client, i *ListInput, page int) (*Workspaces, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
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
