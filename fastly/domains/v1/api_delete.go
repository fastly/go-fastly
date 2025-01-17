package v1

import (
	"net/http"

	"github.com/fastly/go-fastly/v9/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// DomainID of the domain to delete (required).
	DomainID *string
}

// Delete deletes the specified domain.
func Delete(c *fastly.Client, i *DeleteInput) error {
	if i.DomainID == nil {
		return fastly.ErrMissingDomainID
	}

	path := fastly.ToSafeURL("domains", "v1", *i.DomainID)

	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
