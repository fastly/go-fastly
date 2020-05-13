package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Datadog represents a Datadog response from the Fastly API.
type Datadog struct {
	ServiceID string `mapstructure:"service_id"`
	Version   int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Token             string     `mapstructure:"token"`
	Region            string     `mapstructure:"region"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Placement         string     `mapstructure:"placement"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// datadogByName is a sortable list of Datadog.
type datadogByName []*Datadog

// Len, Swap, and Less implement the sortable interface.
func (s datadogByName) Len() int      { return len(s) }
func (s datadogByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s datadogByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListDatadogInput is used as input to the ListDatadog function.
type ListDatadogInput struct {
	// Service is the ID of the service (required).
	Service string

	// Version is the specific configuration version (required).
	Version int
}

// ListDatadog returns the list of Datadog for the configuration version.
func (c *Client) ListDatadog(i *ListDatadogInput) ([]*Datadog, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/datadog", i.Service, i.Version)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var d []*Datadog
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	sort.Stable(datadogByName(d))
	return d, nil
}

// CreateDatadogInput is used as input to the CreateDatadog function.
type CreateDatadogInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	Name              *string `form:"name,omitempty"`
	Token             *string `form:"token,omitempty"`
	Region            *string `form:"region,omitempty"`
	Format            *string `form:"format,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// CreateDatadog creates a new Datadog logging endpoint on a Fastly service version.
func (c *Client) CreateDatadog(i *CreateDatadogInput) (*Datadog, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/datadog", i.Service, i.Version)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var d *Datadog
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// GetDatadogInput is used as input to the GetDatadog function.
type GetDatadogInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the Datadog to fetch.
	Name string
}

// GetDatadog gets the Datadog configuration with the given parameters.
func (c *Client) GetDatadog(i *GetDatadogInput) (*Datadog, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/datadog/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var d *Datadog
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// UpdateDatadogInput is used as input to the UpdateDatadog function.
type UpdateDatadogInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the Datadog to update.
	Name string

	NewName           *string `form:"name,omitempty"`
	Token             *string `form:"token,omitempty"`
	Region            *string `form:"region,omitempty"`
	Format            *string `form:"format,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// UpdateDatadog updates a Datadog logging endpoint on a Fastly service version.
func (c *Client) UpdateDatadog(i *UpdateDatadogInput) (*Datadog, error) {
	if i.Service == "" {
		return nil, ErrMissingService
	}

	if i.Version == 0 {
		return nil, ErrMissingVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/datadog/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var d *Datadog
	if err := decodeBodyMap(resp.Body, &d); err != nil {
		return nil, err
	}
	return d, nil
}

// DeleteDatadogInput is the input parameter to DeleteDatadog.
type DeleteDatadogInput struct {
	// Service is the ID of the service. Version is the specific configuration
	// version. Both fields are required.
	Service string
	Version int

	// Name is the name of the Datadog to delete (required).
	Name string
}

// DeleteDatadog deletes a Datadog logging endpoint on a Fastly service version.
func (c *Client) DeleteDatadog(i *DeleteDatadogInput) error {
	if i.Service == "" {
		return ErrMissingService
	}

	if i.Version == 0 {
		return ErrMissingVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/datadog/%s", i.Service, i.Version, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrStatusNotOk
	}
	return nil
}
