package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// HealthCheck represents a health check response from the Fastly API.
type HealthCheck struct {
	CheckInterval    uint       `mapstructure:"check_interval"`
	Comment          string     `mapstructure:"comment"`
	CreatedAt        *time.Time `mapstructure:"created_at"`
	DeletedAt        *time.Time `mapstructure:"deleted_at"`
	ExpectedResponse uint       `mapstructure:"expected_response"`
	HTTPVersion      string     `mapstructure:"http_version"`
	Headers          []string   `mapstructure:"headers"`
	Host             string     `mapstructure:"host"`
	Initial          uint       `mapstructure:"initial"`
	Method           string     `mapstructure:"method"`
	Name             string     `mapstructure:"name"`
	Path             string     `mapstructure:"path"`
	ServiceID        string     `mapstructure:"service_id"`
	ServiceVersion   int        `mapstructure:"version"`
	Threshold        uint       `mapstructure:"threshold"`
	Timeout          uint       `mapstructure:"timeout"`
	UpdatedAt        *time.Time `mapstructure:"updated_at"`
	Window           uint       `mapstructure:"window"`
}

// healthChecksByName is a sortable list of health checks.
type healthChecksByName []*HealthCheck

// Len implement the sortable interface.
func (s healthChecksByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s healthChecksByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s healthChecksByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListHealthChecksInput is used as input to the ListHealthChecks function.
type ListHealthChecksInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListHealthChecks retrieves all resources.
// version.
func (c *Client) ListHealthChecks(i *ListHealthChecksInput) ([]*HealthCheck, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/healthcheck", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var hcs []*HealthCheck
	if err := decodeBodyMap(resp.Body, &hcs); err != nil {
		return nil, err
	}
	sort.Stable(healthChecksByName(hcs))
	return hcs, nil
}

// CreateHealthCheckInput is used as input to the CreateHealthCheck function.
type CreateHealthCheckInput struct {
	CheckInterval    *uint    `url:"check_interval,omitempty"`
	Comment          string   `url:"comment,omitempty"`
	ExpectedResponse *uint    `url:"expected_response,omitempty"`
	HTTPVersion      string   `url:"http_version,omitempty"`
	Headers          []string `url:"headers,omitempty"`
	Host             string   `url:"host,omitempty"`
	Initial          *uint    `url:"initial,omitempty"`
	Method           string   `url:"method,omitempty"`
	Name             string   `url:"name,omitempty"`
	Path             string   `url:"path,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Threshold      *uint `url:"threshold,omitempty"`
	Timeout        *uint `url:"timeout,omitempty"`
	Window         *uint `url:"window,omitempty"`
}

// CreateHealthCheck creates a new resource.
func (c *Client) CreateHealthCheck(i *CreateHealthCheckInput) (*HealthCheck, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	ro := new(RequestOptions)
	ro.HealthCheckHeaders = true

	path := fmt.Sprintf("/service/%s/version/%d/healthcheck", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *HealthCheck
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// GetHealthCheckInput is used as input to the GetHealthCheck function.
type GetHealthCheckInput struct {
	// Name is the name of the health check to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetHealthCheck gets the health check configuration with the given parameters.
func (c *Client) GetHealthCheck(i *GetHealthCheckInput) (*HealthCheck, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/healthcheck/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *HealthCheck
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// UpdateHealthCheckInput is used as input to the UpdateHealthCheck function.
type UpdateHealthCheckInput struct {
	CheckInterval    *uint     `url:"check_interval,omitempty"`
	Comment          *string   `url:"comment,omitempty"`
	ExpectedResponse *uint     `url:"expected_response,omitempty"`
	HTTPVersion      *string   `url:"http_version,omitempty"`
	Headers          *[]string `url:"headers,omitempty"`
	Host             *string   `url:"host,omitempty"`
	Initial          *uint     `url:"initial,omitempty"`
	Method           *string   `url:"method,omitempty"`
	// Name is the name of the health check to update.
	Name    string
	NewName *string `url:"name,omitempty"`
	Path    *string `url:"path,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Threshold      *uint `url:"threshold,omitempty"`
	Timeout        *uint `url:"timeout,omitempty"`
	Window         *uint `url:"window,omitempty"`
}

// UpdateHealthCheck updates a specific health check.
func (c *Client) UpdateHealthCheck(i *UpdateHealthCheckInput) (*HealthCheck, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	ro := new(RequestOptions)
	ro.HealthCheckHeaders = true

	path := fmt.Sprintf("/service/%s/version/%d/healthcheck/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var h *HealthCheck
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// DeleteHealthCheckInput is the input parameter to DeleteHealthCheck.
type DeleteHealthCheckInput struct {
	// Name is the name of the health check to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteHealthCheck deletes the specified resource.
func (c *Client) DeleteHealthCheck(i *DeleteHealthCheckInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/healthcheck/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
