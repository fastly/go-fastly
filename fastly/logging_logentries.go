package fastly

import (
	"context"
	"strconv"
	"time"
)

// Logentries represents a logentries response from the Fastly API.
type Logentries struct {
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	Name              *string    `mapstructure:"name"`
	Placement         *string    `mapstructure:"placement"`
	Port              *int       `mapstructure:"port"`
	Region            *string    `mapstructure:"region"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	Token             *string    `mapstructure:"token"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	UseTLS            *bool      `mapstructure:"use_tls"`
}

// ListLogentriesInput is used as input to the ListLogentries function.
type ListLogentriesInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListLogentries retrieves all resources.
func (c *Client) ListLogentries(i *ListLogentriesInput) ([]*Logentries, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "logentries")
	resp, err := c.Get(path, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer CheckCloseForErr(resp.Body.Close)

	var ls []*Logentries
	if err := DecodeBodyMap(resp.Body, &ls); err != nil {
		return nil, err
	}
	return ls, nil
}

// CreateLogentriesInput is used as input to the CreateLogentries function.
type CreateLogentriesInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `url:"-"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// Port is the port number.
	Port *int `url:"port,omitempty"`
	// Region is the region to which to stream logs.
	Region *string `url:"region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Token is token based authentication
	Token *string `url:"token,omitempty"`
	// UseTLS is whether to use TLS (0: do not use, 1: use).
	UseTLS *Compatibool `url:"use_tls,omitempty"`
}

// CreateLogentries creates a new resource.
func (c *Client) CreateLogentries(i *CreateLogentriesInput) (*Logentries, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "logentries")
	resp, err := c.PostForm(path, i, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer CheckCloseForErr(resp.Body.Close)

	var l *Logentries
	if err := DecodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// GetLogentriesInput is used as input to the GetLogentries function.
type GetLogentriesInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Name is the name of the logentries to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetLogentries retrieves the specified resource.
func (c *Client) GetLogentries(i *GetLogentriesInput) (*Logentries, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "logentries", i.Name)
	resp, err := c.Get(path, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer CheckCloseForErr(resp.Body.Close)

	var l *Logentries
	if err := DecodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// UpdateLogentriesInput is used as input to the UpdateLogentries function.
type UpdateLogentriesInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `url:"-"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name of the logentries to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// Port is the port number.
	Port *int `url:"port,omitempty"`
	// Region is the region to which to stream logs.
	Region *string `url:"region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Token is token based authentication
	Token *string `url:"token,omitempty"`
	// UseTLS is whether to use TLS (0: do not use, 1: use).
	UseTLS *Compatibool `url:"use_tls,omitempty"`
}

// UpdateLogentries updates the specified resource.
func (c *Client) UpdateLogentries(i *UpdateLogentriesInput) (*Logentries, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "logentries", i.Name)
	resp, err := c.PutForm(path, i, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer CheckCloseForErr(resp.Body.Close)

	var l *Logentries
	if err := DecodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// DeleteLogentriesInput is the input parameter to DeleteLogentries.
type DeleteLogentriesInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Name is the name of the logentries to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteLogentries deletes the specified resource.
func (c *Client) DeleteLogentries(i *DeleteLogentriesInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "logentries", i.Name)
	resp, err := c.Delete(path, CreateRequestOptions(i.Context))
	if err != nil {
		return err
	}
	defer CheckCloseForErr(resp.Body.Close)

	var r *statusResp
	if err := DecodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
