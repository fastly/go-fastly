package draft

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v17/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string
}

// Delete discards the draft version of the specified routing config,
// reverting it back to match the active version.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.RoutingConfigID == nil || *i.RoutingConfigID == "" {
		return fastly.ErrMissingRoutingConfigID
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "draft")

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
