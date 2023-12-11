package fastly

import (
	"fmt"
	"time"
)

// Service represents a server response from the Fastly API.
type Service struct {
	ActiveVersion *int       `mapstructure:"version"`
	Comment       *string    `mapstructure:"comment"`
	CreatedAt     *time.Time `mapstructure:"created_at"`
	CustomerID    *string    `mapstructure:"customer_id"`
	DeletedAt     *time.Time `mapstructure:"deleted_at"`
	ServiceID     *string    `mapstructure:"id"`
	Name          *string    `mapstructure:"name"`
	Type          *string    `mapstructure:"type"`
	UpdatedAt     *time.Time `mapstructure:"updated_at"`
	Versions      []*Version `mapstructure:"versions"`
}

// ServiceDetail represents a server response from the Fastly API.
type ServiceDetail struct {
	ActiveVersion *Version   `mapstructure:"active_version"`
	Comment       *string    `mapstructure:"comment"`
	CreatedAt     *time.Time `mapstructure:"created_at"`
	CustomerID    *string    `mapstructure:"customer_id"`
	DeletedAt     *time.Time `mapstructure:"deleted_at"`
	ServiceID     *string    `mapstructure:"id"`
	Name          *string    `mapstructure:"name"`
	Type          *string    `mapstructure:"type"`
	UpdatedAt     *time.Time `mapstructure:"updated_at"`
	Version       *Version   `mapstructure:"version"`
	Versions      []*Version `mapstructure:"versions"`
}

// ServiceDomain represents a server response from the Fastly API.
type ServiceDomain struct {
	Comment        *string    `mapstructure:"comment"`
	CreatedAt      *time.Time `mapstructure:"created_at"`
	DeletedAt      *time.Time `mapstructure:"deleted_at"`
	Locked         *bool      `mapstructure:"locked"`
	Name           *string    `mapstructure:"name"`
	ServiceID      *string    `mapstructure:"service_id"`
	ServiceVersion *int64     `mapstructure:"version"`
	UpdatedAt      *time.Time `mapstructure:"updated_at"`
}

// ServiceDomainsList represents a list of service domains.
type ServiceDomainsList []*ServiceDomain

// GetServicesInput is used as input to the GetServices function.
type GetServicesInput struct {
	// Direction is the direction in which to sort results.
	Direction *string
	// Page is the current page.
	Page *int
	// PerPage is the number of records per page.
	PerPage *int
	// Sort is the field on which to sort.
	Sort *string
}

// GetServices returns a ListPaginator for paginating through the resources.
func (c *Client) GetServices(i *GetServicesInput) *ListPaginator[Service] {
	input := ListOpts{}
	if i.Direction != nil {
		input.Direction = *i.Direction
	}
	if i.Sort != nil {
		input.Sort = *i.Sort
	}
	if i.Page != nil {
		input.Page = *i.Page
	}
	if i.PerPage != nil {
		input.PerPage = *i.PerPage
	}
	return NewPaginator[Service](c, input, "/service")
}

// ListServicesInput is used as input to the ListServices function.
type ListServicesInput struct {
	// Direction is the direction in which to sort results.
	Direction *string
	// Sort is the field on which to sort.
	Sort *string
}

// ListServices retrieves all resources. Not suitable for large collections.
func (c *Client) ListServices(i *ListServicesInput) ([]*Service, error) {
	p := c.GetServices(&GetServicesInput{
		Direction: i.Direction,
		Sort:      i.Sort,
	})
	var results []*Service
	for p.HasNext() {
		data, err := p.GetNext()
		if err != nil {
			return nil, fmt.Errorf("failed to get next page (remaining: %d): %s", p.Remaining(), err)
		}
		results = append(results, data...)
	}
	return results, nil
}

// CreateServiceInput is used as input to the CreateService function.
type CreateServiceInput struct {
	// Comment is a freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// Name is the name of the service.
	Name *string `url:"name,omitempty"`
	// Type is the type of this service (vcl, wasm).
	Type *string `url:"type,omitempty"`
}

// CreateService creates a new resource.
func (c *Client) CreateService(i *CreateServiceInput) (*Service, error) {
	resp, err := c.PostForm("/service", i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Service
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetServiceInput is used as input to the GetService function.
type GetServiceInput struct {
	// ServiceID is an alphanumeric string identifying the service (required).
	ServiceID string
}

// GetService retrieves the specified resource.
//
// If no service exists for the given id, the API returns a 400 response not 404.
func (c *Client) GetService(i *GetServiceInput) (*Service, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s", i.ServiceID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Service
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}

	// NOTE: GET /service/:service_id endpoint does not return the "version" field
	// unlike other GET service endpoints (/service, /service/:service_id/details).
	// Therefore, ActiveVersion is always zero value when GetService is called.
	// We work around this by manually finding the active version number from the
	// "versions" array in the returned JSON response.
	for i := range s.Versions {
		if *s.Versions[i].Active {
			s.ActiveVersion = s.Versions[i].Number
			break
		}
	}

	return s, nil
}

// GetServiceDetails retrieves the specified resource.
//
// If no service exists for the given id, the API returns a 400 response not 404.
func (c *Client) GetServiceDetails(i *GetServiceInput) (*ServiceDetail, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/details", i.ServiceID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *ServiceDetail
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}

	return s, nil
}

// UpdateServiceInput is used as input to the UpdateService function.
type UpdateServiceInput struct {
	// Comment is a freeform descriptive note.
	Comment *string `url:"comment,omitempty"`
	// Name is the name of the service.
	Name *string `url:"name,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// UpdateService updates the specified resource.
func (c *Client) UpdateService(i *UpdateServiceInput) (*Service, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s", i.ServiceID)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Service
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteServiceInput is used as input to the DeleteService function.
type DeleteServiceInput struct {
	// ServiceID is an alphanumeric string identifying the service (required).
	ServiceID string
}

// DeleteService deletes the specified resource.
func (c *Client) DeleteService(i *DeleteServiceInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s", i.ServiceID)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}

// SearchServiceInput is used as input to the SearchService function.
type SearchServiceInput struct {
	// Name is the name of the service (required).
	Name string
}

// SearchService retrieves the specified resource.
//
// If no service exists by that name, the API returns a 400 response not a 404.
func (c *Client) SearchService(i *SearchServiceInput) (*Service, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}

	resp, err := c.Get("/service/search", &RequestOptions{
		Params: map[string]string{
			"name": i.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Service
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}

	return s, nil
}

// ListServiceDomainInput is the input parameter to the ListServiceDomains
// function.
type ListServiceDomainInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// ListServiceDomains retrieves all resources.
func (c *Client) ListServiceDomains(i *ListServiceDomainInput) (ServiceDomainsList, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := fmt.Sprintf("/service/%s/domain", i.ServiceID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ds ServiceDomainsList

	if err := decodeBodyMap(resp.Body, &ds); err != nil {
		return nil, err
	}

	return ds, nil
}
