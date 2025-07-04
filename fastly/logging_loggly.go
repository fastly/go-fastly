package fastly

import (
	"context"
	"strconv"
	"time"
)

// Loggly represents a loggly response from the Fastly API.
type Loggly struct {
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	Name              *string    `mapstructure:"name"`
	Placement         *string    `mapstructure:"placement"`
	ProcessingRegion  *string    `mapstructure:"log_processing_region"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	Token             *string    `mapstructure:"token"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// ListLogglyInput is used as input to the ListLoggly function.
type ListLogglyInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListLoggly retrieves all resources.
func (c *Client) ListLoggly(i *ListLogglyInput) ([]*Loggly, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "loggly")
	resp, err := c.Get(path, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ls []*Loggly
	if err := DecodeBodyMap(resp.Body, &ls); err != nil {
		return nil, err
	}
	return ls, nil
}

// CreateLogglyInput is used as input to the CreateLoggly function.
type CreateLogglyInput struct {
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
	// ProcessingRegion is the region where logs will be processed before streaming to Loggly.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Token is the token to use for authentication.
	Token *string `url:"token,omitempty"`
}

// CreateLoggly creates a new resource.
func (c *Client) CreateLoggly(i *CreateLogglyInput) (*Loggly, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "loggly")
	resp, err := c.PostForm(path, i, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var l *Loggly
	if err := DecodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// GetLogglyInput is used as input to the GetLoggly function.
type GetLogglyInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Name is the name of the loggly to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetLoggly retrieves the specified resource.
func (c *Client) GetLoggly(i *GetLogglyInput) (*Loggly, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "loggly", i.Name)
	resp, err := c.Get(path, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var l *Loggly
	if err := DecodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// UpdateLogglyInput is used as input to the UpdateLoggly function.
type UpdateLogglyInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `url:"-"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// Name is the name of the loggly to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// ProcessingRegion is the region where logs will be processed before streaming to Loggly.
	ProcessingRegion *string `url:"log_processing_region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// Token is the token to use for authentication.
	Token *string `url:"token,omitempty"`
}

// UpdateLoggly updates the specified resource.
func (c *Client) UpdateLoggly(i *UpdateLogglyInput) (*Loggly, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "loggly", i.Name)
	resp, err := c.PutForm(path, i, CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var l *Loggly
	if err := DecodeBodyMap(resp.Body, &l); err != nil {
		return nil, err
	}
	return l, nil
}

// DeleteLogglyInput is the input parameter to DeleteLoggly.
type DeleteLogglyInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Name is the name of the loggly to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteLoggly deletes the specified resource.
func (c *Client) DeleteLoggly(i *DeleteLogglyInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "loggly", i.Name)
	resp, err := c.Delete(path, CreateRequestOptions(i.Context))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := DecodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrStatusNotOk
	}
	return nil
}
