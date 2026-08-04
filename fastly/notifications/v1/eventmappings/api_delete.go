package eventmappings

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v17/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// MappingID is the ID of the event mapping (required).
	MappingID *string
}

// Delete deletes the specified event mapping.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.MappingID == nil {
		return fastly.ErrMissingMappingID
	}

	path := fastly.ToSafeURL("notifications", "v1", "event-mappings", *i.MappingID)

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
