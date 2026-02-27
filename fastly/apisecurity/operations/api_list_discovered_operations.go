package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/fastly/go-fastly/v13/fastly"
)

// ListDiscoveredInput specifies the information needed for the ListDiscovered()
// function to perform the operation.
type ListDiscoveredInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string
	// Status filters discovered operations by status (required).
	Status *string
	// Method filters operations by HTTP method. If multiple values are provided,
	// they will be sent as a comma-separated string (e.g. "GET,POST").
	Method []string
	// Domain filters operations by domain (exact match). If multiple values are
	// provided, they will be sent as a comma-separated string.
	Domain []string
	// Path filters operations by path (exact match).
	Path *string
	// Page is the page number to return.
	Page *int
	// Limit is the maximum number of results per page.
	Limit *int
}

// ListDiscovered lists discovered operations associated with a service.
//
// The API requires the "status" query parameter to be set for this endpoint.
func ListDiscovered(ctx context.Context, c *fastly.Client, i *ListDiscoveredInput) (*DiscoveredOperations, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}
	if i.Status == nil || *i.Status == "" {
		return nil, fastly.ErrMissingStatus
	}

	opts := fastly.CreateRequestOptions()
	opts.Params["status"] = *i.Status

	if len(i.Method) > 0 {
		opts.Params["method"] = strings.Join(i.Method, ",")
	}
	if len(i.Domain) > 0 {
		opts.Params["domain"] = strings.Join(i.Domain, ",")
	}
	if i.Path != nil && *i.Path != "" {
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
