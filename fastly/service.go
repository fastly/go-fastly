package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/peterhellberg/link"
)

// Service represents a single service for the Fastly account.
type Service struct {
	ID            string     `mapstructure:"id"`
	Name          string     `mapstructure:"name"`
	Type          string     `mapstructure:"type"`
	Comment       string     `mapstructure:"comment"`
	CustomerID    string     `mapstructure:"customer_id"`
	CreatedAt     *time.Time `mapstructure:"created_at"`
	UpdatedAt     *time.Time `mapstructure:"updated_at"`
	DeletedAt     *time.Time `mapstructure:"deleted_at"`
	ActiveVersion uint       `mapstructure:"version"`
	Versions      []*Version `mapstructure:"versions"`
}

type ServiceDetail struct {
	ID            string     `mapstructure:"id"`
	Name          string     `mapstructure:"name"`
	Type          string     `mapstructure:"type"`
	Comment       string     `mapstructure:"comment"`
	CustomerID    string     `mapstructure:"customer_id"`
	ActiveVersion Version    `mapstructure:"active_version"`
	Version       Version    `mapstructure:"version"`
	Versions      []*Version `mapstructure:"versions"`
	CreatedAt     *time.Time `mapstructure:"created_at"`
	UpdatedAt     *time.Time `mapstructure:"updated_at"`
	DeletedAt     *time.Time `mapstructure:"deleted_at"`
}

type ServiceDomain struct {
	Locked         bool       `mapstructure:"locked"`
	Name           string     `mapstructure:"name"`
	DeletedAt      *time.Time `mapstructure:"deleted_at"`
	ServiceID      string     `mapstructure:"service_id"`
	ServiceVersion int64      `mapstructure:"version"`
	CreatedAt      *time.Time `mapstructure:"created_at"`
	Comment        string     `mapstructure:"comment"`
	UpdatedAt      *time.Time `mapstructure:"updated_at"`
}
type ServiceDomainsList []*ServiceDomain

// servicesByName is a sortable list of services.
type servicesByName []*Service

// Len, Swap, and Less implement the sortable interface.
func (s servicesByName) Len() int      { return len(s) }
func (s servicesByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s servicesByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListServicesInput is used as input to the ListServices function.
type ListServicesInput struct {
	Direction string
	PerPage   int
	Page      int
	Sort      string
}

// ListServices returns the full list of services for the current account.
func (c *Client) ListServices(i *ListServicesInput) ([]*Service, error) {
	resp, err := c.Get("/service", nil)
	if err != nil {
		return nil, err
	}

	var s []*Service
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	sort.Stable(servicesByName(s))
	return s, nil
}

type ListServicesPaginator struct {
	consumed    bool
	CurrentPage int
	NextPage    int
	LastPage    int
	client      *Client
	options     *ListServicesInput
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

	if i.Page <= 0 && p.CurrentPage == 0 {
		p.CurrentPage = 1
	} else {
		p.CurrentPage = p.CurrentPage + 1
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
	Name    string `url:"name,omitempty"`
	Type    string `url:"type,omitempty"`
	Comment string `url:"comment,omitempty"`
}

// CreateService creates a new service with the given information.
func (c *Client) CreateService(i *CreateServiceInput) (*Service, error) {
	resp, err := c.PostForm("/service", i, nil)
	if err != nil {
		return nil, err
	}

	var s *Service
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetServiceInput is used as input to the GetService function.
type GetServiceInput struct {
	ID string
}

// GetService retrieves the service information for the service with the given
// id. If no service exists for the given id, the API returns a 400 response
// (not a 404).
func (c *Client) GetService(i *GetServiceInput) (*Service, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/service/%s", i.ID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

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
			s.ActiveVersion = uint(s.Versions[i].Number)
			break
		}
	}

	return s, nil
}

// GetService retrieves the details for the service with the given id. If no
// service exists for the given id, the API returns a 400 response (not a 404).
func (c *Client) GetServiceDetails(i *GetServiceInput) (*ServiceDetail, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/service/%s/details", i.ID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var s *ServiceDetail
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}

	return s, nil
}

// UpdateServiceInput is used as input to the UpdateService function.
type UpdateServiceInput struct {
	ServiceID string

	Name    *string `url:"name,omitempty"`
	Comment *string `url:"comment,omitempty"`
}

// UpdateService updates the service with the given input.
func (c *Client) UpdateService(i *UpdateServiceInput) (*Service, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.Name == nil && i.Comment == nil {
		return nil, ErrMissingOptionalNameComment
	}

	if i.Name != nil && *i.Name == "" {
		return nil, ErrMissingNameValue
	}

	path := fmt.Sprintf("/service/%s", i.ServiceID)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var s *Service
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteServiceInput is used as input to the DeleteService function.
type DeleteServiceInput struct {
	ID string
}

// DeleteService updates the service with the given input.
func (c *Client) DeleteService(i *DeleteServiceInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/service/%s", i.ID)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

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
	Name string
}

// SearchService gets a specific service by name. If no service exists by that
// name, the API returns a 400 response (not a 404).
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

	var s *Service
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}

	return s, nil
}

type ListServiceDomainInput struct {
	ID string
}

// ListServiceDomains lists all domains associated with a given service
func (c *Client) ListServiceDomains(i *ListServiceDomainInput) (ServiceDomainsList, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}
	path := fmt.Sprintf("/service/%s/domain", i.ID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var ds ServiceDomainsList

	if err := decodeBodyMap(resp.Body, &ds); err != nil {
		return nil, err
	}

	return ds, nil
}
