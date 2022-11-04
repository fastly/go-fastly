package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/peterhellberg/link"
)

// Service represents a server response from the Fastly API.
type Service struct {
	ActiveVersion int        `mapstructure:"version"`
	Comment       string     `mapstructure:"comment"`
	CreatedAt     *time.Time `mapstructure:"created_at"`
	CustomerID    string     `mapstructure:"customer_id"`
	DeletedAt     *time.Time `mapstructure:"deleted_at"`
	ID            string     `mapstructure:"id"`
	Name          string     `mapstructure:"name"`
	Type          string     `mapstructure:"type"`
	UpdatedAt     *time.Time `mapstructure:"updated_at"`
	Versions      []*Version `mapstructure:"versions"`
}

// ServiceDetail represents a server response from the Fastly API.
type ServiceDetail struct {
	ActiveVersion Version    `mapstructure:"active_version"`
	Comment       string     `mapstructure:"comment"`
	CreatedAt     *time.Time `mapstructure:"created_at"`
	CustomerID    string     `mapstructure:"customer_id"`
	DeletedAt     *time.Time `mapstructure:"deleted_at"`
	ID            string     `mapstructure:"id"`
	Name          string     `mapstructure:"name"`
	Type          string     `mapstructure:"type"`
	UpdatedAt     *time.Time `mapstructure:"updated_at"`
	Version       Version    `mapstructure:"version"`
	Versions      []*Version `mapstructure:"versions"`
}

// ServiceDomain represents a server response from the Fastly API.
type ServiceDomain struct {
	Comment        string     `mapstructure:"comment"`
	CreatedAt      *time.Time `mapstructure:"created_at"`
	DeletedAt      *time.Time `mapstructure:"deleted_at"`
	Locked         bool       `mapstructure:"locked"`
	Name           string     `mapstructure:"name"`
	ServiceID      string     `mapstructure:"service_id"`
	ServiceVersion int64      `mapstructure:"version"`
	UpdatedAt      *time.Time `mapstructure:"updated_at"`
}

// ServiceDomainsList represents a list of service domains.
type ServiceDomainsList []*ServiceDomain

// servicesByName is a sortable list of services.
type servicesByName []*Service

// Len implement the sortable interface.
func (s servicesByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s servicesByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s servicesByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListServicesInput is used as input to the ListServices function.
type ListServicesInput struct {
	// Direction is the direction in which to sort results.
	Direction string
	// Page is the current page.
	Page int
	// PerPage is the number of records per page.
	PerPage int
	// Sort is the field on which to sort.
	Sort string
}

func (l *ListServicesInput) formatFilters() map[string]string {
	m := make(map[string]string)

	if l.Direction != "" {
		m["direction"] = l.Direction
	}
	if l.Page != 0 {
		m["page"] = strconv.Itoa(l.Page)
	}
	if l.PerPage != 0 {
		m["per_page"] = strconv.Itoa(l.PerPage)
	}
	if l.Sort != "" {
		m["sort"] = l.Sort
	}

	return m
}

// ListServices retrieves all resources.
func (c *Client) ListServices(i *ListServicesInput) ([]*Service, error) {
	ro := new(RequestOptions)
	ro.Params = i.formatFilters()

	resp, err := c.Get("/service", ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s []*Service
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	sort.Stable(servicesByName(s))
	return s, nil
}

// ListServicesPaginator implements the PaginatorServices interface.
type ListServicesPaginator struct {
	CurrentPage int
	LastPage    int
	NextPage    int

	// Private
	client   *Client
	consumed bool
	options  *ListServicesInput
}

// HasNext returns a boolean indicating whether more pages are available
func (p *ListServicesPaginator) HasNext() bool {
	return !p.consumed || p.Remaining() != 0
}

// Remaining returns the remaining page count
func (p *ListServicesPaginator) Remaining() int {
	if p.LastPage == 0 {
		return 0
	}
	return p.LastPage - p.CurrentPage
}

// GetNext retrieves data in the next page
func (p *ListServicesPaginator) GetNext() ([]*Service, error) {
	return p.client.listServicesWithPage(p.options, p)
}

// NewListServicesPaginator returns a new paginator
func (c *Client) NewListServicesPaginator(i *ListServicesInput) PaginatorServices {
	return &ListServicesPaginator{
		client:  c,
		options: i,
	}
}

// listServicesWithPage return a list of services
func (c *Client) listServicesWithPage(i *ListServicesInput, p *ListServicesPaginator) ([]*Service, error) {
	var perPage int
	const maxPerPage = 100
	if i.PerPage <= 0 {
		perPage = maxPerPage
	} else {
		perPage = i.PerPage
	}

	// page is not specified, fetch from the beginning
	if i.Page <= 0 && p.CurrentPage == 0 {
		p.CurrentPage = 1
	} else {
		// page is specified, fetch from a given page
		if !p.consumed {
			p.CurrentPage = i.Page
		} else {
			p.CurrentPage = p.CurrentPage + 1
		}
	}

	requestOptions := &RequestOptions{
		Params: map[string]string{
			"per_page": strconv.Itoa(perPage),
			"page":     strconv.Itoa(p.CurrentPage),
		},
	}

	if i.Direction != "" {
		requestOptions.Params["direction"] = i.Direction
	}
	if i.Sort != "" {
		requestOptions.Params["sort"] = i.Sort
	}

	resp, err := c.Get("/service", requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	for _, l := range link.ParseResponse(resp) {
		// indicates the Link response header contained the next page instruction
		if l.Rel == "next" {
			u, _ := url.Parse(l.URI)
			query := u.Query()
			p.NextPage, _ = strconv.Atoi(query["page"][0])
		}
		// indicates the Link response header contained the last page instruction
		if l.Rel == "last" {
			u, _ := url.Parse(l.URI)
			query := u.Query()
			p.LastPage, _ = strconv.Atoi(query["page"][0])
		}
	}

	p.consumed = true

	var s []*Service
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}

	sort.Stable(servicesByName(s))

	return s, nil
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
	// ID is an alphanumeric string identifying the service (required).
	ID string
}

// GetService retrieves the specified resource.
//
// If no service exists for the given id, the API returns a 400 response not 404.
func (c *Client) GetService(i *GetServiceInput) (*Service, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/service/%s", i.ID)
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
		if s.Versions[i].Active {
			s.ActiveVersion = int(s.Versions[i].Number)
			break
		}
	}

	return s, nil
}

// GetServiceDetails retrieves the specified resource.
//
// If no service exists for the given id, the API returns a 400 response not 404.
func (c *Client) GetServiceDetails(i *GetServiceInput) (*ServiceDetail, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/service/%s/details", i.ID)
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
	// ID is an alphanumeric string identifying the service (required).
	ID string
}

// DeleteService deletes the specified resource.
func (c *Client) DeleteService(i *DeleteServiceInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/service/%s", i.ID)
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
