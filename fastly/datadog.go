package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Datadog represents a Datadog response from the Fastly API.
type Datadog struct {
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	Name              string     `mapstructure:"name"`
	Placement         string     `mapstructure:"placement"`
	Region            string     `mapstructure:"region"`
	ResponseCondition string     `mapstructure:"response_condition"`
	ServiceID         string     `mapstructure:"service_id"`
	ServiceVersion    int        `mapstructure:"version"`
	Token             string     `mapstructure:"token"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// datadogByName is a sortable list of Datadog.
type datadogByName []*Datadog

// Len implement the sortable interface.
func (s datadogByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s datadogByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s datadogByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListDatadogInput is used as input to the ListDatadog function.
type ListDatadogInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListDatadog retrieves all resources.
func (c *Client) ListDatadog(i *ListDatadogInput) ([]*Datadog, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/datadog", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d []*Datadog
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	sort.Stable(datadogByName(d))
	return d, nil
}

// CreateDatadogInput is used as input to the CreateDatadog function.
type CreateDatadogInput struct {
	// Format is a Fastly log format string. Must produce valid JSON that Datadog can ingest.
	Format string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion uint `url:"format_version,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement string `url:"placement,omitempty"`
	// Region is the region that log data will be sent to.
	Region string `url:"region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	// Token is the API key from your Datadog account.
	Token string `url:"token,omitempty"`
}

// CreateDatadog creates a new resource.
func (c *Client) CreateDatadog(i *CreateDatadogInput) (*Datadog, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Token == "" {
		return nil, ErrMissingToken
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/datadog", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Datadog
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// GetDatadogInput is used as input to the GetDatadog function.
type GetDatadogInput struct {
	// Name is the name of the Datadog to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetDatadog retrieves the specified resource.
func (c *Client) GetDatadog(i *GetDatadogInput) (*Datadog, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/datadog/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Datadog
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// UpdateDatadogInput is used as input to the UpdateDatadog function.
type UpdateDatadogInput struct {
	// Format is a Fastly log format string. Must produce valid JSON that Datadog can ingest.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *uint `url:"format_version,omitempty"`
	// Name is the name of the Datadog to update.
	Name string
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// Region is the region that log data will be sent to.
	Region *string `url:"region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	// Token is the API key from your Datadog account.
	Token *string `url:"token,omitempty"`
}

// UpdateDatadog updates the specified resource.
func (c *Client) UpdateDatadog(i *UpdateDatadogInput) (*Datadog, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	if i.Token != nil && *i.Token == "" {
		return nil, ErrTokenEmpty
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/datadog/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Datadog
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// DeleteDatadogInput is the input parameter to DeleteDatadog.
type DeleteDatadogInput struct {
	// Name is the name of the Datadog to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteDatadog deletes the specified resource.
func (c *Client) DeleteDatadog(i *DeleteDatadogInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/datadog/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
		return ErrStatusNotOk
	}
	return nil
}
