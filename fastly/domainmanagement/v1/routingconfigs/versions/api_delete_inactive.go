package versions

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v17/fastly"
)

// DeleteInactiveInput specifies the information needed for the
// DeleteInactive() function to perform the operation.
type DeleteInactiveInput struct {
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string
}

// DeleteInactive permanently deletes all inactive versions of the specified
// routing config. This destroys rollback history and cannot be undone.
func DeleteInactive(ctx context.Context, c *fastly.Client, i *DeleteInactiveInput) error {
	if i.RoutingConfigID == nil {
		return fastly.ErrMissingRoutingConfigID
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "versions", "inactive")

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
