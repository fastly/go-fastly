package fastly

import (
	"context"
	"fmt"
	"time"
)

// Service represents a server response from the Fastly API.
type Service struct {
	ActiveVersion *int           `mapstructure:"version"`
	Comment       *string        `mapstructure:"comment"`
	CreatedAt     *time.Time     `mapstructure:"created_at"`
	CustomerID    *string        `mapstructure:"customer_id"`
	DeletedAt     *time.Time     `mapstructure:"deleted_at"`
	ServiceID     *string        `mapstructure:"id"`
	Name          *string        `mapstructure:"name"`
	Type          *string        `mapstructure:"type"`
	UpdatedAt     *time.Time     `mapstructure:"updated_at"`
	Versions      []*Version     `mapstructure:"versions"`
	Environments  []*Environment `mapstructure:"environments"`
}

// ServiceDetail represents a server response from the Fastly API.
type ServiceDetail struct {
	ActiveVersion *Version       `mapstructure:"active_version"`
	Comment       *string        `mapstructure:"comment"`
	CreatedAt     *time.Time     `mapstructure:"created_at"`
	CustomerID    *string        `mapstructure:"customer_id"`
	DeletedAt     *time.Time     `mapstructure:"deleted_at"`
	ServiceID     *string        `mapstructure:"id"`
	Name          *string        `mapstructure:"name"`
	Type          *string        `mapstructure:"type"`
	UpdatedAt     *time.Time     `mapstructure:"updated_at"`
	Version       *Version       `mapstructure:"version"`
	Versions      []*Version     `mapstructure:"versions"`
	Environments  []*Environment `mapstructure:"environments"`
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
func (c *Client) GetServices(ctx context.Context, i *GetServicesInput) *ListPaginator[Service] {
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
	return NewPaginator[Service](ctx, c, input, "/service")
}

// ListServicesInput is used as input to the ListServices function.
type ListServicesInput struct {
	// Direction is the direction in which to sort results.
	Direction *string
	// Sort is the field on which to sort.
	Sort *string
}

// ListServices retrieves all resources. Not suitable for large collections.
func (c *Client) ListServices(ctx context.Context, i *ListServicesInput) ([]*Service, error) {
	p := c.GetServices(ctx, &GetServicesInput{
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
func (c *Client) CreateService(ctx context.Context, i *CreateServiceInput) (*Service, error) {
	resp, err := c.PostForm(ctx, "/service", i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Service
	if err := DecodeBodyMap(resp.Body, &s); err != nil {
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
func (c *Client) GetService(ctx context.Context, i *GetServiceInput) (*Service, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID)

	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Service
	if err := DecodeBodyMap(resp.Body, &s); err != nil {
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
func (c *Client) GetServiceDetails(ctx context.Context, i *GetServiceInput) (*ServiceDetail, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID, "details")

	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *ServiceDetail
	if err := DecodeBodyMap(resp.Body, &s); err != nil {
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
func (c *Client) UpdateService(ctx context.Context, i *UpdateServiceInput) (*Service, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID)

	resp, err := c.PutForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Service
	if err := DecodeBodyMap(resp.Body, &s); err != nil {
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
func (c *Client) DeleteService(ctx context.Context, i *DeleteServiceInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID)

	resp, err := c.Delete(ctx, path, CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := DecodeBodyMap(resp.Body, &r); err != nil {
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
func (c *Client) SearchService(ctx context.Context, i *SearchServiceInput) (*Service, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}

	requestOptions := CreateRequestOptions()
	requestOptions.Params["name"] = i.Name

	resp, err := c.Get(ctx, "/service/search", requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Service
	if err := DecodeBodyMap(resp.Body, &s); err != nil {
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
func (c *Client) ListServiceDomains(ctx context.Context, i *ListServiceDomainInput) (ServiceDomainsList, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	path := ToSafeURL("service", i.ServiceID, "domain")

	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ds ServiceDomainsList

	if err := DecodeBodyMap(resp.Body, &ds); err != nil {
		return nil, err
	}

	return ds, nil
}
