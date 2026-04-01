package operations

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v14/fastly"
)

// DeleteTagInput specifies the information needed for the DeleteTag() function
// to perform the operation.
type DeleteTagInput struct {
	// ServiceID is the unique identifier of the service (required).
	ServiceID *string
	// TagID is the unique identifier of the operation tag (required).
	TagID *string
}

// DeleteTag deletes an existing operation tag associated with a service.
func DeleteTag(ctx context.Context, c *fastly.Client, i *DeleteTagInput) error {
	if i.ServiceID == nil {
		return fastly.ErrMissingServiceID
	}
	if i.TagID == nil {
		return fastly.ErrMissingID
	}

	path := fastly.ToSafeURL(
		"api-security", "v1", "services", *i.ServiceID, "tags", *i.TagID,
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
