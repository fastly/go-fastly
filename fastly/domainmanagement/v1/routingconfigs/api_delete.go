package routingconfigs

import (
	"context"
	"net/http"
	"strconv"

	"github.com/fastly/go-fastly/v17/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// Force, when true, allows deletion of a routing config that has an active
	// version, bypassing the active-version check.
	Force *bool
	// RoutingConfigID of the routing config to delete (required).
	RoutingConfigID *string
}

// Delete deletes the specified routing config.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.RoutingConfigID == nil {
		return fastly.ErrMissingRoutingConfigID
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID)

	ro := fastly.CreateRequestOptions()
	if i.Force != nil {
		ro.Params["force"] = strconv.FormatBool(*i.Force)
	}

	resp, err := c.Delete(ctx, path, ro)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
