package paths

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v17/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// PathID of the path to delete (required).
	PathID *string
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string
}

// Delete deletes the specified path.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.RoutingConfigID == nil {
		return fastly.ErrMissingRoutingConfigID
	}
	if i.PathID == nil {
		return fastly.ErrMissingPathID
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "paths", *i.PathID)

	resp, err := c.Delete(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
