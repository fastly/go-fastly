package fastly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// DomainsV1Response is the API response structure.
type DomainsV1Response struct {
	// Data contains the API data.
	Data []DomainsV1Data `json:"data"`
	// Meta contains metadata related to paginating the full dataset.
	Meta DomainsV1Meta `json:"meta"`
}

// DomainsV1Data is a subset of the API response structure containing the
// specific API data itself.
type DomainsV1Data struct {
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt time.Time `json:"created_at"`
	// ID is the domain identifier (UUID).
	DomainID string `json:"id"`
	// FQDN is the fully-qualified domain name of the domain. Read-only
	// after creation.
	FQDN string `json:"fqdn"`
	// ServiceID is the service_id associated with the domain or nil if there
	// is no association.
	ServiceID *string `json:"service_id"`
	// UpdatedAt is the date and time in ISO 8601 format.
	UpdatedAt time.Time `json:"updated_at"`
}

// DomainsV1Meta is a subset of the API response structure containing metadata
// related to paginating the full dataset.
type DomainsV1Meta struct {
	// Limit is how many results are included in this response.
	Limit int `json:"limit"`
	// NextCursor is the cursor value used to retrieve the next page.
	NextCursor string `json:"next_cursor"`
	// Sort is the field used to order the response by.
	Sort string `json:"sort"`
	// Total is the total number of results.
	Total int `json:"total"`
}

// ListDomainsV1Input is used as input to the ListDomainsV1 function.
type ListDomainsV1Input struct {
	// Cursor is the cursor value from the next_cursor field of a previous
	// response, used to retrieve the next page. To request the first page, this
	// should be an empty string or nil.
	Cursor *string
	// FQDN filters results by the FQDN using a fuzzy/partial match (optional).
	FQDN *string
	// Limit is the maximum number of results to return (optional).
	Limit *int
	// ServiceID filter results based on a service_id (optional).
	ServiceID *string
	// Sort is the order in which to list the results (optional).
	Sort *string
}

// ListDomainsV1 retrieves a list of domains, with optional filtering and pagination.
func (c *Client) ListDomainsV1(i *ListDomainsV1Input) (*DomainsV1Response, error) {
	ro := &RequestOptions{
		Params: map[string]string{},
	}
	if i.Cursor != nil {
		ro.Params["cursor"] = *i.Cursor
	}
	if i.Limit != nil {
		ro.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.FQDN != nil {
		ro.Params["fqdn"] = *i.FQDN
	}
	if i.ServiceID != nil {
		ro.Params["service_id"] = *i.ServiceID
	}
	if i.Sort != nil {
		ro.Params["sort"] = *i.Sort
	}

	resp, err := c.Get("/domains/v1", ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dr *DomainsV1Response
	if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return dr, nil
}

// CreateDomainsV1Input is used as input to the CreateDomainsV1 function.
type CreateDomainsV1Input struct {
	// FQDN is the fully-qualified domain name of the domain (required).
	FQDN *string `json:"fqdn"`
	// ServiceID is the service_id associated with the domain or nil if there
	// is no association (optional)
	ServiceID *string `json:"service_id"`
}

// CreateDomainsV1 creates a new domain.
func (c *Client) CreateDomainsV1(i *CreateDomainsV1Input) (*DomainsV1Data, error) {
	resp, err := c.PostJSON("/domains/v1", i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dd *DomainsV1Data
	if err := json.NewDecoder(resp.Body).Decode(&dd); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}
	return dd, nil
}

// GetDomainsV1Input is used as input to the GetDomainsV1 function.
type GetDomainsV1Input struct {
	// DomainID is the domain identifier (required).
	DomainID *string
}

// GetDomainsV1 retrieves a specified domain.
func (c *Client) GetDomainsV1(i *GetDomainsV1Input) (*DomainsV1Data, error) {
	if i.DomainID == nil {
		return nil, ErrMissingDomainID
	}

	path := ToSafeURL("domains", "v1", *i.DomainID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dd *DomainsV1Data
	if err := json.NewDecoder(resp.Body).Decode(&dd); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return dd, nil
}

// UpdateDomainsV1Input is used as input to the UpdateDomainsV1 function.
type UpdateDomainsV1Input struct {
	// DomainID is the domain identifier (required).
	DomainID *string `json:"-"`
	// ServiceID is the service_id associated with the domain or nil if there
	// is no association (optional)
	ServiceID *string `json:"service_id"`
}

// UpdateDomainsV1 updates the specified domain.
func (c *Client) UpdateDomainsV1(i *UpdateDomainsV1Input) (*DomainsV1Data, error) {
	if i.DomainID == nil {
		return nil, ErrMissingDomainID
	}

	path := ToSafeURL("domains", "v1", *i.DomainID)

	resp, err := c.PatchJSON(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dd *DomainsV1Data
	if err := json.NewDecoder(resp.Body).Decode(&dd); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}
	return dd, nil
}

// DeleteDomainsV1Input is used as input to the DeleteDomainsV1 function.
type DeleteDomainsV1Input struct {
	// DomainID of the domain to delete (required).
	DomainID *string
}

// DeleteDomainsV1 deletes the specified domain.
func (c *Client) DeleteDomainsV1(i *DeleteDomainsV1Input) error {
	if i.DomainID == nil {
		return ErrMissingDomainID
	}

	path := ToSafeURL("domains", "v1", *i.DomainID)

	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return NewHTTPError(resp)
	}

	return nil
}
