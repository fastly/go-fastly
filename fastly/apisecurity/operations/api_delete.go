package operations

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v13/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string
	// OperationID is the unique identifier of the operation (required).
	OperationID *string
}

// Delete deletes an existing operation associated with a service.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.ServiceID == nil {
		return fastly.ErrMissingServiceID
	}
	if i.OperationID == nil {
		return fastly.ErrMissingID
	}

	path := fastly.ToSafeURL(
		"api-security", "v1", "services", *i.ServiceID, "operations", *i.OperationID,
	)

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
