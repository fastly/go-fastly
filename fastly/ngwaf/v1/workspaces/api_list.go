package workspaces

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v11/fastly"
)

// ListInput specifies the information needed for the List() function to perform
// the operation.
type ListInput struct {
	// Limit how many results are returned.
	Limit *int
	// Mode filter results based on mode.
	Mode *string
	// Page number of the collection to request.
	Page *int
}

// List retrieves a list of workspaces, with optional filtering and pagination.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Workspaces, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Mode != nil {
		requestOptions.Params["mode"] = *i.Mode
	}
	if i.Page != nil {
		requestOptions.Params["page"] = strconv.Itoa(*i.Page)
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
