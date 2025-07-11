package v1

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
	// Cursor is the cursor value from the next_cursor field of a previous
	// response, used to retrieve the next page. To request the first page, this
	// should be an empty string or nil.
	Cursor *string
	// FQDN filters results by the FQDN using a fuzzy/partial match (optional).
	FQDN *string
	// Limit is the maximum number of results to return (optional).
	Limit *int
	// ServiceID filter results based on a service_id (optional).
	ServiceID *string
	// Sort is the order in which to list the results (optional).
	Sort *string
}

// List retrieves a list of domains, with optional filtering and pagination.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Collection, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.Cursor != nil {
		requestOptions.Params["cursor"] = *i.Cursor
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.FQDN != nil {
		requestOptions.Params["fqdn"] = *i.FQDN
	}
	if i.ServiceID != nil {
		requestOptions.Params["service_id"] = *i.ServiceID
	}
	if i.Sort != nil {
		requestOptions.Params["sort"] = *i.Sort
	}

	resp, err := c.Get(ctx, "/domains/v1", requestOptions)
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
