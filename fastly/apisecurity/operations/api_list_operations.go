package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/fastly/go-fastly/v13/fastly"
)

// ListOperationsInput specifies the information needed for the ListOperations()
// function to perform the operation.
type ListOperationsInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string
	// TagID filters operations by tag ID.
	TagID *string
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

// ListOperations lists all operations associated with a service.
func ListOperations(ctx context.Context, c *fastly.Client, i *ListOperationsInput) (*Operations, error) {
	if i.ServiceID == nil {
		return nil, fastly.ErrMissingServiceID
	}

	opts := fastly.CreateRequestOptions()

	if i.TagID != nil && *i.TagID != "" {
		opts.Params["tag_id"] = *i.TagID
	}

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

	path := fastly.ToSafeURL("api-security", "v1", "services", *i.ServiceID, "operations")

	resp, err := c.Get(ctx, path, opts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ops *Operations
	if err := json.NewDecoder(resp.Body).Decode(&ops); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return ops, nil
}
