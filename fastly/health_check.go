package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// HealthCheck represents a health check response from the Fastly API.
type HealthCheck struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name             string     `mapstructure:"name"`
	Comment          string     `mapstructure:"comment"`
	Method           string     `mapstructure:"method"`
	Host             string     `mapstructure:"host"`
	Path             string     `mapstructure:"path"`
	HTTPVersion      string     `mapstructure:"http_version"`
	Timeout          uint       `mapstructure:"timeout"`
	CheckInterval    uint       `mapstructure:"check_interval"`
	ExpectedResponse uint       `mapstructure:"expected_response"`
	Window           uint       `mapstructure:"window"`
	Threshold        uint       `mapstructure:"threshold"`
	Initial          uint       `mapstructure:"initial"`
	CreatedAt        *time.Time `mapstructure:"created_at"`
	UpdatedAt        *time.Time `mapstructure:"updated_at"`
	DeletedAt        *time.Time `mapstructure:"deleted_at"`
}

// healthChecksByName is a sortable list of health checks.
type healthChecksByName []*HealthCheck

// Len, Swap, and Less implement the sortable interface.
func (s healthChecksByName) Len() int      { return len(s) }
func (s healthChecksByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
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

// ListHealthChecks returns the list of health checks for the configuration
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name             string `form:"name,omitempty"`
	Comment          string `form:"comment,omitempty"`
	Method           string `form:"method,omitempty"`
	Host             string `form:"host,omitempty"`
	Path             string `form:"path,omitempty"`
	HTTPVersion      string `form:"http_version,omitempty"`
	Timeout          uint   `form:"timeout,omitempty"`
	CheckInterval    uint   `form:"check_interval,omitempty"`
	ExpectedResponse uint   `form:"expected_response,omitempty"`
	Window           uint   `form:"window,omitempty"`
	Threshold        uint   `form:"threshold,omitempty"`
	Initial          uint   `form:"initial,omitempty"`
}

// CreateHealthCheck creates a new Fastly health check.
func (c *Client) CreateHealthCheck(i *CreateHealthCheckInput) (*HealthCheck, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/healthcheck", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var h *HealthCheck
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// GetHealthCheckInput is used as input to the GetHealthCheck function.
type GetHealthCheckInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the health check to fetch.
	Name string
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

	var h *HealthCheck
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// UpdateHealthCheckInput is used as input to the UpdateHealthCheck function.
type UpdateHealthCheckInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the health check to update.
	Name string

	NewName          *string `form:"name,omitempty"`
	Comment          *string `form:"comment,omitempty"`
	Method           *string `form:"method,omitempty"`
	Host             *string `form:"host,omitempty"`
	Path             *string `form:"path,omitempty"`
	HTTPVersion      *string `form:"http_version,omitempty"`
	Timeout          *uint   `form:"timeout,omitempty"`
	CheckInterval    *uint   `form:"check_interval,omitempty"`
	ExpectedResponse *uint   `form:"expected_response,omitempty"`
	Window           *uint   `form:"window,omitempty"`
	Threshold        *uint   `form:"threshold,omitempty"`
	Initial          *uint   `form:"initial,omitempty"`
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

	path := fmt.Sprintf("/service/%s/version/%d/healthcheck/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var h *HealthCheck
	if err := decodeBodyMap(resp.Body, &h); err != nil {
		return nil, err
	}
	return h, nil
}

// DeleteHealthCheckInput is the input parameter to DeleteHealthCheck.
type DeleteHealthCheckInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the health check to delete (required).
	Name string
}

// DeleteHealthCheck deletes the given health check.
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

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
