package contacts

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v15/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// CustomerID is the alphanumeric identifier of the customer (required).
	CustomerID *string
	// ContactID is the alphanumeric identifier of the contact (required).
	ContactID *string
}

// Delete deletes the specified contact.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.CustomerID == nil {
		return fastly.ErrMissingCustomerID
	}
	if i.ContactID == nil {
		return fastly.ErrMissingContactID
	}

	path := fastly.ToSafeURL("customer", *i.CustomerID, "contacts", *i.ContactID)

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
