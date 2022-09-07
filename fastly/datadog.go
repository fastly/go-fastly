package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// Datadog represents a Datadog response from the Fastly API.
type Datadog struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListDatadog returns the list of Datadog for the configuration version.
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `url:"name,omitempty"`
	Token             string `url:"token,omitempty"`
	Region            string `url:"region,omitempty"`
	Format            string `url:"format,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	Placement         string `url:"placement,omitempty"`
}

// CreateDatadog creates a new Datadog logging endpoint on a Fastly service version.
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Datadog to fetch.
	Name string
}

// GetDatadog gets the Datadog configuration with the given parameters.
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Datadog to update.
	Name string

	NewName           *string `url:"name,omitempty"`
	Token             *string `url:"token,omitempty"`
	Region            *string `url:"region,omitempty"`
	Format            *string `url:"format,omitempty"`
	FormatVersion     *uint   `url:"format_version,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	Placement         *string `url:"placement,omitempty"`
}

// UpdateDatadog updates a Datadog logging endpoint on a Fastly service version.
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the Datadog to delete (required).
	Name string
}

// DeleteDatadog deletes a Datadog logging endpoint on a Fastly service version.
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
