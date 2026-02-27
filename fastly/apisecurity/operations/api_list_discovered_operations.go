package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v13/fastly"
)

// ListDiscoveredInput specifies the information needed for the ListDiscovered()
// function to perform the operation.
type ListDiscoveredInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string
	// Status filters discovered operations by status (required).
	Status *string
	// Method filters operations by HTTP method.
	Method *string
	// Domain filters operations by domain (exact match).
	Domain *string
	// Path filters operations by path (exact match).
	Path *string
	// Page is the page number to return.
	Page *int
	// Limit is the maximum number of results per page.
	Limit *int
}

// ListDiscovered lists discovered operations associated with a service.
//
// NOTE: Although the API description says status is optional, the backend may
// return an error when it is omitted. For reliability, the client requires it.
func ListDiscovered(ctx context.Context, c *fastly.Client, i *ListDiscoveredInput) (*DiscoveredOperations, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}
	if i.Status == nil || *i.Status == "" {
		return nil, fastly.ErrMissingStatus
	}

	opts := fastly.CreateRequestOptions()
	opts.Params["status"] = *i.Status

	if i.Method != nil {
		opts.Params["method"] = *i.Method
	}
	if i.Domain != nil {
		opts.Params["domain"] = *i.Domain
	}
	if i.Path != nil {
		opts.Params["path"] = *i.Path
	}
	if i.Page != nil {
		opts.Params["page"] = strconv.Itoa(*i.Page)
	}
	if i.Limit != nil {
		opts.Params["limit"] = strconv.Itoa(*i.Limit)
	}

	path := fastly.ToSafeURL("api-security", "v1", "services", *i.ServiceID, "discovered-operations")

	resp, err := c.Get(ctx, path, opts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ops *DiscoveredOperations
	if err := json.NewDecoder(resp.Body).Decode(&ops); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return ops, nil
}
