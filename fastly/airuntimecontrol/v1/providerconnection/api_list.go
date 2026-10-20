package providerconnection

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v18/fastly"
)

// ListInput specifies the information needed for the List() function to
// perform the operation.
type ListInput struct {
	// IncludeDeleted includes deleted provider connections in the results.
	// Deleted connections have a `deleted_at` value.
	IncludeDeleted *bool
	// Sort is the sort field. Prefix with `-` for descending order (e.g.,
	// `sort=-name`). One of `created_at`, `updated_at`, `name` (each optionally
	// prefixed with `-`).
	Sort *string
}

// List lists the provider connections belonging to your account.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*ProviderConnections, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.IncludeDeleted != nil {
		requestOptions.Params["include_deleted"] = strconv.FormatBool(*i.IncludeDeleted)
	}
	if i.Sort != nil && *i.Sort != "" {
		requestOptions.Params["sort"] = *i.Sort
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "provider-connections")

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pcs *ProviderConnections
	if err := json.NewDecoder(resp.Body).Decode(&pcs); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return pcs, nil
}
